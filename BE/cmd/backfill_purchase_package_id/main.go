// Skrip one-off: backfill purchase_items.package_id untuk data lama.
//
// Latar belakang (celah #15, docs/RENCANA_PERBAIKAN_STOK_PRESISI.md):
// purchase_items.package_id 100% NULL di data live -- jalur tulis sudah
// diperbaiki (purchase_repo.go) supaya ke depan selalu terisi, tapi baris
// LAMA tetap perlu dicocokkan manual dari kolom `unit` (teks) yang tersimpan,
// karena kolom itu tidak bisa di-join balik dari dirinya sendiri. Ini
// prasyarat WAJIB sebelum Fase 3 (backfill stok dari riwayat) bisa jalan,
// karena backfill butuh package_id buat resolve faktor konversi tiap
// kejadian pembelian.
//
// Jalankan: go run ./cmd/backfill_purchase_package_id
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	model_product "pos_api/domain/product/model"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:@tcp(127.0.0.1:3306)/pos_retail_db?charset=utf8&parseTime=True&loc=Local"

type purchaseItemRow struct {
	ID            int
	ProductID     int
	Unit          string
	ConversionQty float64
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

	rows, err := db.Query(`SELECT id, product_id, unit, conversion_qty FROM purchase_items WHERE package_id IS NULL`)
	if err != nil {
		log.Fatalf("query purchase_items: %v", err)
	}
	var items []purchaseItemRow
	for rows.Next() {
		var it purchaseItemRow
		if err := rows.Scan(&it.ID, &it.ProductID, &it.Unit, &it.ConversionQty); err != nil {
			log.Fatalf("scan: %v", err)
		}
		items = append(items, it)
	}
	rows.Close()

	fmt.Printf("Total baris purchase_items.package_id NULL: %d\n", len(items))

	packagesByProduct := make(map[int][]*model_product.ProductPackage)
	var updated, skippedNoMatch, skippedAmbiguous int
	var skippedDetail []string

	for _, it := range items {
		pkgs, ok := packagesByProduct[it.ProductID]
		if !ok {
			pkgs, err = loadPackages(db, it.ProductID)
			if err != nil {
				log.Fatalf("load packages produk %d: %v", it.ProductID, err)
			}
			packagesByProduct[it.ProductID] = pkgs
		}

		candidates := make([]*model_product.ProductPackage, 0, 2)
		for _, p := range pkgs {
			if p.UnitName == it.Unit {
				candidates = append(candidates, p)
			}
		}

		var chosen *model_product.ProductPackage
		switch len(candidates) {
		case 0:
			skippedNoMatch++
			skippedDetail = append(skippedDetail, fmt.Sprintf("purchase_item %d (produk %d): unit '%s' tidak ditemukan di product_packages produk ini", it.ID, it.ProductID, it.Unit))
			continue
		case 1:
			chosen = candidates[0]
		default:
			// Lebih dari 1 paket pakai unit yang sama utk produk ini (mis.
			// varian "Slop (Lama)" vs "Slop") -- pilih yang faktor
			// konversinya paling dekat dengan conversion_qty tersimpan.
			bestDiff := math.MaxFloat64
			var ambiguous bool
			for _, c := range candidates {
				factor, err := model_product.ResolvePackageFactor(pkgs, c.ID)
				if err != nil {
					continue
				}
				diff := math.Abs(factor - it.ConversionQty)
				if diff < bestDiff-1e-9 {
					bestDiff = diff
					chosen = c
					ambiguous = false
				} else if math.Abs(diff-bestDiff) < 1e-9 {
					ambiguous = true
				}
			}
			if chosen == nil || ambiguous {
				skippedAmbiguous++
				skippedDetail = append(skippedDetail, fmt.Sprintf("purchase_item %d (produk %d): %d paket sama-sama unit '%s', tidak bisa dipastikan otomatis", it.ID, it.ProductID, len(candidates), it.Unit))
				continue
			}
		}

		if _, err := db.Exec(`UPDATE purchase_items SET package_id = ? WHERE id = ?`, chosen.ID, it.ID); err != nil {
			log.Fatalf("update purchase_item %d: %v", it.ID, err)
		}
		updated++
	}

	fmt.Printf("\n=== Laporan Backfill purchase_items.package_id ===\n")
	fmt.Printf("Berhasil diisi otomatis : %d\n", updated)
	fmt.Printf("Dilewati (unit tak ketemu): %d\n", skippedNoMatch)
	fmt.Printf("Dilewati (ambigu)         : %d\n", skippedAmbiguous)
	if len(skippedDetail) > 0 {
		fmt.Println("\nDetail yang dilewati (perlu ditinjau manual):")
		for _, d := range skippedDetail {
			fmt.Println(" - " + d)
		}
	}
}

func loadPackages(db *sql.DB, productID int) ([]*model_product.ProductPackage, error) {
	rows, err := db.Query(`
		SELECT pp.id, pp.ref_package_id, pp.qty, pp.ref_qty, COALESCE(u.name,'') AS unit_name, pp.is_default
		FROM product_packages pp
		JOIN units u ON u.id = pp.unit_id
		WHERE pp.product_id = ? AND pp.is_active = 1`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model_product.ProductPackage
	for rows.Next() {
		p := &model_product.ProductPackage{ProductID: productID}
		var refQty sql.NullFloat64
		if err := rows.Scan(&p.ID, &p.RefPackageID, &p.Qty, &refQty, &p.UnitName, &p.IsDefault); err != nil {
			return nil, err
		}
		if refQty.Valid {
			v := refQty.Float64
			p.RefQty = &v
		}
		out = append(out, p)
	}
	return out, nil
}
