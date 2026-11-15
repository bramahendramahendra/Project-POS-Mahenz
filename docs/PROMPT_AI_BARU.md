# Cara Pindah ke Claude Code

## Persiapan
1. Sudah commit (selesai). MySQL WAMP hidup.
2. Buka terminal di `d:\Develop\Project_pos_mahenz`, ketik `claude`.

## Satu prompt (salin seluruh blok, tempel ke Claude)

```
Kamu melanjutkan proyek POS ini dari AI lain (Kiro) yang kehabisan token. Balas dalam bahasa Indonesia.

1. Baca berurutan: CLAUDE.md, docs/HANDOFF_AI.md, docs/RANCANGAN_PORTAL_TOKO_COFFEE.md, README.md.
2. Sebelum mengubah apa pun, jawab singkat poin per poin (ambil dari dokumen, jangan menebak;
   tulis "tidak ada di dokumen" bila memang tidak ada):
   a. Port BE dan FE tiap aplikasi, dan nama DB tiap aplikasi.
   b. Role di Portal apa saja, dan bagaimana role Toko dipetakan ke Portal.
   c. Cara kerja SSO: umur kode, siapa yang menyimpan kunci aplikasi, dan di mana.
   d. Apa yang terjadi bila akun lokal Toko punya username sama dengan user SSO baru.
   e. Siapa yang boleh mengubah toggle login langsung.
   f. Tahap mana yang sudah selesai, mana yang menunggu persetujuan, mana berikutnya.
   g. Apa saja yang BELUM terverifikasi.
   h. Aturan kerja: soal commit, deploy, perubahan Toko, dan pencatatan keputusan.
   i. Cek beberapa klaim penting di dokumen langsung ke kode (misalnya port di config, nama tabel
      migrasi 011 Toko, endpoint SSO). Sebutkan bila ada yang tidak cocok.
3. Setelah menjawab, BERHENTI dan tunggu saya. Jangan mengubah kode atau menjalankan server dulu.
4. Setelah saya menjawab "benar", kerjakan tahap 3.6 (uji ujung ke ujung di dev; jalankan keempat server
   sendiri). Anggap tahap 3.5 sudah saya setujui. Ikuti "Ritme kerja" di CLAUDE.md: setelah tahap selesai,
   berhenti, lapor skenario satu per satu beserta hasilnya, dan tanya apakah lanjut.
```

## Setelah Claude menjawab
Cocokkan jawaban a sampai e dengan yang Anda tahu. Bila benar, ketik: `Benar, lanjut.` Bila ada yang salah, sebutkan bagian yang salah.

## Setelah itu
Pola yang sama seperti selama ini: Claude mengerjakan satu tahap, melapor, lalu bertanya "lanjut?".
Anda cukup menjawab "lanjut" (atau koreksi), dan commit sendiri memakai pesan yang ia sarankan.

Jika sesi ditutup lalu dibuka lagi, ketik:
`Baca CLAUDE.md dan docs/HANDOFF_AI.md, ringkas status 5 baris, lalu lanjutkan tahap berikutnya.`
