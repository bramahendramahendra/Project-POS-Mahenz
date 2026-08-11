-- Menambahkan 'void' ke enum receivables.status.
-- Sebelumnya UPDATE receivables SET status = 'void' (saat transaksi kredit di-void)
-- disimpan MySQL sebagai string kosong karena 'void' tidak ada di enum, yang lalu
-- meruntuhkan StatusBadge di FE (key '' tidak ada di STATUS_MAP).
ALTER TABLE receivables MODIFY COLUMN status ENUM('unpaid','partial','paid','void') DEFAULT 'unpaid';

-- Perbaiki data yang sudah terlanjur korup jadi '' akibat bug di atas.
UPDATE receivables SET status = 'void' WHERE status = '';
