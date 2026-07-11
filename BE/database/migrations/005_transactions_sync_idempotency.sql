-- Dedupe untuk transaksi yang masuk lewat sync offline (BE/domain/sync). Tanpa ini,
-- push yang di-retry (mis. network flaky, klien tidak menerima response sukses lalu
-- coba lagi) bisa membuat transaksi dobel dan memotong stok dua kali. Kolom nullable
-- supaya transaksi checkout langsung (bukan dari sync) tidak terpengaruh — NULL tidak
-- dianggap bentrok oleh unique index di MySQL/MariaDB.
ALTER TABLE transactions
    ADD COLUMN sync_device_id VARCHAR(100) NULL AFTER device_source,
    ADD COLUMN sync_local_id  VARCHAR(36)  NULL AFTER sync_device_id,
    ADD UNIQUE KEY unique_sync_origin (sync_device_id, sync_local_id);
