// Package syncmap menyediakan mekanisme dedupe/idempotency + resolve-ID generik untuk
// sync offline: pemetaan (device_id, local_id, entity_type) -> server_id, disimpan di
// tabel sync_id_map. Dipakai oleh SEMUA entity yang sync offline (transaksi, dan nanti
// shift, cash_drawer) supaya satu entity bisa "dicari" oleh entity lain dalam batch push
// yang sama lewat local_id-nya (mis. transaksi merujuk shift yang juga baru dibuka
// offline), dan supaya retry push tidak memproses ulang item yang sudah pernah berhasil.
//
// Sengaja jadi package berdiri sendiri (bukan bagian dari domain/sync atau
// domain/transaction) supaya bisa dipakai lintas-domain tanpa membuat kedua domain itu
// saling bergantung satu sama lain (sync_service sudah bergantung ke transaction repo;
// kalau mekanisme ini ditaruh di domain/sync, transaction_repo yang butuh memanggilnya
// akan menciptakan dependensi melingkar).
package syncmap

import "gorm.io/gorm"

const (
	resolveQuery = `SELECT server_id FROM sync_id_map WHERE device_id = ? AND local_id = ? AND entity_type = ? LIMIT 1`
	recordQuery  = `INSERT INTO sync_id_map (device_id, local_id, entity_type, server_id) VALUES (?, ?, ?, ?)`
)

// Resolve mencari server_id yang sudah pernah dicatat untuk kombinasi
// (deviceID, localID, entityType). found=false kalau belum pernah tercatat sama sekali
// (belum dedupe-worthy, atau referensi ke entity lain belum tersinkron).
func Resolve(db *gorm.DB, deviceID, localID, entityType string) (serverID int, found bool, err error) {
	if deviceID == "" || localID == "" {
		return 0, false, nil
	}
	row := db.Raw(resolveQuery, deviceID, localID, entityType).Row()
	if scanErr := row.Scan(&serverID); scanErr != nil {
		return 0, false, nil
	}
	return serverID, true, nil
}

// Record mencatat pemetaan (deviceID, localID, entityType) -> serverID setelah entity
// berhasil diterapkan (create). Dipanggil sekali per entity per local_id — panggilan
// kedua untuk kombinasi yang sama akan gagal karena unique constraint (seharusnya tidak
// pernah terjadi kalau caller sudah mengecek Resolve() lebih dulu sebelum apply).
func Record(db *gorm.DB, deviceID, localID, entityType string, serverID int) error {
	if deviceID == "" || localID == "" {
		return nil
	}
	return db.Exec(recordQuery, deviceID, localID, entityType, serverID).Error
}
