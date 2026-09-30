-- =============================================================
-- remove_sync_center_20261001.sql
-- Menghapus TOTAL fitur "Sync Center" dari DATABASE production:
--   - seed menu 'operasional.sync' + akses role + route_registry '/sync'
--   - 4 tabel sync: sync_conflicts, sync_queue, sync_history, sync_id_map
--
-- Kode BE & FE untuk sync SUDAH dihapus (folder domain/sync, pkg/syncmap,
-- fitur FE operational/sync, metode sync di transaction/expense repo).
-- File migrasi 001 & 002 juga sudah dibersihkan untuk instalasi BARU — TAPI
-- migrasi pakai CREATE TABLE IF NOT EXISTS, jadi TIDAK menghapus tabel yang
-- SUDAH ada di DB production. Karena itu drop-nya harus MANUAL lewat skrip ini.
--
-- TIDAK dijalankan otomatis (lihat README). Jalankan manual + PENGAWASAN.
--
-- URUTAN AMAN:
--   1) BACKUP DATABASE DULU (wajib).
--   2) Jalankan BAGIAN 1 (CEK) dulu — lihat apakah tabel sync ada isinya.
--        - Kalau SEMUA 0 baris  -> aman, lanjut BAGIAN 2.
--        - Kalau ADA isinya      -> STOP, konsultasikan dulu (mungkin ada
--          data dari device desktop/Android yang belum tersinkron). Drop akan
--          MENGHILANGKAN data itu permanen.
--   3) Jalankan BAGIAN 2 (HAPUS).
--   4) Jalankan BAGIAN 3 (VERIFIKASI) — pastikan semua 0.
-- =============================================================

-- -------------------------------------------------------------
-- BAGIAN 1 — CEK (read-only, aman). Jalankan & baca dulu.
-- -------------------------------------------------------------
SELECT '== isi tabel sync (harus 0 semua kalau tak pakai offline) ==' AS info;
SELECT
  (SELECT COUNT(*) FROM sync_conflicts) AS conflicts,
  (SELECT COUNT(*) FROM sync_queue)     AS queue,
  (SELECT COUNT(*) FROM sync_history)   AS history,
  (SELECT COUNT(*) FROM sync_id_map)    AS id_map;

SELECT '== menu & route sync yang akan dihapus ==' AS info;
SELECT id, key_name, label, path FROM menus WHERE key_name = 'operasional.sync' OR path = '/sync';
SELECT * FROM route_registry WHERE path = '/sync';

-- >>> Kalau BAGIAN 1 menunjukkan ADA isi di tabel sync, BERHENTI di sini. <<<

-- -------------------------------------------------------------
-- BAGIAN 2 — HAPUS. Jalankan hanya setelah BAGIAN 1 dipastikan aman.
-- -------------------------------------------------------------
-- 2a) Seed menu, akses role, route registry
DELETE rma FROM role_menu_access rma
  JOIN menus m ON m.id = rma.menu_id
  WHERE m.key_name = 'operasional.sync';
DELETE FROM menus WHERE key_name = 'operasional.sync';
DELETE FROM route_registry WHERE path = '/sync';

-- 2b) Tabel sync (indeks ikut terhapus otomatis; tak ada FK antar tabel sync)
DROP TABLE IF EXISTS sync_conflicts;
DROP TABLE IF EXISTS sync_queue;
DROP TABLE IF EXISTS sync_history;
DROP TABLE IF EXISTS sync_id_map;

-- -------------------------------------------------------------
-- BAGIAN 3 — VERIFIKASI (semua harus 0)
-- -------------------------------------------------------------
SELECT '== sisa menu sync (harus 0) ==' AS info;
SELECT COUNT(*) AS menu_sync FROM menus WHERE key_name = 'operasional.sync';
SELECT '== sisa route /sync (harus 0) ==' AS info;
SELECT COUNT(*) AS route_sync FROM route_registry WHERE path = '/sync';
SELECT '== tabel sync tersisa (harus 0) ==' AS info;
SELECT COUNT(*) AS tabel_sync FROM information_schema.tables
  WHERE table_schema = DATABASE()
    AND table_name IN ('sync_conflicts','sync_queue','sync_history','sync_id_map');
