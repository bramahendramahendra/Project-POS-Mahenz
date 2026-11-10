# Cara Pindah ke Claude Code

## Persiapan
1. Sudah commit (selesai). MySQL WAMP hidup.
2. Buka terminal di `d:\Develop\Project_pos_mahenz`, ketik `claude`.

## Satu prompt (salin seluruh blok, tempel ke Claude)

```
Kamu melanjutkan proyek POS ini dari AI lain (Kiro) yang kehabisan token. Balas dalam bahasa Indonesia.

1. Baca berurutan: CLAUDE.md, docs/HANDOFF_AI.md, docs/RANCANGAN_PORTAL_TOKO_COFFEE.md, README.md.
2. Ringkas pemahamanmu maksimal 15 baris (tiga aplikasi, status tahap, keputusan terpenting,
   hal belum terverifikasi, dan pertanyaan bila dokumen dan kode bertentangan).
3. Lalu kerjakan tahap 3.6 (uji ujung ke ujung di dev; jalankan keempat server sendiri).
   Ikuti "Ritme kerja" di CLAUDE.md: setelah tahap selesai, berhenti, lapor, dan tanya apakah lanjut.
```

## Setelah itu
Pola yang sama seperti selama ini: Claude mengerjakan satu tahap, melapor, lalu bertanya "lanjut?".
Anda cukup menjawab "lanjut" (atau koreksi), dan commit sendiri memakai pesan yang ia sarankan.

Jika sesi ditutup lalu dibuka lagi, ketik:
`Baca CLAUDE.md dan docs/HANDOFF_AI.md, ringkas status 5 baris, lalu lanjutkan tahap berikutnya.`
