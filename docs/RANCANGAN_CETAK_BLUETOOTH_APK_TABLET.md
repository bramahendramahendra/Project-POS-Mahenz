# Rancangan Cetak Struk Bluetooth untuk Tablet (Bungkus Web jadi APK)

> Status: **RANCANGAN — belum ada implementasi.**
> Tujuan: papan kerja bersama untuk menambah cetak struk langsung ke printer
> Bluetooth EP58M **dari tablet Android**, dengan cara membungkus aplikasi web
> menjadi APK (WebView + jembatan Bluetooth native). Dibahas dulu sampai
> disepakati sebelum menulis kode.
> Bahasa sederhana + catatan teknis agar bisa direview.

---

## 1. Latar Belakang / Masalah

- Aplikasi POS berjalan sebagai **web** (React/Vite di FE, Go + MySQL di BE pada VPS).
- Dipakai di **desktop/laptop** dan **tablet**.
- Printer struk yang dipakai: **EPPOS EP58M** (thermal 58mm, Bluetooth).
  - No. Product: BG100452.
  - Jenis koneksi: **Bluetooth Classic / SPP** (bukan Bluetooth LE).
- **Desktop sudah BERES**: cetak via jalur "Cetak via Server" (BE menulis ESC/POS
  langsung ke COM port, printer terhubung ke komputer via Bluetooth yang dipetakan
  ke COM3). Hasil tajam & tidak terpotong. **Desktop tidak dibahas lagi di sini.**
- **Tablet BELUM beres**: dari browser (Chrome) tidak ada jalur yang bisa mencetak
  langsung & native ke EP58M.

### Akar masalah tablet (penting)
Browser **tidak bisa** mengakses **Bluetooth Classic/SPP**. Web Bluetooth hanya
mendukung **Bluetooth LE (BLE)**. Karena EP58M memakai SPP, maka:
- Metode "Printer Bluetooth (Web Bluetooth)" hampir pasti **gagal** di tablet
  (printer tidak terdeteksi).
- Ini **batas teknologi browser**, bukan bug kode aplikasi.

---

## 2. Metode Cetak yang SUDAH Ada (per 2026, sebelum rancangan ini)

Saat tombol **Cetak** ditekan, muncul modal pemilih metode (manual, tidak ada
auto-print). Total **8 metode**:

| # | Metode | Inti | Cocok Tablet + EP58M? |
|---|--------|------|-----------------------|
| 1 | Cetak via Server | BE tulis ESC/POS ke COM port komputer server | ⚠️ Bisa, tapi butuh 1 komputer selalu nyala + terhubung printer + 1 jaringan |
| 2 | Printer Jaringan / LAN | BE kirim ESC/POS via TCP ke IP printer (port 9100) | ❌ EP58M tidak punya LAN |
| 3 | Printer USB / Serial | Web Serial dari browser ke COM | ❌ Web Serial tidak ada di Android |
| 4 | Printer Bluetooth (Web Bluetooth) | BLE langsung dari browser | ❌ EP58M = SPP, bukan BLE |
| 5 | Printer via RawBT | Lewat aplikasi RawBT (Android) | ⚠️ Secara teori bisa, tapi di lapangan gagal terus |
| 6 | Simpan / Bagikan PDF | Buat file PDF struk | ✅ Jalan, tapi bukan cetak langsung (buka PDF lalu cetak manual) |
| 7 | Kirim via WhatsApp | Bagikan ringkasan struk | ✅ Jalan, tapi bukan cetak kertas |
| 8 | Browser / Dialog Print | Dialog print bawaan browser | ⚠️ Ke Bluetooth di tablet biasanya ribet |

**Kesimpulan:** belum ada satu pun jalur yang mencetak **langsung + native** ke
EP58M dari tablet. Itulah yang mau diselesaikan rancangan ini.

---

## 3. Solusi yang Dipilih: Bungkus Web jadi APK (WebView + jembatan Bluetooth)

### Konsep
- **Bukan** menulis ulang aplikasi jadi native Android. Kode FE **tetap React/Vite
  yang sekarang**.
- APK hanya "cangkang": berisi **WebView** yang menampilkan aplikasi web, **plus**
  **jembatan Bluetooth native** yang punya akses ke Bluetooth Classic/SPP (yang
  tidak bisa diakses browser).
- Aplikasi web memanggil jembatan itu (lewat JavaScript) untuk mengirim ESC/POS ke
  EP58M → cetak langsung, native, tajam (seperti e-nota), tanpa RawBT.

### Keputusan yang SUDAH disepakati dalam diskusi
1. **Pakai Opsi A**: WebView menunjuk ke **URL server/domain** (bukan bundle offline).
   - FE tetap dilayani dari VPS/domain seperti sekarang.
   - Update FE di server → **otomatis** tampil di tablet, tanpa rebuild APK.
   - APK tipis (cuma cangkang + jembatan BT).
   - Konsekuensi: tablet **butuh internet** (sama seperti sekarang).
2. **VPS + domain TETAP diperlukan.** APK tidak menggantikan server.
   - BE (Go) + MySQL tetap di VPS (ini "otak" aplikasi: login, data produk,
     transaksi, dll).
   - Domain untuk FE tetap dipakai (WebView membukanya).
3. **APK WAJIB diinstall di tiap tablet** (sekali). Kemampuan Bluetooth ada di
   jembatan native di dalam APK; tanpa install APK, tablet = web biasa (tidak bisa
   cetak Bluetooth EP58M).
   - Setelah terpasang: update web otomatis. Install ulang APK hanya bila bagian
     **native** berubah (ganti plugin BT, ganti domain hardcoded, dsb) — jarang.

### Framework
- **Rekomendasi: Capacitor** (penerus Cordova, cocok dengan stack React/Vite,
  ekosistem plugin Bluetooth matang, mudah panggil native dari JS).
- **TWA (Trusted Web Activity) TIDAK dipakai** — TWA pada dasarnya tetap Chrome,
  jadi Bluetooth Classic tetap tidak bisa.

### Plugin Bluetooth
- Butuh plugin **Bluetooth Classic / SPP** (bukan BLE) yang menyediakan:
  scan device ter-pair, connect, dan **write bytes** (untuk kirim ESC/POS).
- Alur pakai: pair EP58M di Settings Android (sekali) → pilih EP58M di aplikasi →
  kirim byte ESC/POS → cetak.
- **Risiko utama (harus dicek di awal implementasi):** memastikan plugin SPP yang
  dipilih masih terawat & kompatibel dengan versi Capacitor + versi Android tablet.

---

## 4. Dampak ke Sistem yang Sudah Jalan (jaminan)

Rancangan ini bersifat **aditif** — menambah, bukan mengganti.

| Skenario | Website (browser) | 8 metode print lama | Metode Bluetooth APK |
|----------|-------------------|---------------------|----------------------|
| APK tidak diinstall, buka di browser | Normal | Semua jalan | Tidak muncul |
| APK diinstall, buka lewat APK | Normal | Semua jalan | Muncul (tambahan) |

Prinsip yang akan dijaga saat implementasi:
1. Metode "Bluetooth Native (APK)" ditambahkan sebagai **pilihan tambahan** di modal
   cetak, dan **hanya muncul bila aplikasi sedang jalan di dalam APK** (dideteksi
   lewat flag Capacitor).
2. Modul jembatan Bluetooth **terisolasi**: bila tidak ada (di browser), aplikasi
   tidak error — pilihannya cukup tidak ditampilkan.
3. Tidak ada metode lama yang dicabut/diubah perilakunya.

---

## 5. Distribusi APK (opsi)
- **Sideload** (copy file `.apk` ke tablet lalu install): paling simpel, cocok untuk
  beberapa tablet toko sendiri. **Kemungkinan cukup untuk kebutuhan sekarang.**
- **Play Store** (internal/publik): lebih resmi & update terkelola, butuh akun
  developer Google (biaya sekali). Opsional, untuk skala besar.

---

## 6. Trade-off / Konsekuensi (jujur)
1. Menambah **jalur rilis baru**: perlu build tool Android (Android Studio / Android
   SDK + JDK) di komputer build. Setup sekali, tapi tidak ringan.
2. Install APK sekali per tablet; rebuild+install ulang hanya bila bagian native
   berubah.
3. Ada **kode native tambahan** yang perlu dirawat (plugin BT, versi Android,
   signing key APK).
4. Effort lebih besar dari sekadar util cetak (setup Capacitor, pasang plugin,
   lapisan jembatan di FE, build, uji di tablet).

### Alternatif yang DITOLAK / ditunda
- **RawBT**: sudah dicoba, gagal terus di lapangan. Tidak diandalkan.
- **Printer WiFi/LAN (ganti hardware)**: solusi paling bersih untuk tablet tanpa
  coding native, TAPI butuh beli printer baru + ada WiFi. Ditunda karena user ingin
  tetap pakai EP58M.

---

## 7. Pertanyaan Terbuka (dibahas sebelum/di awal implementasi)
1. Tablet Android versi berapa? (Android 12+ punya model izin Bluetooth runtime yang
   berbeda — mempengaruhi kode izin & kompatibilitas plugin.)
2. Tim siap dengan setup build Android (Android Studio) + proses signing & distribusi
   APK?
3. Domain yang akan dibuka WebView: pakai domain produksi yang sekarang? (idealnya
   HTTPS.)
4. Verifikasi plugin SPP Capacitor yang masih terawat & jalan di Android target.

---

## 8. Rencana Teknis Detail

Disusun berdasarkan pemeriksaan kode nyata (FE Vite/React + axios) dan verifikasi
plugin. **Belum ada kode ditulis** — ini blueprint untuk disepakati.

### 8.0 Temuan kode yang relevan
- FE pakai Vite. Alamat API dibaca dari `import.meta.env.VITE_API_URL`
  (`src/services/api.client.ts` → `axios.create({ baseURL: VITE_API_URL })`).
  - `.env.development`: `VITE_API_URL=http://localhost:8080/api`
  - `.env.production`: `VITE_API_URL=/api` ← **RELATIF**.
- Sudah ada `VITE_PLATFORM=web` di env (bisa dimanfaatkan sebagai penanda).
- Byte struk ESC/POS sudah tersedia: `buildReceiptEscposBytes(data)` di
  `FE/src/features/sales/cashier/blePrinter.ts` → tinggal dipakai ulang.
- Modal pemilih metode + daftar `PRINT_OPTIONS` + `runPrintMethod` ada di
  `FE/src/features/sales/cashier/components/ReceiptPrint.tsx`.

> **Catatan penting (Opsi A + API relatif):** di produksi `VITE_API_URL=/api`
> (relatif). Ini jalan di browser karena FE & API diakses lewat domain yang sama
> (nginx reverse proxy). Di dalam APK Opsi A, WebView membuka domain yang sama,
> jadi path relatif `/api` **tetap valid** (base URL = domain yang dibuka WebView).
> Jadi kemungkinan besar **tidak perlu ubah konfigurasi API**. Ini akan diverifikasi
> saat uji APK; bila bermasalah, siapkan env khusus APK dengan URL absolut.

### 8.1 Plugin yang dipilih (hasil verifikasi)
- **`@kduma-autoid/capacitor-bluetooth-printer`** — plugin Capacitor untuk printer
  **SPP (Bluetooth Classic)**, cocok dengan EP58M. API:
  - `list()` → daftar device ter-pair (punya `type`: `classic` | `le` | `dual` | `unknown`).
  - `connect({ address })`, `disconnect()`
  - `print({ data })`, `connectAndPrint({ address, data })`
  - `data` berupa string (format—base64/teks—**dikonfirmasi saat implementasi**;
    kemungkinan base64 dari byte ESC/POS).
- Alternatif cadangan bila plugin di atas bermasalah: `Glavotaner/bluetooth-serial`
  (Cordova→Capacitor, kompatibel Android API 31) atau `capacitor-thermal-printer`.

### 8.2 Langkah setup Capacitor (di proyek FE)
1. Install: `@capacitor/core`, `@capacitor/cli`, `@capacitor/android`, dan plugin
   `@kduma-autoid/capacitor-bluetooth-printer`.
2. `npx cap init` → buat `capacitor.config.ts`:
   - `appId` (mis. `id.mahenz.pos`), `appName`.
   - `server.url = https://<domain-produksi>` dan `server.cleartext = false`
     (Opsi A: WebView membuka web produksi langsung).
3. `npx cap add android` → membuat folder `android/` (proyek Gradle).
4. `AndroidManifest.xml`: tambah izin Bluetooth
   - Android ≤ 11: `BLUETOOTH`, `BLUETOOTH_ADMIN`.
   - Android 12+ (API 31+): `BLUETOOTH_CONNECT`, `BLUETOOTH_SCAN` (izin runtime).
5. Signing key untuk rilis APK (keystore) — disimpan aman, dipakai tiap build rilis.

### 8.3 Lapisan "jembatan cetak" di FE (kode baru, terisolasi)
Buat modul baru, mis. `FE/src/features/sales/cashier/printNativeBt.ts`:
- `isNativeApp(): boolean` — deteksi apakah jalan di dalam APK Capacitor
  (`Capacitor.isNativePlatform()`), supaya pilihan hanya muncul di APK.
- `printViaNativeBt(data: ReceiptData, address?: string): Promise<void>`
  - Susun byte via `buildReceiptEscposBytes(data)`.
  - Encode ke format yang diminta plugin (mis. base64).
  - Panggil `list()` → pilih device (ingat pilihan terakhir), lalu
    `connectAndPrint({ address, data })`.
- Impor plugin dibuat **dynamic import** / dibungkus try-catch agar di browser
  (plugin tak ada) tidak error — cukup `isNativeApp()` false.

### 8.4 Integrasi ke modal cetak (`ReceiptPrint.tsx`)
- Tambah `PrintChoice` baru: `'native-bt'`.
- Tambah kartu di `PRINT_OPTIONS`:
  - Judul: "Printer Bluetooth (Aplikasi)" / "Cetak Bluetooth Langsung".
  - Keterangan: hanya untuk tablet/HP via aplikasi terpasang.
  - **Kartu hanya dirender bila `isNativeApp()` true** (difilter saat map, atau
    daftar opsi dihitung kondisional).
- Tambah handler `printNativeBt()` + cabang di `runPrintMethod`.
- **Tidak mengubah** 8 metode lama.

### 8.5 Build & distribusi
- Build web: `npm run build` (menghasilkan `dist/`). Untuk Opsi A, isi `dist`
  sebenarnya tidak dipakai WebView (WebView buka domain), tapi `cap sync` tetap
  butuh ada.
- `npx cap sync android` → sinkron plugin & config ke proyek Android.
- Build APK: via Android Studio (`Build > Generate Signed APK`) atau Gradle CLI.
- Distribusi: **sideload** `.apk` ke tablet (cukup untuk toko sendiri).

### 8.6 Uji di tablet (kriteria lulus)
1. Pair EP58M di Settings Android tablet.
2. Buka APK → aplikasi web tampil normal, login & transaksi jalan (lewat domain).
3. Di modal cetak, muncul kartu "Printer Bluetooth (Aplikasi)".
4. Pilih → device EP58M terdeteksi → struk tercetak **tajam & tidak terpotong**.
5. Di browser biasa (tanpa APK): kartu itu **tidak muncul**, 8 metode lama tetap jalan.

---

## 9. Pembagian Kerja / Urutan Implementasi (saat disetujui)
1. **FE**: tambah dependency Capacitor + plugin; buat `capacitor.config.ts` (Opsi A).
2. **FE**: modul `printNativeBt.ts` (deteksi native + fungsi cetak) — terisolasi.
3. **FE**: tambah kartu + handler di `ReceiptPrint.tsx` (muncul hanya di APK).
4. **Android**: `cap add android`, izin manifest, keystore, build APK.
5. **Uji** di tablet dengan EP58M (lihat 8.6).
6. **Docs**: tulis SOP build & install APK.
7. Verifikasi: FE `type-check` + `lint` tetap hijau; metode lama tak berubah.

---

## 10. Lingkungan & Cara Build (dikonfirmasi dalam diskusi)

### 10.1 Lingkungan terkonfirmasi
- **Domain WebView**: `https://mahenz.duckdns.org` — **HTTPS sudah jalan**. Dipakai
  sebagai `server.url` di `capacitor.config` (Opsi A).
- **Versi Android tablet**: beragam / tidak dipatok. Satu APK akan menangani Android
  lama maupun 12+ (deteksi versi untuk model izin Bluetooth runtime).
- **Repo**: aplikasi web sudah ada di **GitHub**. User terbiasa `git push` biasa,
  tapi **belum pernah** build APK.

### 10.2 Cara build APK: GitHub Actions (cloud) — dipilih
Alasan: laptop user **tidak kuat Android Studio**, jadi compile APK dilempar ke cloud.
- **Pembagian kerja:**
  - **Di laptop (ringan):** install Capacitor + plugin, bikin `capacitor.config.ts`,
    generate folder `android/`, tulis kode jembatan cetak. Tidak perlu Android Studio.
  - **Di GitHub Actions (berat):** compile `android/` menjadi file `.apk`. User unduh
    hasilnya dari halaman Actions (artifact) atau Release.
- **Keystore (tanda tangan APK):** dibuat **sekali** dengan satu perintah `keytool`
  (ikut JDK, ringan — tidak butuh Android Studio). File keystore + password disimpan
  sebagai **GitHub Secrets**, dipakai workflow untuk menandatangani APK rilis.
- **Workflow** (`.github/workflows/android-apk.yml`, dibuat saat implementasi) akan:
  1. checkout repo, setup Node + JDK + Android SDK (disediakan runner GitHub),
  2. `npm ci` + `npm run build` (FE),
  3. `npx cap sync android`,
  4. build APK via Gradle, tandatangani pakai keystore dari Secrets,
  5. unggah `.apk` sebagai artifact/Release untuk diunduh & di-sideload ke tablet.
- **Distribusi:** sideload `.apk` ke tablet (cukup untuk toko sendiri).

> Konsekuensi: setiap kali bagian **native** berubah (plugin/izin/config APK), build
> APK baru lewat Actions lalu install ulang di tablet. Perubahan **web** biasa tetap
> otomatis (WebView buka domain), tidak perlu rebuild APK.

---

## 10.3 Panduan Langkah Manual (yang dikerjakan user)

Bagian ini memisahkan dengan jelas **apa yang dikerjakan AI (lapis kode)** dan **apa
yang harus dikerjakan user**. Semua perintah ditulis lengkap agar bisa diikuti
tanpa pengalaman APK sebelumnya. (Perintah persis bisa sedikit menyesuaikan saat
implementasi; ini kerangka SOP.)

### Lapis yang dikerjakan AI (di dalam proyek)
- Install Capacitor + plugin, buat `capacitor.config.ts`, generate folder `android/`.
- Modul jembatan cetak FE (`printNativeBt.ts`) + kartu metode di `ReceiptPrint.tsx`.
- File workflow GitHub Actions (`.github/workflows/android-apk.yml`).
- Dokumentasi.

### Lapis yang dikerjakan USER (AI tidak bisa) — dengan panduan

**LANGKAH 1 — Install JDK (sekali, ringan).**
Dibutuhkan untuk perintah `keytool` (pembuat keystore). JDK jauh lebih ringan dari
Android Studio.
- Unduh & install **Temurin JDK 17** (adoptium.net) atau OpenJDK 17.
- Cek berhasil: buka PowerShell → `keytool -help` (muncul bantuan = sudah ada).

**LANGKAH 2 — Buat keystore (tanda tangan APK, sekali seumur hidup aplikasi).**
Jalankan di PowerShell (ganti nilai sesuai keinginan, CATAT passwordnya):
```
keytool -genkeypair -v -keystore pos-release.keystore -alias pos -keyalg RSA -keysize 2048 -validity 10000
```
- Akan ditanya password + beberapa data (nama, organisasi) — isi bebas.
- Hasil: file `pos-release.keystore`. **SIMPAN BAIK-BAIK** (kalau hilang, tidak bisa
  update APK yang sudah terpasang di tablet). JANGAN commit ke git.

**LANGKAH 3 — Ubah keystore jadi teks Base64 (untuk disimpan di GitHub Secrets).**
GitHub Secrets hanya menyimpan teks, jadi file keystore diubah ke base64 dulu:
```
[Convert]::ToBase64String([IO.File]::ReadAllBytes("pos-release.keystore")) | Set-Content pos-release.keystore.b64.txt
```
- Buka file `pos-release.keystore.b64.txt`, salin seluruh isinya untuk Langkah 4.

**LANGKAH 4 — Simpan Secrets di GitHub.**
Di repo GitHub → Settings → Secrets and variables → Actions → New repository secret.
Buat 4 secret ini:
| Nama Secret | Isi |
|-------------|-----|
| `ANDROID_KEYSTORE_B64` | isi file base64 dari Langkah 3 |
| `ANDROID_KEYSTORE_PASSWORD` | password keystore (Langkah 2) |
| `ANDROID_KEY_ALIAS` | `pos` (alias di Langkah 2) |
| `ANDROID_KEY_PASSWORD` | password key (biasanya sama dgn keystore) |

**LANGKAH 5 — Jalankan build APK (lewat GitHub Actions).**
Setelah AI menambahkan file workflow & semua kode, dan Anda `git push`:
- Buka repo GitHub → tab **Actions** → pilih workflow "Build Android APK".
- Klik **Run workflow** (atau otomatis jalan saat push, tergantung setelan).
- Tunggu selesai (beberapa menit). Bila hijau (sukses), buka run tsb → bagian
  **Artifacts** → unduh file `app-release.apk`.

**LANGKAH 6 — Pasang APK di tablet (sideload).**
- Pindahkan `app-release.apk` ke tablet (USB / Google Drive / WhatsApp ke diri sendiri).
- Di tablet, buka file → Install. Bila muncul "Install from unknown sources",
  izinkan untuk aplikasi pembuka (mis. File Manager/Chrome) lalu lanjut.

**LANGKAH 7 — Pair printer & uji.**
- Di tablet: Settings → Bluetooth → pair **EP58M** (sekali).
- Buka aplikasi POS (ikon APK) → login → transaksi → **Cetak** → pilih
  "Printer Bluetooth (Aplikasi)" → pilih EP58M → struk tercetak.
- Bila gagal connect di Android 12+, pastikan saat pertama kali aplikasi meminta
  **izin Bluetooth**, pilih **Izinkan**.

**Update ke depan:**
- Ganti tampilan/fitur **web** → cukup deploy web seperti biasa; tablet dapat
  otomatis (tidak perlu APK baru).
- Ganti **bagian native** (plugin/izin/domain) → ulangi Langkah 5–6 (build APK baru
  + install ulang). Keystore (Langkah 2) **tidak** dibuat ulang — pakai yang sama.

---

## 11. Ringkasan Keputusan (status terkini)
- [x] Pilih pendekatan: **bungkus web jadi APK (WebView + jembatan Bluetooth native)**.
- [x] Pilih **Opsi A** (WebView → URL domain; update web otomatis).
- [x] VPS + domain **tetap** dipakai.
- [x] APK **wajib** diinstall di tablet (sekali).
- [x] Bersifat aditif: website & 8 metode print lama **tidak terganggu**.
- [x] Framework: **Capacitor**. Plugin: **`@kduma-autoid/capacitor-bluetooth-printer`**
      (SPP, API `list/connect/print/connectAndPrint`) — terverifikasi cocok untuk EP58M.
- [x] Domain WebView: `https://mahenz.duckdns.org` (HTTPS aktif).
- [x] Android: satu APK menangani versi lama & 12+.
- [x] Build: **GitHub Actions** (cloud, gratis) — laptop tidak perlu Android Studio.
- [x] Rencana teknis detail ditulis (bagian 8, 9, 10).
- [ ] Konfirmasi format `data` plugin (base64 vs teks) — saat implementasi.
- [ ] Siap mulai implementasi? (menunggu aba-aba user)

---

## 11. Catatan Risiko
- **Plugin pihak ketiga**: `@kduma-autoid/...` relatif baru. Perlu verifikasi
  kompatibilitas Capacitor terbaru + Android target saat implementasi. Ada
  cadangan (`Glavotaner/bluetooth-serial`, `capacitor-thermal-printer`).
- **Izin Bluetooth Android 12+**: butuh runtime permission `BLUETOOTH_CONNECT`/
  `BLUETOOTH_SCAN`; harus di-handle di kode agar tidak gagal diam-diam.
- **Opsi A + path `/api` relatif**: kemungkinan besar aman (domain sama), tapi
  wajib diuji; siapkan fallback env URL absolut bila perlu.
- **Signing key APK**: harus disimpan aman & konsisten; kalau hilang, update APK
  di perangkat lama bisa bermasalah.
