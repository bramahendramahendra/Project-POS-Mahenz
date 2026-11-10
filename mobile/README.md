# POS Mahenz — Pembungkus APK (Capacitor)

Folder ini membungkus aplikasi web POS menjadi APK Android (Opsi A: WebView membuka
`https://mahenz.duckdns.org`), plus plugin Bluetooth SPP agar bisa cetak langsung ke
printer EP58M dari tablet.

> Rancangan lengkap: lihat `../docs/RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`.

## Isi
- `capacitor.config.ts` — konfigurasi Capacitor (appId, server.url).
- `package.json` — dependency Capacitor + plugin `@kduma-autoid/capacitor-bluetooth-printer`.
- `www/` — placeholder (tidak ditampilkan; WebView memuat dari server.url).
- `android/` — di-generate oleh `npx cap add android` (tidak di-commit; dibuat ulang CI).

## Build APK (ringkas)
Build dilakukan di **GitHub Actions** (tidak perlu Android Studio di laptop).
1. Siapkan keystore + simpan GitHub Secrets (lihat dokumen rancangan bagian 10.3).
2. Push repo → buka tab **Actions** → jalankan workflow **Build Android APK**.
3. Unduh artifact `pos-mahenz-apk` → sideload `.apk` ke tablet.

Tanpa secret keystore, workflow tetap menghasilkan **APK debug** yang bisa di-install
untuk uji coba.

## Build lokal (opsional, bila punya Android SDK)
```bash
npm install
npx cap add android
npx cap sync android
cd android && ./gradlew assembleDebug
```

## Cetak Bluetooth di aplikasi
- Kode jembatan cetak ada di FE: `TOKO/FE/src/features/sales/cashier/printNativeBt.ts`
  (mengakses plugin lewat `window.Capacitor` saat runtime).
- Di modal cetak, pilihan **"Printer Bluetooth (Aplikasi)"** hanya muncul saat
  aplikasi dibuka lewat APK ini (bukan browser).
