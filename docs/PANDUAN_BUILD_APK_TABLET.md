# Panduan Build & Pasang APK Tablet (Cetak Bluetooth EP58M)

> Status: **PANDUAN LANGKAH MANUAL** untuk dikerjakan user.
> Lapis kode (folder `mobile/`, modul FE, workflow GitHub Actions) sudah selesai.
> Dokumen ini hanya berisi langkah yang harus dikerjakan sendiri di laptop, GitHub, dan tablet.
> Latar belakang & rancangan lengkap: `docs/RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`.

---

## Gambaran Singkat

Tujuan: tablet bisa cetak struk langsung ke printer **EPPOS EP58M** (Bluetooth Classic/SPP)
lewat aplikasi APK pembungkus web POS.

Alur besar:

1. Siapkan alat di laptop (JDK) dan buat **keystore** (tanda tangan APK).
2. Simpan keystore ke **GitHub Secrets**.
3. Jalankan build APK di **GitHub Actions** (tidak berat di laptop).
4. Unduh APK, pasang di tablet, pair printer, uji cetak.

Yang dikerjakan AI (sudah selesai): kode FE, folder `mobile/`, file workflow CI.
Yang dikerjakan user (dokumen ini): Langkah 1 sampai 7 di bawah.

---

## Ringkasan 7 Langkah

| # | Langkah | Di mana | Frekuensi |
|---|---------|---------|-----------|
| 1 | Install JDK 17 | Laptop | Sekali |
| 2 | Buat keystore | Laptop | Sekali seumur aplikasi |
| 3 | Ubah keystore jadi Base64 | Laptop | Sekali |
| 4 | Simpan 4 GitHub Secrets | GitHub | Sekali |
| 5 | Jalankan build APK | GitHub Actions | Tiap ada perubahan native |
| 6 | Pasang APK di tablet | Tablet | Tiap tablet, sekali |
| 7 | Pair printer & uji cetak | Tablet | Sekali per tablet |

---

## LANGKAH 1 - Install JDK 17

JDK dibutuhkan hanya untuk perintah `keytool` (pembuat keystore). Jauh lebih ringan dari
Android Studio. Compile APK dikerjakan GitHub Actions, bukan laptop.

**Sebelum install, cek dulu apakah sudah ada:**

```
keytool -help
```

Kalau muncul daftar bantuan, JDK sudah terpasang. Lanjut ke Langkah 2.

**Cara install (lokasi bawaan installer sudah benar, tidak perlu dipindah):**

1. Buka **adoptium.net**, unduh **Temurin 17 (LTS)** untuk Windows x64 (file `.msi`).
2. Jalankan installer. Di langkah "Custom Setup", aktifkan:
   - **Add to PATH**
   - **Set JAVA_HOME variable**
3. Selesaikan instalasi.
4. **Tutup semua jendela PowerShell lalu buka yang baru** (PATH baru hanya terbaca di jendela baru).
5. Cek berhasil: jalankan `keytool -help`.

Alternatif lewat terminal (bila `winget` tersedia; ID paket belum diverifikasi, bila gagal pakai cara `.msi`):

```
winget install EclipseAdoptium.Temurin.17.JDK
```

---

## LANGKAH 2 - Buat Keystore

Keystore adalah "tanda tangan digital" APK. Android tidak mau memasang APK tanpa tanda tangan.
Dibuat **sekali**, lalu dipakai terus.

Jalankan di PowerShell (disarankan di folder aman di luar proyek, misalnya `D:\Backup\keystore\`):

```
keytool -genkeypair -v -keystore pos-release.keystore -alias pos -keyalg RSA -keysize 2048 -validity 10000
```

- Akan ditanya **password** dan beberapa data (nama, organisasi). Isi bebas.
- **CATAT password-nya.** Nanti dipakai di Langkah 4.
- Hasil: file `pos-release.keystore`.

**Penting:**
- Simpan file ini baik-baik dan buat cadangan (flashdisk / cloud pribadi).
- **Jangan simpan di dalam folder proyek dan jangan di-commit ke git.**
- Kalau file hilang, APK yang sudah terpasang di tablet tidak bisa di-update.
  Harus uninstall lalu pasang ulang.

---

## LANGKAH 3 - Ubah Keystore Jadi Base64

GitHub Secrets hanya menyimpan teks, jadi file keystore diubah ke teks base64 dulu.

Jalankan di folder tempat `pos-release.keystore` berada:

```
[Convert]::ToBase64String([IO.File]::ReadAllBytes("pos-release.keystore")) | Set-Content pos-release.keystore.b64.txt
```

- Buka file `pos-release.keystore.b64.txt`, salin **seluruh isinya** (satu baris panjang).
- Isi file ini dipakai di Langkah 4.
- Setelah selesai, hapus `pos-release.keystore.b64.txt` dari komputer
  (isinya sama sensitifnya dengan keystore).

---

## LANGKAH 4 - Simpan 4 Secret di GitHub

Di repo GitHub: **Settings > Secrets and variables > Actions > New repository secret**.

Buat 4 secret berikut (nama harus persis sama):

| Nama Secret | Isi |
|-------------|-----|
| `ANDROID_KEYSTORE_B64` | Isi file base64 dari Langkah 3 |
| `ANDROID_KEYSTORE_PASSWORD` | Password keystore (Langkah 2) |
| `ANDROID_KEY_ALIAS` | `pos` (alias yang dipakai di Langkah 2) |
| `ANDROID_KEY_PASSWORD` | Password key (biasanya sama dengan password keystore) |

Catatan: kalau 4 secret ini belum dibuat, workflow tetap jalan tapi hanya menghasilkan
**APK debug** (bisa dipasang untuk uji coba, tetapi bukan APK rilis bertanda tangan).

---

## LANGKAH 5 - Jalankan Build APK

Pastikan semua perubahan sudah di-`git push` ke GitHub (folder `mobile/`, workflow
`.github/workflows/android-apk.yml`, dan perubahan FE).

1. Buka repo GitHub, tab **Actions**.
2. Pilih workflow **Build Android APK**.
3. Klik **Run workflow** (workflow juga otomatis jalan bila ada push ke `main`
   yang menyentuh folder `mobile/`).
4. Tunggu beberapa menit sampai status hijau (sukses).
5. Buka run tersebut, bagian **Artifacts**, unduh **`pos-mahenz-apk`**.
6. Ekstrak zip, di dalamnya ada file `.apk`.

Kalau build merah (gagal), buka log run-nya dan kirim pesan error-nya ke AI untuk diperbaiki.
Kandidat penyebab: versi Capacitor/plugin di `mobile/package.json` tidak cocok.

---

## LANGKAH 6 - Pasang APK di Tablet

1. Pindahkan file `.apk` ke tablet (kabel USB, Google Drive, atau WhatsApp ke diri sendiri).
2. Di tablet, buka file `.apk` lalu tap **Install**.
3. Kalau muncul peringatan "Install from unknown sources", izinkan untuk aplikasi pembuka
   (misalnya File Manager atau Chrome), lalu ulangi install.
4. Aplikasi **POS Mahenz** muncul di layar tablet.

---

## LANGKAH 7 - Pair Printer dan Uji Cetak

1. Di tablet: **Settings > Bluetooth**, pair **EP58M** (sekali saja).
2. Buka aplikasi **POS Mahenz** (ikon APK, bukan Chrome), login.
3. Lakukan transaksi, tekan **Cetak**.
4. Pilih **"Printer Bluetooth (Aplikasi)"**.
   (Kartu ini **hanya muncul bila dibuka lewat APK**, tidak muncul di Chrome biasa.)
5. Saat pertama kali, aplikasi meminta **izin Bluetooth**. Pilih **Izinkan**.
6. Struk harus tercetak tajam dan tidak terpotong.

**Kriteria lulus:**
- Aplikasi tampil normal dan login/transaksi jalan.
- Kartu "Printer Bluetooth (Aplikasi)" muncul di modal cetak.
- EP58M terdeteksi dan struk tercetak.
- Di Chrome biasa (tanpa APK), kartu itu tidak muncul dan 8 metode lama tetap jalan.

---

## Kalau Ada Masalah Saat Uji

| Gejala | Kemungkinan penyebab | Tindakan |
|--------|----------------------|----------|
| Hasil cetak kacau / huruf aneh | Format data plugin bukan base64 | Beri tahu AI, sesuaikan `bytesToBase64` di `FE/src/features/sales/cashier/printNativeBt.ts` |
| Gagal connect, tidak ada pesan jelas (Android 12+) | Izin Bluetooth belum diberikan / belum ada di manifest | Cek izin aplikasi di Settings tablet; bila perlu AI menambah izin di `AndroidManifest` |
| "Tidak ada printer Bluetooth yang dipasangkan" | EP58M belum di-pair di tablet | Pair dulu di Settings > Bluetooth |
| Kartu "Printer Bluetooth (Aplikasi)" tidak muncul | Aplikasi dibuka lewat Chrome, bukan APK | Buka lewat ikon aplikasi POS Mahenz |
| Build Actions merah | Versi Capacitor/plugin tidak cocok | Kirim log error ke AI |
| Halaman tidak mau terbuka di APK | Domain/HTTPS bermasalah | Pastikan `https://mahenz.duckdns.org` bisa dibuka di browser tablet |

---

## Update ke Depan

- **Ubah tampilan/fitur web** (FE/BE): cukup deploy seperti biasa. Tablet langsung dapat
  versi baru otomatis. **Tidak perlu APK baru.**
- **Ubah bagian native** (plugin, izin, domain di `capacitor.config.ts`): ulangi Langkah 5 dan 6
  (build APK baru lalu pasang ulang). **Keystore tidak dibuat ulang**, pakai yang sama.

---

## Catatan Keamanan

- Jangan commit `*.keystore`, `*.b64.txt`, atau password ke repo (sudah di-ignore di `mobile/.gitignore`).
- Simpan cadangan keystore + password di tempat aman.
- Hapus file `pos-release.keystore.b64.txt` setelah secret tersimpan.
