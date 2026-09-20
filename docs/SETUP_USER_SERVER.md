# Cara Membuat User Baru di Server (Linux)

Dokumen pendamping [DEPLOYMENT_PROD.md](./DEPLOYMENT_PROD.md). Server production sebaiknya tidak dioperasikan terus-menerus sebagai `root` — dokumen ini menjelaskan cara membuat user baru untuk keperluan deploy sehari-hari.

---

## 1. Kenapa Tidak Boleh Selalu Pakai Root?

`root` punya akses penuh tanpa batas ke seluruh sistem. Kalau sesi SSH Anda diretas, salah ketik perintah destruktif (`rm -rf`), atau ada aplikasi yang tereksploitasi saat berjalan sebagai root — dampaknya bisa merusak seluruh server. Dengan user terbatas + `sudo` (yang butuh password setiap kali dipakai untuk aksi administratif), risiko human error dan blast radius serangan jadi jauh lebih kecil.

---

## 2. Buat User Baru untuk Deploy Sehari-hari

Login sebagai `root` (karena saat ini itu satu-satunya akses Anda), lalu jalankan:

```bash
adduser deploy
```

Perintah ini interaktif — akan menanyakan:
- Password baru untuk user `deploy` (isi password kuat)
- Full Name, Room Number, dll → boleh dikosongkan (tekan Enter saja)
- Konfirmasi "Is the information correct? [Y/n]" → ketik `Y`

**Beri akses `sudo`** supaya user ini bisa menjalankan perintah administratif (install package, restart service, dll) tanpa harus login sebagai root:

```bash
usermod -aG sudo deploy
```

Verifikasi user berhasil dibuat dan masuk grup `sudo`:

```bash
id deploy
# Output harus memuat: groups=...,sudo
```

---

## 3. (Disarankan) Setup Login SSH via Key, Bukan Password

Supaya lebih aman dan tidak perlu ketik password setiap SSH:

**Di komputer lokal Anda** (bukan di server), generate key jika belum punya:

```bash
ssh-keygen -t ed25519 -C "email-anda@example.com"
```

Copy public key ke server (masih login sebagai `root` di server):

```bash
# Cara 1 — dari komputer lokal langsung ke server (paling mudah)
ssh-copy-id deploy@ALAMAT_IP_SERVER

# Cara 2 — manual kalau ssh-copy-id tidak tersedia
mkdir -p /home/deploy/.ssh
nano /home/deploy/.ssh/authorized_keys   # paste isi file ~/.ssh/id_ed25519.pub dari komputer lokal
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys
```

Test login dengan user baru:

```bash
ssh deploy@ALAMAT_IP_SERVER
```

Kalau berhasil masuk tanpa diminta password (atau minta passphrase key, bukan password server), setup berhasil.

---

## 4. Pindah ke User `deploy` untuk Kerja Sehari-hari

Mulai sekarang, gunakan `deploy` untuk:
- Clone repository
- Build backend (`go build`) dan frontend (`npm run build`)
- Jalankan perintah `sudo` saat perlu (misal `sudo systemctl restart pos-backend`, `sudo cp ... /var/www/...`)

```bash
su - deploy
# atau langsung SSH:
ssh deploy@ALAMAT_IP_SERVER
```

Semua contoh perintah `sudo ...` di [DEPLOYMENT_PROD.md](./DEPLOYMENT_PROD.md) sekarang dijalankan sebagai user `deploy`, bukan `root` langsung.

---

## 5. (Opsional, Lebih Aman) Nonaktifkan Login SSH Root

Setelah yakin user `deploy` bisa login dan pakai `sudo` dengan baik, sebaiknya matikan akses SSH langsung ke `root` — supaya penyerang tidak bisa mencoba brute-force akun `root` yang paling berharga.

```bash
sudo nano /etc/ssh/sshd_config
```

Cari baris `PermitRootLogin`, ubah jadi:

```
PermitRootLogin no
```

Simpan, lalu restart SSH:

```bash
sudo systemctl restart sshd
```

> ⚠️ **Sebelum menutup sesi terminal Anda saat ini**, buka satu sesi SSH baru dulu dan pastikan login `deploy` + `sudo` benar-benar berfungsi. Kalau langkah ini salah dan Anda terlanjur logout, bisa saja Anda terkunci dari server (kecuali punya akses console dari provider VPS).

---

## 6. Ringkasan Pembagian User

| User | Dipakai untuk |
|---|---|
| `root` | Hanya darurat / akses awal. Setelah setup selesai, sebaiknya tidak dipakai untuk SSH harian |
| `deploy` (baru dibuat) | Login SSH sehari-hari: git pull, build, `sudo systemctl restart ...`, copy file deploy |
| `www-data` (sudah ada di sistem, bawaan Nginx) | User yang benar-benar menjalankan proses backend Go (via systemd) dan Nginx — didefinisikan di `pos-backend.service` pada [DEPLOYMENT_PROD.md](./DEPLOYMENT_PROD.md#6-menjalankan-backend-sebagai-service-systemd) |

`deploy` tidak perlu diberi akses menjalankan aplikasi secara langsung — cukup hak untuk `sudo systemctl restart pos-backend` dan menulis ke folder deploy (`/opt/pos-mahenz`, `/var/www/pos-web`). Aplikasi tetap berjalan sebagai `www-data` untuk membatasi hak akses proses yang benar-benar terekspos ke internet.
