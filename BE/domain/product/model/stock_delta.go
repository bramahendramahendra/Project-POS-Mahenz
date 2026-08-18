package model

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// StockDelta -- fungsi terpusat perhitungan stok per-level (Fase 2,
// docs/RENCANA_PERBAIKAN_STOK_PRESISI.md). Ini murni logika kalkulasi, tidak
// menyentuh DB sama sekali, supaya bisa diuji langsung tanpa mock/koneksi --
// pemanggilan DB (locking FOR UPDATE, baca/tulis) ada di lapisan repo
// terpisah yang memanggil ComputeStockDelta setelah data dimuat.

type StockDirection string

const (
	StockIn  StockDirection = "in"  // pembelian, restore/void penjualan, release retur
	StockOut StockDirection = "out" // penjualan, reduce retur, write-off
)

var (
	ErrInsufficientStock = errors.New("stok tidak mencukupi")
	ErrBranchingChain    = errors.New("rantai satuan produk bercabang (ref_package_id tidak linear), tidak bisa diproses otomatis")
	ErrNeedsStockReview  = errors.New("produk ini ditandai perlu ditinjau manual (needs_stock_review), operasi stok diblokir")
	ErrPackageNotFound   = errors.New("paket satuan tidak ditemukan pada produk ini")
	ErrNoAnchorPackage   = errors.New("produk tidak punya baris satuan anchor (is_default)")
	ErrNoPackages        = errors.New("produk tidak punya baris product_packages sama sekali")
)

// StockDeltaInput adalah input murni (in-memory) untuk satu operasi perubahan
// stok pada satu produk. Packages harus berisi SEMUA baris product_packages
// aktif (is_active=true) milik produk itu, sudah dimuat FOR UPDATE oleh
// pemanggil (lapisan repo) sebelum function ini dipanggil -- supaya locking
// dan kalkulasi konsisten dalam satu transaksi DB yang sama.
type StockDeltaInput struct {
	Packages         []*ProductPackage
	ContinuousUnitID map[int]bool // unit_id -> units.is_continuous
	PackageID        int          // paket acuan transaksi (mis. yang dijual/dibeli/diretur)
	Quantity         float64      // qty di satuan PackageID (integer untuk satuan diskrit)
	Direction        StockDirection
	NeedsStockReview bool // products.needs_stock_review milik produk ini
}

// PackageStockUpdate adalah nilai stok baru untuk satu baris product_packages,
// siap dipakai pemanggil untuk UPDATE ke DB.
type PackageStockUpdate struct {
	PackageID int
	Stock     float64
}

// StockDeltaResult adalah hasil kalkulasi lengkap. StockBeforeAnchor/After/
// MutationQtyAnchor dinyatakan dalam satuan anchor produk (bukan satuan
// terkecil) supaya konsisten dengan semantik lama stock_mutations
// (Aturan Operasional #5 di dokumen) -- audit trail tetap "total di satuan
// anchor", sekarang plus package_id per baris (celah #16) yang diisi
// pemanggil dari PackageID input asli (bukan dari sini).
type StockDeltaResult struct {
	Updates           []PackageStockUpdate
	AnchorPackageID   int
	StockBeforeAnchor float64
	StockAfterAnchor  float64
	MutationQtyAnchor float64
}

// ComputeStockDelta menghitung perubahan stok lintas-level untuk satu
// transaksi (Aturan Operasional #3, #4, #5, #6, #7 di dokumen desain):
//   1. Tolak total kalau produk ditandai needs_stock_review (celah #20).
//   2. Tolak total kalau rantai ref_package_id bercabang, bukan linear (celah #10).
//   3. Kurangi reserved_qty dulu dari pool yang boleh dipakai (celah #13).
//   4. Turunkan semua level ke total di satuan terkecil (integer, kecuali
//      level kontinu di posisi terkecil boleh desimal) -- Aturan #4.
//   5. Faktor konversi SELALU dihitung ulang fresh dari qty/ref_qty (lewat
//      ResolvePackageFactor), TIDAK PERNAH membaca nilai tersimpan yang
//      berpotensi sudah dibulatkan -- ini yang menutup celah #21 (drift saat
//      void) sekaligus bug precision utama yang jadi alasan seluruh rencana
//      ini (docs/RENCANA_PERBAIKAN_STOK_PRESISI.md bagian "Masalah").
func ComputeStockDelta(in StockDeltaInput) (*StockDeltaResult, error) {
	if in.NeedsStockReview {
		return nil, ErrNeedsStockReview
	}
	if len(in.Packages) == 0 {
		return nil, ErrNoPackages
	}

	anchor, smallest, factorToAnchor, err := analyzePackageChain(in.Packages)
	if err != nil {
		return nil, err
	}

	byID := make(map[int]*ProductPackage, len(in.Packages))
	for _, p := range in.Packages {
		byID[p.ID] = p
	}
	target, ok := byID[in.PackageID]
	if !ok {
		return nil, ErrPackageNotFound
	}

	smallestFactor := factorToAnchor[smallest.ID]

	unitsPerSmallest := func(p *ProductPackage) float64 {
		return factorToAnchor[p.ID] / smallestFactor
	}

	const epsilon = 1e-6

	var totalSmallestBefore, freeSmallestBefore float64
	for _, p := range in.Packages {
		ups := unitsPerSmallest(p)
		totalSmallestBefore += p.Stock * ups
		free := p.Stock - p.ReservedQty
		if free < 0 {
			free = 0
		}
		freeSmallestBefore += free * ups
	}

	targetUPS := unitsPerSmallest(target)
	deltaSmallest := in.Quantity * targetUPS

	var totalSmallestAfter, freeSmallestAfter float64
	switch in.Direction {
	case StockOut:
		if deltaSmallest > freeSmallestBefore+epsilon {
			return nil, ErrInsufficientStock
		}
		totalSmallestAfter = totalSmallestBefore - deltaSmallest
		freeSmallestAfter = freeSmallestBefore - deltaSmallest
	case StockIn:
		totalSmallestAfter = totalSmallestBefore + deltaSmallest
		freeSmallestAfter = freeSmallestBefore + deltaSmallest
	default:
		return nil, fmt.Errorf("arah stok tidak dikenal: %q", in.Direction)
	}
	if freeSmallestAfter < 0 {
		freeSmallestAfter = 0
	}
	if totalSmallestAfter < 0 {
		totalSmallestAfter = 0
	}

	// Breakdown ulang pool bebas (freeSmallestAfter) ke tiap level, dari
	// satuan terbesar ke terkecil (cascading floor) -- Aturan Operasional #4.
	// reserved_qty tiap baris TIDAK ikut dipecah ulang, ditambahkan balik
	// utuh di akhir supaya bagian yang ditahan retur tetap protected (celah #13).
	ordered := make([]*ProductPackage, len(in.Packages))
	copy(ordered, in.Packages)
	sort.Slice(ordered, func(i, j int) bool {
		return factorToAnchor[ordered[i].ID] > factorToAnchor[ordered[j].ID]
	})

	remaining := freeSmallestAfter
	updates := make([]PackageStockUpdate, 0, len(ordered))
	for i, p := range ordered {
		ups := unitsPerSmallest(p)
		isLast := i == len(ordered)-1
		continuous := in.ContinuousUnitID[p.UnitID]

		var freeQty float64
		if isLast && continuous {
			freeQty = remaining / ups
			remaining = 0
		} else {
			freeQty = math.Floor(remaining/ups + epsilon)
			remaining -= freeQty * ups
		}
		updates = append(updates, PackageStockUpdate{
			PackageID: p.ID,
			Stock:     freeQty + p.ReservedQty,
		})
	}
	// Sisa pecahan yang tidak habis dibagi (kasus satuan terkecil bukan
	// kontinu tapi rasio rantainya tidak bulat -- data lama yang belum
	// dibetulkan) numpuk ke level terkecil apa adanya, supaya tidak ada
	// stok yang hilang diam-diam tanpa jejak.
	if remaining > epsilon && len(updates) > 0 {
		lastUPS := unitsPerSmallest(ordered[len(ordered)-1])
		updates[len(updates)-1].Stock += remaining / lastUPS
	}

	anchorUPS := unitsPerSmallest(anchor)
	return &StockDeltaResult{
		Updates:           updates,
		AnchorPackageID:   anchor.ID,
		StockBeforeAnchor: totalSmallestBefore / anchorUPS,
		StockAfterAnchor:  totalSmallestAfter / anchorUPS,
		MutationQtyAnchor: deltaSmallest / anchorUPS,
	}, nil
}

// analyzePackageChain menemukan baris anchor & satuan terkecil (leaf) pada
// rantai product_packages sebuah produk, plus faktor konversi tiap baris
// relatif ke anchor (dihitung ulang fresh lewat ResolvePackageFactor, bukan
// baca nilai tersimpan). Dipakai bersama oleh ComputeStockDelta (jalur tulis)
// dan ComputeStockSummary (jalur baca, Fase 5) supaya logika "apa itu
// anchor/leaf/faktor" cuma ada satu tempat.
func analyzePackageChain(packages []*ProductPackage) (anchor, smallest *ProductPackage, factorToAnchor map[int]float64, err error) {
	childCount := make(map[int]int)
	for _, p := range packages {
		if p.IsDefault {
			anchor = p
		}
		if p.RefPackageID != nil {
			childCount[*p.RefPackageID]++
		}
	}
	if anchor == nil {
		return nil, nil, nil, ErrNoAnchorPackage
	}
	for _, count := range childCount {
		if count > 1 {
			return nil, nil, nil, ErrBranchingChain
		}
	}

	factorToAnchor = make(map[int]float64, len(packages))
	for _, p := range packages {
		f, ferr := ResolvePackageFactor(packages, p.ID)
		if ferr != nil {
			return nil, nil, nil, fmt.Errorf("resolve factor paket %d: %w", p.ID, ferr)
		}
		factorToAnchor[p.ID] = f
	}

	isParent := make(map[int]bool, len(packages))
	for _, p := range packages {
		if p.RefPackageID != nil {
			isParent[*p.RefPackageID] = true
		}
	}
	for _, p := range packages {
		if !isParent[p.ID] {
			smallest = p
			break
		}
	}
	if smallest == nil {
		return nil, nil, nil, errors.New("tidak ditemukan satuan terkecil pada rantai produk ini")
	}
	if factorToAnchor[smallest.ID] <= 0 {
		return nil, nil, nil, errors.New("faktor satuan terkecil tidak valid (<= 0)")
	}

	return anchor, smallest, factorToAnchor, nil
}

// StockSummary adalah hasil agregasi baca-saja lintas semua level satuan
// aktif sebuah produk (Fase 5, docs/RENCANA_PERBAIKAN_STOK_PRESISI.md).
// Tidak pernah menulis apa pun -- murni untuk jalur baca (list/detail produk,
// laporan, dashboard low-stock).
type StockSummary struct {
	AnchorStock    float64 // total semua level, dikonversi ke satuan anchor -- untuk tampilan angka tunggal (kompatibel dengan products.stock lama)
	AnchorReserved float64 // total reserved_qty semua level, dikonversi ke satuan anchor
	IsLowStock     bool    // celah #12: dibandingkan di satuan TERKECIL, bukan anchor-vs-anchor, supaya sisa < 1 unit anchor tidak salah alarm
}

// ComputeStockSummary menghitung StockSummary sebuah produk dari baris
// product_packages aktifnya. minStockAnchor adalah products.min_stock (selalu
// diinput di satuan anchor per Aturan Operasional #dsb) -- dikonversi ke
// satuan terkecil di sini sebelum dibandingkan (celah #12).
func ComputeStockSummary(packages []*ProductPackage, continuousUnitID map[int]bool, minStockAnchor float64) (*StockSummary, error) {
	active := make([]*ProductPackage, 0, len(packages))
	for _, p := range packages {
		if p.IsActive {
			active = append(active, p)
		}
	}
	if len(active) == 0 {
		return nil, ErrNoPackages
	}

	_, smallest, factorToAnchor, err := analyzePackageChain(active)
	if err != nil {
		return nil, err
	}
	smallestFactor := factorToAnchor[smallest.ID]

	unitsPerSmallest := func(p *ProductPackage) float64 {
		return factorToAnchor[p.ID] / smallestFactor
	}

	var anchorStock, anchorReserved, smallestFree float64
	for _, p := range active {
		f := factorToAnchor[p.ID]
		anchorStock += p.Stock * f
		anchorReserved += p.ReservedQty * f

		free := p.Stock - p.ReservedQty
		if free < 0 {
			free = 0
		}
		smallestFree += free * unitsPerSmallest(p)
	}

	// minStockAnchor (satuan anchor) -> satuan terkecil: kalikan dengan
	// unitsPerSmallest(anchor) = 1/smallestFactor (anchor punya factorToAnchor=1).
	minStockSmallest := minStockAnchor / smallestFactor
	const epsilon = 1e-6
	isLowStock := smallestFree <= minStockSmallest+epsilon

	return &StockSummary{
		AnchorStock:    anchorStock,
		AnchorReserved: anchorReserved,
		IsLowStock:     isLowStock,
	}, nil
}
