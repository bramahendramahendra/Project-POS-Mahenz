# One-off SQL Scripts

Folder ini berisi skrip SQL **koreksi data satu kali (one-off)** yang berkaitan
dengan insiden/perbaikan tertentu di production.

## PENTING: skrip di sini TIDAK dijalankan otomatis

Runner migrasi (`database/migrate.go` → `RunMigrations`) HANYA membaca file dari
`database/migrations/`. Skrip di folder `oneoff/` ini **tidak** ikut jalan saat
`go run` / boot BE. Semuanya harus dijalankan **manual** lewat klien MySQL,
dengan pengawasan.

## Kenapa dipisah dari `migrations/`

- Skrip di sini koreksi **data spesifik** akibat insiden (id baris tertentu),
  bukan evolusi skema yang harus ada di semua database.
- Butuh kontrol manusia: banyak yang menampilkan verifikasi `SELECT` sebelum
  `COMMIT` sehingga operator bisa memutuskan `COMMIT` atau `ROLLBACK`.
- Runner migrasi memecah file per `;` dan menjalankan tiap statement terpisah,
  sehingga blok `START TRANSACTION ... COMMIT` tidak berlaku utuh di sana —
  skrip di `oneoff/` dijalankan langsung oleh klien MySQL agar transaksinya utuh.

## Daftar skrip

- `cleanup_transaksi_cangkang_20260913.sql` — void 3 transaksi cangkang (header
  tanpa item) akibat bug atomicity, + koreksi `closing_balance` cash_drawer 21.
  Lihat `docs/BUG_TRANSAKSI_CANGKANG_HEADER_TANPA_ITEM.md`.

## Cara menjalankan

```powershell
# contoh (sesuaikan kredensial & path mysql client)
Get-Content .\cleanup_transaksi_cangkang_20260913.sql | & "C:\path\to\mysql.exe" -u <user> -p <database>
```

Baca output verifikasi. Jika ada baris "VERIFIKASI GAGAL", ganti `COMMIT;` di
akhir skrip menjadi `ROLLBACK;` sebelum menjalankan ulang, lalu investigasi.
