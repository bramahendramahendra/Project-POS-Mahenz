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

// Skenario 6: produk dengan dua satuan berfaktor konversi SAMA PERSIS (ambigu)
// ditolak diproses otomatis, bukan menebak jawaban. (Revisi celah #10: yang
// ditolak sekarang hanya ambiguitas faktor kembar, bukan sekadar percabangan.)
func TestComputeStockDelta_RejectAmbiguousEqualFactor(t *testing.T) {
	refQty1 := 1.0
	// BotolA & BotolB dua-duanya 12 unit = 1 Kardus -> faktor identik 1/12.
	// Tidak ada cara memutuskan "satuan terkecil" -> harus ditolak.
	kardus := &ProductPackage{ID: 1, UnitID: 100, IsDefault: true, Qty: 1, Stock: 5}
	botolA := &ProductPackage{ID: 2, UnitID: 101, RefPackageID: intPtr(1), Qty: 12, RefQty: &refQty1, Stock: 0}
	botolB := &ProductPackage{ID: 3, UnitID: 102, RefPackageID: intPtr(1), Qty: 12, RefQty: &refQty1, Stock: 0}
	packages := []*ProductPackage{kardus, botolA, botolB}

	_, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 2,
		Quantity:  1,
		Direction: StockOut,
	})
	if !errors.Is(err, ErrBranchingChain) {
		t.Fatalf("faktor kembar harus ErrBranchingChain, dapat: %v", err)
	}
}

// Skenario 6d: dua turunan berbeda yang sama-sama LEBIH KECIL dari anchor dan
// menunjuk anchor langsung (kasus Kardus/Renteng/Sachet dari produk nyata).
// Faktor berbeda (Sachet 1/200 < Renteng 1/20 < Kardus 1) -> tidak ambigu,
// harus BISA diproses. Ini yang dulu keliru ditolak guard lama.
func TestComputeStockDelta_TwoSmallerLeaves_DistinctFactor_Allowed(t *testing.T) {
	refQty1 := 1.0
	kardus := &ProductPackage{ID: 1, UnitID: 100, IsDefault: true, Qty: 1, Stock: 0}
	renteng := &ProductPackage{ID: 2, UnitID: 101, RefPackageID: intPtr(1), Qty: 20, RefQty: &refQty1, Stock: 0} // 20 Renteng = 1 Kardus
	sachet := &ProductPackage{ID: 3, UnitID: 102, RefPackageID: intPtr(1), Qty: 200, RefQty: &refQty1, Stock: 0} // 200 Sachet = 1 Kardus
	packages := []*ProductPackage{kardus, renteng, sachet}

	// Stok awal 1 Kardus = 200 Sachet. Jual 1 Sachet -> sisa 199 Sachet.
	kardus.Stock = 1
	result, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 3, // Sachet
		Quantity:  1,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("dua turunan faktor-beda harus bisa diproses, dapat error: %v", err)
	}
	applyUpdates(packages, result.Updates)

	// Breakdown 199 Sachet, urut terbesar->terkecil:
	// Kardus (200 sachet): floor(199/200)=0, sisa 199.
	// Renteng (10 sachet): floor(199/10)=19 (190), sisa 9.
	// Sachet (1): 9.
	if !float64AlmostEqual(kardus.Stock, 0) {
		t.Fatalf("Kardus harus 0, dapat %v", kardus.Stock)
	}
	if !float64AlmostEqual(renteng.Stock, 19) {
		t.Fatalf("Renteng harus 19, dapat %v", renteng.Stock)
	}
	if !float64AlmostEqual(sachet.Stock, 9) {
		t.Fatalf("Sachet harus 9, dapat %v", sachet.Stock)
	}
	// Total setara Kardus: 199/200.
	if !float64AlmostEqual(result.StockAfterAnchor, 199.0/200.0) {
		t.Fatalf("StockAfterAnchor harus 199/200, dapat %v", result.StockAfterAnchor)
	}
}

// Skenario 6b: "star topology" yang VALID (kasus rokok Surya 12) -- anchor di
// tengah rantai (Pack), dengan turunan lebih kecil (Batang) DAN lebih besar
// (Slop) sama-sama menunjuk anchor. Meski dua paket menunjuk satu induk, ini
// TIDAK ambigu karena hanya ada satu daun sejati (Batang) dan urutan faktornya
// tunggal (Batang 1/12 < Pack 1 < Slop 10). Harus BISA diproses, bukan ditolak.
func TestComputeStockDelta_StarTopology_AnchorInMiddle_Allowed(t *testing.T) {
	// Pack = anchor (id 32). Batang: 12 Batang = 1 Pack (id 354).
	// Slop: 1 Slop = 10 Pack (id 31). Batang & Slop dua-duanya ref ke Pack.
	refBatang := 1.0 // 12 Batang = 1 Pack
	refSlop := 10.0  // 1 Slop  = 10 Pack
	pack := &ProductPackage{ID: 32, UnitID: 300, IsDefault: true, Qty: 1, Stock: 0}
	batang := &ProductPackage{ID: 354, UnitID: 301, RefPackageID: intPtr(32), Qty: 12, RefQty: &refBatang, Stock: 0}
	slop := &ProductPackage{ID: 31, UnitID: 302, RefPackageID: intPtr(32), Qty: 1, RefQty: &refSlop, Stock: 0}
	packages := []*ProductPackage{pack, batang, slop}

	// Set stok awal: 2 Pack + 2 Batang (Slop 0). Total di satuan terkecil (Batang):
	// 2 Pack * 12 + 2 Batang = 26 Batang.
	pack.Stock = 2
	batang.Stock = 2

	// Jual 1 Batang -> harus berhasil (tidak ErrBranchingChain), sisa 25 Batang.
	result, err := ComputeStockDelta(StockDeltaInput{
		Packages:  packages,
		PackageID: 354, // Batang
		Quantity:  1,
		Direction: StockOut,
	})
	if err != nil {
		t.Fatalf("star topology valid harus bisa diproses, dapat error: %v", err)
	}
	applyUpdates(packages, result.Updates)

	// 25 Batang = 2 Slop? tidak. Breakdown terbesar->terkecil: Slop(120 batang) dulu,
	// tak cukup utk 1 Slop (butuh 120), lalu Pack: floor(25/12)=2 Pack (24 batang),
	// sisa 1 Batang. Jadi Slop=0, Pack=2, Batang=1.
	if !float64AlmostEqual(slop.Stock, 0) {
		t.Fatalf("Slop harus 0, dapat %v", slop.Stock)
	}
	if !float64AlmostEqual(pack.Stock, 2) {
		t.Fatalf("Pack harus 2, dapat %v", pack.Stock)
	}
	if !float64AlmostEqual(batang.Stock, 1) {
		t.Fatalf("Batang harus 1, dapat %v", batang.Stock)
	}

	// Total setara anchor (Pack) harus konsisten: 25 Batang = 25/12 Pack.
	if !float64AlmostEqual(result.StockAfterAnchor, 25.0/12.0) {
		t.Fatalf("StockAfterAnchor harus 25/12 Pack, dapat %v", result.StockAfterAnchor)
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
