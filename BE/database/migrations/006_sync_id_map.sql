-- Mekanisme dedupe/idempotency + resolve-ID lintas-entity generik untuk sync offline.
-- Menggantikan kolom khusus per-tabel (lihat migrasi 005) dengan satu tabel pemetaan
-- (device_id, local_id, entity_type) -> server_id yang dipakai SEMUA entity yang sync
-- offline (transaksi, dan nanti shift, cash_drawer) — bukan kolom unik per tabel, supaya
-- entity baru tidak perlu migrasi skema lagi tiap kali ditambahkan, dan entity lain bisa
-- "mencari" server_id sebuah item lewat local_id-nya (dibutuhkan saat mis. transaksi
-- offline merujuk shift yang JUGA baru dibuka offline dalam batch yang sama).
CREATE TABLE IF NOT EXISTS sync_id_map (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    device_id   VARCHAR(100) NOT NULL,
    local_id    VARCHAR(36)  NOT NULL,
    entity_type VARCHAR(50)  NOT NULL,
    server_id   INT          NOT NULL,
    created_at  DATETIME     DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_sync_origin (device_id, local_id, entity_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- sync_queue perlu tahu local_id item juga (sebelumnya tidak tersimpan sama sekali di
-- sini, cuma di sync_conflicts) supaya bisa dilacak balik ke sync_id_map dan supaya
-- transaksi -- yang sebelumnya tidak pernah tercatat di sync_queue sama sekali -- ikut
-- punya jejak audit yang sama seperti entity lain.
ALTER TABLE sync_queue
    ADD COLUMN local_id VARCHAR(36) NULL AFTER device_id;

-- Membatalkan migrasi 005: kolom sync_device_id/sync_local_id di transactions jadi
-- mubazir setelah sync_id_map ada (perannya digantikan sepenuhnya oleh tabel itu, dicek
-- lewat ResolveSyncMapping/RecordSyncMapping, bukan query langsung ke tabel transactions).
-- Dedupe transaksi TETAP ada, cuma pindah rumah ke mekanisme yang seragam dengan entity
-- lain (shift, cash_drawer) yang akan disusul di fase berikutnya.
ALTER TABLE transactions
    DROP INDEX unique_sync_origin,
    DROP COLUMN sync_device_id,
    DROP COLUMN sync_local_id;
