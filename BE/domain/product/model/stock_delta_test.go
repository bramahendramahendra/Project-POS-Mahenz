package model

import (
	"errors"
	"math"
	"testing"
)

// helper: bikin rantai 2 level Kardus (anchor) -> Botol, 24 Botol = 1 Kardus
// (rasio non-terminating 1/24 -- kasus bug precision utama yang jadi alasan
// seluruh rencana ini, lihat docs/RENCANA_PERBAIKAN_STOK_PRESISI.md).
func kardusBotolChain(kardusStock, kardusReserved, botolStock, botolReserved float64) []*ProductPackage {
	refQty1 := 1.0
	kardus := &ProductPackage{ID: 1, UnitID: 100, IsDefault: true, Qty: 1, Stock: kardusStock, ReservedQty: kardusReserved}
	botol := &ProductPackage{ID: 2, UnitID: 101, RefPackageID: intPtr(1), Qty: 24, RefQty: &refQty1, Stock: botolStock, ReservedQty: botolReserved}
	return []*ProductPackage{kardus, botol}
}

func intPtr(v int) *int { return &v }

func float64AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

// Skenario 1: rasio non-terminating (1/24) -- presisi tidak boleh hilang
// walau dipakai berulang kali (24x jual 1 Botol dari 1 Kardus harus persis
// habis ke 0, tidak ada sisa/drift akibat pembulatan seperti sistem lama).
func TestComputeStockDelta_NonTerminatingRatio_NoDriftAcrossRepeatedOps(t *testing.T) {
	packages := kardusBotolChain(1, 0, 0, 0) // total = 24 Botol

	for i := 0; i < 24; i++ {
		result, err := ComputeStockDelta(StockDeltaInput{
			Packages:  packages,
			PackageID: 2, // Botol
			Quantity:  1,
			Direction: StockOut,
		})
		if err != nil {
			t.Fatalf("jual ke-%d gagal: %v", i+1, err)
		}
		applyUpdates(packages, result.Updates)
	}

	kardus, botol := packages[0], packages[1]
	if !float64AlmostEqual(kardus.Stock, 0) || !float64AlmostEqual(botol.Stock, 0) {
		t.Fatalf("stok belum habis persis setelah 24x jual 1 Botol: Kardus=%v Botol=%v (harus 0,0 tanpa drift)", kardus.Stock, botol.Stock)
	}
}

// Skenario 1b: rasio non-terminating 1/3 -- pastikan ResolvePackageFactor
// (dipakai ulang oleh ComputeStockDelta) menghasilkan faktor presisi penuh,
// bukan yang sudah dibulatkan ke 3 desimal seperti kolom DB lama.
func TestComputeStockDelta_OneThirdRatio_PreciseFactor(t *testing.T) {
	refQty1 := 1.0
	slop := &ProductPackage{ID: 1, UnitID: 200, IsDefault: true, Qty: 1, Stock: 1, ReservedQty: 0}
	pcs := &ProductPackage{ID: 2, UnitID: 201, RefPackageID: intPtr(1), Qty: 3, RefQty: &refQty1, Stock: 0, ReservedQty: 0}
	packages := []*ProductPackage{slop, pcs}

	result, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  1,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	applyUpdates(packages, result.Updates)
	// 1 Slop = 3 Pcs. Jual 1 Pcs -> sisa 2 Pcs, 0 Slop.
	if !float64AlmostEqual(slop.Stock, 0) || !float64AlmostEqual(pcs.Stock, 2) {
		t.Fatalf("hasil salah: Slop=%v Pcs=%v (harus 0,2)", slop.Stock, pcs.Stock)
	}
}

// Skenario 2: jual lintas-level -- sisa lepasan Botol tidak cukup, harus
// "buka" 1 Kardus secara otomatis lewat cascading breakdown.
func TestComputeStockDelta_CrossLevelBorrow(t *testing.T) {
	packages := kardusBotolChain(1, 0, 4, 0) // 1 Kardus + 4 Botol = 28 Botol total

	result, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2, // Botol
		Quantity:  6,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	applyUpdates(packages, result.Updates)

	kardus, botol := packages[0], packages[1]
	if !float64AlmostEqual(kardus.Stock, 0) || !float64AlmostEqual(botol.Stock, 22) {
		t.Fatalf("borrow lintas-level salah: Kardus=%v Botol=%v (harus 0,22)", kardus.Stock, botol.Stock)
	}
}

// Skenario 3: penolakan saat stok tidak cukup -- tidak boleh jadi minus.
func TestComputeStockDelta_RejectInsufficientStock(t *testing.T) {
	packages := kardusBotolChain(0, 0, 2, 0) // cuma 2 Botol

	_, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  5,
		Direction: StockOut,
	})
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("harus ErrInsufficientStock, dapat: %v", err)
	}
	// pastikan tidak ada perubahan (fungsi murni, tidak punya efek samping)
	if packages[1].Stock != 2 {
		t.Fatalf("stok tidak boleh berubah saat ditolak, dapat Botol=%v", packages[1].Stock)
	}
}

// Skenario 5: reserved_qty tidak boleh ikut kepakai/kepinjam (celah #13) --
// kasus persis dari verifikasi ke-5: Kardus stock=1 reserved=1, Botol
// stock=23 (total 47 setara Botol), tapi yang BEBAS cuma 23. Jual 30 Botol
// harus ditolak walau totalnya (termasuk yang ditahan) cukup secara matematis.
func TestComputeStockDelta_ReservedQtyExcludedFromPool(t *testing.T) {
	packages := kardusBotolChain(1, 1, 23, 0) // free = 0*24 + 23 = 23 Botol

	_, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  30,
		Direction: StockOut,
	})
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("harus ditolak karena melebihi stok BEBAS (23), dapat: %v", err)
	}

	// tapi jual 23 (pas free-nya) harus berhasil, dan reserved tetap 1 di Kardus
	result, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  23,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("jual pas 23 (batas bebas) harusnya berhasil: %v", err)
	}
	applyUpdates(packages, result.Updates)
	if packages[0].Stock != 1 || packages[0].ReservedQty != 1 {
		t.Fatalf("baris Kardus yang ditahan retur tidak boleh ikut terpakai: Stock=%v ReservedQty=%v (harus tetap 1,1)", packages[0].Stock, packages[0].ReservedQty)
	}
	if packages[1].Stock != 0 {
		t.Fatalf("Botol harus habis ke 0, dapat %v", packages[1].Stock)
	}
}

// Skenario 6: produk dengan ref_package_id bercabang (celah #10) ditolak
// diproses otomatis, bukan menebak jawaban.
func TestComputeStockDelta_RejectBranchingChain(t *testing.T) {
	refQty1 := 1.0
	kardus := &ProductPackage{ID: 1, UnitID: 100, IsDefault: true, Qty: 1, Stock: 5}
	botolA := &ProductPackage{ID: 2, UnitID: 101, RefPackageID: intPtr(1), Qty: 24, RefQty: &refQty1, Stock: 0}
	botolB := &ProductPackage{ID: 3, UnitID: 102, RefPackageID: intPtr(1), Qty: 12, RefQty: &refQty1, Stock: 0} // sama-sama nunjuk Kardus (id 1) -> bercabang
	packages := []*ProductPackage{kardus, botolA, botolB}

	_, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  1,
		Direction: StockOut,
	})
	if !errors.Is(err, ErrBranchingChain) {
		t.Fatalf("harus ErrBranchingChain, dapat: %v", err)
	}
}

// Skenario 7: produk needs_stock_review=true menolak SEMUA operasi stok
// eksplisit (celah #20), tidak cuma andalkan stok=0 yang kebetulan menolak.
func TestComputeStockDelta_RejectNeedsStockReview(t *testing.T) {
	packages := kardusBotolChain(5, 0, 0, 0) // stok banyak, bukan 0 -- harus tetap ditolak

	_, err := ComputeStockDelta(StockDeltaInput{
		Packages:         packages,
		PackageID:        1,
		Quantity:         1,
		Direction:        StockOut,
		NeedsStockReview: true,
	})
	if !errors.Is(err, ErrNeedsStockReview) {
		t.Fatalf("harus ErrNeedsStockReview walau stok non-zero, dapat: %v", err)
	}

	// juga harus ditolak untuk arah IN (pembelian), bukan cuma OUT
	_, err = ComputeStockDelta(StockDeltaInput{
		Packages:         packages,
		PackageID:        1,
		Quantity:         1,
		Direction:        StockIn,
		NeedsStockReview: true,
	})
	if !errors.Is(err, ErrNeedsStockReview) {
		t.Fatalf("harus ErrNeedsStockReview untuk arah IN juga, dapat: %v", err)
	}
}

// Skenario 8: void selalu hitung ulang fresh dari package_id+quantity asli,
// TIDAK baca balik conversion_qty tersimpan -- stok harus kembali PERSIS ke
// angka semula tanpa drift (celah #21). Ini simulasi round-trip jual lalu
// void, keduanya lewat ComputeStockDelta yang sama (fresh setiap kali).
func TestComputeStockDelta_VoidRoundTrip_NoDrift(t *testing.T) {
	packages := kardusBotolChain(2, 0, 0, 0) // total 48 Botol

	sellResult, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  3,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("jual gagal: %v", err)
	}
	applyUpdates(packages, sellResult.Updates)

	// void: paket & qty PERSIS sama seperti transaksi asli, arah dibalik.
	voidResult, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  3,
		Direction: StockIn,
	})
	if err != nil {
		t.Fatalf("void gagal: %v", err)
	}
	applyUpdates(packages, voidResult.Updates)

	kardus, botol := packages[0], packages[1]
	if !float64AlmostEqual(kardus.Stock, 2) || !float64AlmostEqual(botol.Stock, 0) {
		t.Fatalf("stok setelah void harus PERSIS kembali ke semula (2,0), dapat Kardus=%v Botol=%v -- drift terdeteksi", kardus.Stock, botol.Stock)
	}
}

// applyUpdates mensimulasikan apa yang dilakukan lapisan repo: menulis hasil
// StockDeltaResult.Updates balik ke slice ProductPackage in-memory, supaya
// test bisa merangkai beberapa operasi berurutan (seperti transaksi nyata).
func applyUpdates(packages []*ProductPackage, updates []PackageStockUpdate) {
	byID := make(map[int]*ProductPackage, len(packages))
	for _, p := range packages {
		byID[p.ID] = p
	}
	for _, u := range updates {
		if p, ok := byID[u.PackageID]; ok {
			p.Stock = u.Stock
		}
	}
}
