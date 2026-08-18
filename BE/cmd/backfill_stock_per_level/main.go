// Skrip one-off Fase 3: backfill product_packages.stock (per-level) dari
// riwayat transaksi, BUKAN dari products.stock lama yang sudah berpotensi
// kepotong presisi. Lihat docs/RENCANA_PERBAIKAN_STOK_PRESISI.md bagian
// "Migrasi data lama & skema DB" > Bagian 2.
//
// Insight kunci: quantity ASLI di purchase_items/transaction_items/
// supplier_return_items/product_expiry_batches tidak rusak -- yang rusak
// cuma conversion_qty/products.stock hasil hitung yang kepotong. Skrip ini
// mereplay seluruh riwayat (urut stock_mutations.id) lewat ComputeStockDelta
// yang SAMA PERSIS dipakai jalur online nanti (Fase 2) -- bukan reimplementasi
// logika breakdown terpisah -- supaya hasilnya konsisten dengan cara sistem
// akan berjalan ke depan.
//
// Jalankan: go run ./cmd/backfill_stock_per_level
// HANYA untuk DB DEV -- skrip ini TIDAK dijalankan ke production di Fase 3.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"

	model_product "pos_api/domain/product/model"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:@tcp(127.0.0.1:3306)/pos_retail_db?charset=utf8&parseTime=True&loc=Local"

// toleransi selisih anchor-unit yang masih dianggap "efek truncation wajar"
// dari bug lama, bukan indikasi ada masalah data lain.
const reconcileTolerance = 0.01

type productRow struct {
	ID       int
	UnitID   sql.NullInt64
	OldStock float64
}

type mutationRow struct {
	ID            int
	MutationType  string
	ReferenceType string
	ReferenceID   sql.NullInt64
}

type reportEntry struct {
	ProductID int
	Name      string
	OldStock  float64
	NewAnchor float64
	Status    string // "ok" | "flagged"
	Reason    string
}

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	continuous, err := loadContinuousUnits(db)
	if err != nil {
		log.Fatalf("load units.is_continuous: %v", err)
	}

	products, err := loadProducts(db)
	if err != nil {
		log.Fatalf("load products: %v", err)
	}
	fmt.Printf("Total produk: %d\n\n", len(products))

	var okCount, flaggedCount int
	var report []reportEntry

	for _, p := range products {
		entry := reportEntry{ProductID: p.ID, OldStock: p.OldStock}
		if err := db.QueryRow(`SELECT name FROM products WHERE id = ?`, p.ID).Scan(&entry.Name); err != nil {
			entry.Name = fmt.Sprintf("produk %d", p.ID)
		}

		if err := ensureAnchorPackage(db, p); err != nil {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("gagal pastikan baris anchor: %v", err)
			flaggedCount++
			report = append(report, entry)
			continue
		}

		packages, err := loadPackages(db, p.ID)
		if err != nil || len(packages) == 0 {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("gagal muat product_packages: %v", err)
			flaggedCount++
			report = append(report, entry)
			continue
		}

		if branching, parentID := hasBranchingChain(packages); branching {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("rantai ref_package_id bercabang (lebih dari 1 anak menunjuk package id %d) -- celah #10, tidak diproses otomatis", parentID)
			flaggedCount++
			markNeedsReview(db, p.ID, entry.Reason)
			report = append(report, entry)
			continue
		}

		mutations, err := loadMutations(db, p.ID)
		if err != nil {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("gagal muat stock_mutations: %v", err)
			flaggedCount++
			markNeedsReview(db, p.ID, entry.Reason)
			report = append(report, entry)
			continue
		}

		replayOK, reason := replayHistory(db, packages, continuous, mutations, p.ID)
		if !replayOK {
			entry.Status = "flagged"
			entry.Reason = reason
			flaggedCount++
			markNeedsReview(db, p.ID, reason)
			report = append(report, entry)
			continue
		}

		newAnchor := anchorEquivalentTotal(packages)
		entry.NewAnchor = newAnchor

		diff := math.Abs(newAnchor - p.OldStock)
		if diff > reconcileTolerance {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("selisih rekonstruksi vs stok lama di luar wajar: lama=%.4f baru=%.4f selisih=%.4f (bukan sekadar efek truncation)", p.OldStock, newAnchor, diff)
			flaggedCount++
			markNeedsReview(db, p.ID, entry.Reason)
			report = append(report, entry)
			continue
		}

		if err := persistPackages(db, packages); err != nil {
			entry.Status = "flagged"
			entry.Reason = fmt.Sprintf("gagal simpan hasil: %v", err)
			flaggedCount++
			markNeedsReview(db, p.ID, entry.Reason)
			report = append(report, entry)
			continue
		}

		entry.Status = "ok"
		okCount++
		report = append(report, entry)
	}

	fmt.Printf("=== Laporan Backfill Fase 3 ===\n")
	fmt.Printf("Berhasil direkonstruksi & disimpan : %d\n", okCount)
	fmt.Printf("Ditandai perlu ditinjau (needs_stock_review): %d\n\n", flaggedCount)

	fmt.Println("--- Contoh produk kunci ---")
	printExample(report, 104, "Air Nestle Purelife")
	printExample(report, 58, "Coffee Candy Kapal Api")
	printExample(report, 197, "Djarum Super Kretek 12")

	if flaggedCount > 0 {
		fmt.Println("\n--- Semua produk yang ditandai perlu ditinjau ---")
		for _, e := range report {
			if e.Status == "flagged" {
				fmt.Printf(" - [%d] %s: %s\n", e.ProductID, e.Name, e.Reason)
			}
		}
	}
}

func printExample(report []reportEntry, id int, fallbackName string) {
	for _, e := range report {
		if e.ProductID == id {
			fmt.Printf("[%d] %s: lama=%.4f -> baru(anchor)=%.4f status=%s %s\n", e.ProductID, e.Name, e.OldStock, e.NewAnchor, e.Status, e.Reason)
			return
		}
	}
	fmt.Printf("[%d] %s: tidak ditemukan di data ini\n", id, fallbackName)
}

func loadContinuousUnits(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query(`SELECT id, is_continuous FROM units`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[int]bool)
	for rows.Next() {
		var id int
		var c bool
		if err := rows.Scan(&id, &c); err != nil {
			return nil, err
		}
		m[id] = c
	}
	return m, nil
}

func loadProducts(db *sql.DB) ([]productRow, error) {
	rows, err := db.Query(`SELECT id, unit_id, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []productRow
	for rows.Next() {
		var p productRow
		if err := rows.Scan(&p.ID, &p.UnitID, &p.OldStock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// ensureAnchorPackage menjamin tiap produk minimal punya 1 baris
// product_packages is_default=true (Aturan Operasional #1).
func ensureAnchorPackage(db *sql.DB, p productRow) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM product_packages WHERE product_id = ?`, p.ID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if !p.UnitID.Valid {
		return fmt.Errorf("produk tidak punya unit_id, tidak bisa buat anchor otomatis")
	}
	_, err := db.Exec(`INSERT INTO product_packages (product_id, unit_id, qty, ref_qty, ref_package_id, is_default, stock, reserved_qty, is_active) VALUES (?, ?, 1, NULL, NULL, 1, 0, 0, 1)`, p.ID, p.UnitID.Int64)
	return err
}

func loadPackages(db *sql.DB, productID int) ([]*model_product.ProductPackage, error) {
	rows, err := db.Query(`
		SELECT pp.id, pp.unit_id, pp.ref_package_id, pp.qty, pp.ref_qty, pp.is_default, pp.reserved_qty
		FROM product_packages pp
		WHERE pp.product_id = ? AND pp.is_active = 1`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model_product.ProductPackage
	for rows.Next() {
		p := &model_product.ProductPackage{ProductID: productID}
		var refQty sql.NullFloat64
		if err := rows.Scan(&p.ID, &p.UnitID, &p.RefPackageID, &p.Qty, &refQty, &p.IsDefault, &p.ReservedQty); err != nil {
			return nil, err
		}
		if refQty.Valid {
			v := refQty.Float64
			p.RefQty = &v
		}
		p.Stock = 0 // mulai dari nol, direplay dari riwayat -- bukan copy dari products.stock lama
		out = append(out, p)
	}
	return out, nil
}

func hasBranchingChain(packages []*model_product.ProductPackage) (bool, int) {
	childCount := make(map[int]int)
	for _, p := range packages {
		if p.RefPackageID != nil {
			childCount[*p.RefPackageID]++
		}
	}
	for parentID, count := range childCount {
		if count > 1 {
			return true, parentID
		}
	}
	return false, 0
}

func loadMutations(db *sql.DB, productID int) ([]mutationRow, error) {
	rows, err := db.Query(`SELECT id, mutation_type, reference_type, reference_id FROM stock_mutations WHERE product_id = ? ORDER BY id ASC`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []mutationRow
	for rows.Next() {
		var m mutationRow
		if err := rows.Scan(&m.ID, &m.MutationType, &m.ReferenceType, &m.ReferenceID); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

type historyEvent struct {
	PackageID int
	Quantity  float64
	Direction model_product.StockDirection
}

// replayHistory mereplay tiap mutasi lewat ComputeStockDelta yang sama
// dipakai jalur online (Fase 2) -- supaya hasil backfill konsisten dengan
// logika yang akan berjalan ke depan, bukan reimplementasi terpisah.
func replayHistory(db *sql.DB, packages []*model_product.ProductPackage, continuous map[int]bool, mutations []mutationRow, productID int) (ok bool, reason string) {
	for _, m := range mutations {
		if !m.ReferenceID.Valid {
			return false, fmt.Sprintf("mutasi id %d tidak punya reference_id, tidak bisa ditelusuri sumbernya", m.ID)
		}
		refID := int(m.ReferenceID.Int64)

		events, err := resolveEventsForMutation(db, productID, m.MutationType, m.ReferenceType, refID)
		if err != nil {
			return false, fmt.Sprintf("mutasi id %d (%s/%s ref=%d): %v", m.ID, m.MutationType, m.ReferenceType, refID, err)
		}

		for _, ev := range events {
			result, err := model_product.ComputeStockDelta(model_product.StockDeltaInput{
				Packages:         packages,
				ContinuousUnitID: continuous,
				PackageID:        ev.PackageID,
				Quantity:         ev.Quantity,
				Direction:        ev.Direction,
			})
			if err != nil {
				return false, fmt.Sprintf("replay mutasi id %d gagal (%s): %v", m.ID, m.MutationType, err)
			}
			applyUpdates(packages, result.Updates)
		}
	}
	return true, ""
}

func applyUpdates(packages []*model_product.ProductPackage, updates []model_product.PackageStockUpdate) {
	byID := make(map[int]*model_product.ProductPackage, len(packages))
	for _, p := range packages {
		byID[p.ID] = p
	}
	for _, u := range updates {
		if p, ok := byID[u.PackageID]; ok {
			p.Stock = u.Stock
		}
	}
}

func resolveEventsForMutation(db *sql.DB, productID int, mutationType, referenceType string, refID int) ([]historyEvent, error) {
	switch mutationType {
	case "in":
		return purchaseItemEvents(db, productID, refID, model_product.StockIn)
	case "void_purchase":
		return purchaseItemEvents(db, productID, refID, model_product.StockOut)
	case "out":
		return transactionItemEvents(db, productID, refID, model_product.StockOut)
	case "void":
		return transactionItemEvents(db, productID, refID, model_product.StockIn)
	case "return":
		return supplierReturnItemEvents(db, productID, refID, model_product.StockOut)
	case "expired":
		return expiryBatchEvents(db, productID, refID)
	default:
		return nil, fmt.Errorf("tipe mutasi %q tidak dikenal/tidak bisa direkonstruksi otomatis", mutationType)
	}
}

func purchaseItemEvents(db *sql.DB, productID, purchaseID int, dir model_product.StockDirection) ([]historyEvent, error) {
	rows, err := db.Query(`SELECT quantity, package_id FROM purchase_items WHERE purchase_id = ? AND product_id = ?`, purchaseID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []historyEvent
	for rows.Next() {
		var qty float64
		var pkgID sql.NullInt64
		if err := rows.Scan(&qty, &pkgID); err != nil {
			return nil, err
		}
		if !pkgID.Valid {
			return nil, fmt.Errorf("purchase_item produk %d di purchase %d tidak punya package_id (harusnya sudah dibackfill di Prasyarat #1)", productID, purchaseID)
		}
		out = append(out, historyEvent{PackageID: int(pkgID.Int64), Quantity: qty, Direction: dir})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("tidak ada purchase_items utk purchase %d produk %d", purchaseID, productID)
	}
	return out, nil
}

func transactionItemEvents(db *sql.DB, productID, transactionID int, dir model_product.StockDirection) ([]historyEvent, error) {
	rows, err := db.Query(`SELECT quantity, unit_id FROM transaction_items WHERE transaction_id = ? AND product_id = ?`, transactionID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []historyEvent
	for rows.Next() {
		var qty float64
		var pkgID sql.NullInt64
		if err := rows.Scan(&qty, &pkgID); err != nil {
			return nil, err
		}
		if !pkgID.Valid {
			return nil, fmt.Errorf("transaction_item produk %d di transaksi %d tidak punya unit_id (package_id)", productID, transactionID)
		}
		out = append(out, historyEvent{PackageID: int(pkgID.Int64), Quantity: qty, Direction: dir})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("tidak ada transaction_items utk transaksi %d produk %d", transactionID, productID)
	}
	return out, nil
}

func supplierReturnItemEvents(db *sql.DB, productID, returnID int, dir model_product.StockDirection) ([]historyEvent, error) {
	rows, err := db.Query(`SELECT quantity, unit, package_id FROM supplier_return_items WHERE return_id = ? AND product_id = ?`, returnID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type raw struct {
		qty    float64
		unit   string
		pkgID  sql.NullInt64
	}
	var raws []raw
	for rows.Next() {
		var r raw
		if err := rows.Scan(&r.qty, &r.unit, &r.pkgID); err != nil {
			return nil, err
		}
		raws = append(raws, r)
	}
	if len(raws) == 0 {
		return nil, fmt.Errorf("tidak ada supplier_return_items utk retur %d produk %d", returnID, productID)
	}

	packages, err := loadPackages(db, productID)
	if err != nil {
		return nil, err
	}

	var out []historyEvent
	for _, r := range raws {
		if r.pkgID.Valid {
			out = append(out, historyEvent{PackageID: int(r.pkgID.Int64), Quantity: r.qty, Direction: dir})
			continue
		}
		// package_id belum diisi (celah #14 masih Fase 4) -- cocokkan lewat
		// nama unit yang tersimpan, sama seperti backfill purchase_items.
		pkgID, err := matchPackageByUnitName(db, packages, r.unit)
		if err != nil {
			return nil, fmt.Errorf("retur %d: %w", returnID, err)
		}
		out = append(out, historyEvent{PackageID: pkgID, Quantity: r.qty, Direction: dir})
	}
	return out, nil
}

func expiryBatchEvents(db *sql.DB, productID, batchID int) ([]historyEvent, error) {
	var qty float64
	var purchaseItemID int
	if err := db.QueryRow(`SELECT qty, purchase_item_id FROM product_expiry_batches WHERE id = ? AND product_id = ?`, batchID, productID).Scan(&qty, &purchaseItemID); err != nil {
		return nil, fmt.Errorf("batch expired %d tidak ditemukan utk produk %d: %w", batchID, productID, err)
	}
	var pkgID sql.NullInt64
	if err := db.QueryRow(`SELECT package_id FROM purchase_items WHERE id = ?`, purchaseItemID).Scan(&pkgID); err != nil {
		return nil, fmt.Errorf("gagal cari package_id dari purchase_item %d: %w", purchaseItemID, err)
	}
	if !pkgID.Valid {
		return nil, fmt.Errorf("purchase_item %d (sumber batch %d) tidak punya package_id", purchaseItemID, batchID)
	}
	return []historyEvent{{PackageID: int(pkgID.Int64), Quantity: qty, Direction: model_product.StockOut}}, nil
}

func matchPackageByUnitName(db *sql.DB, packages []*model_product.ProductPackage, unitName string) (int, error) {
	rows, err := db.Query(`SELECT pp.id, u.name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.id IN (` + placeholders(len(packages)) + `)`, packageIDArgs(packages)...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var matchID int
	found := 0
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return 0, err
		}
		if name == unitName {
			matchID = id
			found++
		}
	}
	if found == 0 {
		return 0, fmt.Errorf("unit '%s' tidak ditemukan di product_packages", unitName)
	}
	if found > 1 {
		return 0, fmt.Errorf("unit '%s' ambigu (%d paket cocok)", unitName, found)
	}
	return matchID, nil
}

func placeholders(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			s += ","
		}
		s += "?"
	}
	return s
}

func packageIDArgs(packages []*model_product.ProductPackage) []any {
	args := make([]any, len(packages))
	for i, p := range packages {
		args[i] = p.ID
	}
	return args
}

// anchorEquivalentTotal menjumlahkan semua level balik ke satuan anchor,
// dipakai buat rekonsiliasi terhadap products.stock lama (Aturan #5: audit
// tetap "total di satuan anchor").
func anchorEquivalentTotal(packages []*model_product.ProductPackage) float64 {
	var anchor *model_product.ProductPackage
	for _, p := range packages {
		if p.IsDefault {
			anchor = p
			break
		}
	}
	if anchor == nil {
		return 0
	}
	factorToAnchor := make(map[int]float64, len(packages))
	for _, p := range packages {
		f, err := model_product.ResolvePackageFactor(packages, p.ID)
		if err != nil {
			return 0
		}
		factorToAnchor[p.ID] = f
	}
	var total float64
	for _, p := range packages {
		total += p.Stock * factorToAnchor[p.ID]
	}
	return total
}

func persistPackages(db *sql.DB, packages []*model_product.ProductPackage) error {
	sort.Slice(packages, func(i, j int) bool { return packages[i].ID < packages[j].ID })
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, p := range packages {
		if _, err := tx.Exec(`UPDATE product_packages SET stock = ? WHERE id = ?`, p.Stock, p.ID); err != nil {
			tx.Rollback()
			return err
		}
	}
	// bersihkan flag needs_stock_review kalau sebelumnya sempat ditandai
	// (mis. dari percobaan run sebelumnya) tapi sekarang berhasil direkonstruksi
	if len(packages) > 0 {
		if _, err := tx.Exec(`UPDATE products SET needs_stock_review = 0, stock_review_note = NULL WHERE id = ?`, packages[0].ProductID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func markNeedsReview(db *sql.DB, productID int, note string) {
	_, _ = db.Exec(`UPDATE products SET needs_stock_review = 1, stock_review_note = ? WHERE id = ?`, note, productID)
}
