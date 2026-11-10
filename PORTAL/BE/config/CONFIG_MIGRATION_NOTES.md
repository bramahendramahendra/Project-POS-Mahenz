# Config Migration Notes
## Migrasi: permen_api (BRI) → pos_api (POS Retail)

---

## 1. Struct / Tipe Data

### `Env`
| Field | Lama | Baru | Keterangan |
|---|---|---|---|
| `AppName` | Ada | Ada | |
| `AppAuthor` | Ada | Ada | Dipertahankan |
| `AppVersion` | Ada | Ada | |
| `AppHost` | Ada | Ada | Dipertahankan |
| `AppPort` | Ada | Ada | |
| `ReleaseMode` | Ada | Ada | |

### `DatabaseConfig`
| Field | Lama | Baru | Keterangan |
|---|---|---|---|
| `Type` | Ada (tanpa json tag) | Ada (dengan json tag) | |
| `Host` | Ada | Ada | |
| `Port` | Ada | Ada | |
| `User` | Ada | Ada | |
| `Password` | Ada | Ada | |
| `Database` | Ada | Ada | |
| `MaxOpenConns` | Ada | Ada | |
| `MaxIdleConns` | Ada | Ada | |
| `ConnMaxLifeTime` | Ada | Ada | |
| `ConnMaxIdleTime` | Ada | Ada | |
| `MaxLifetime` | Tidak ada | **Baru** | Menggantikan `MaxLifeTime` (typo fix di JSON key) |

### `GeneralConfig` → `generalCompat`
- Nama berubah menjadi `generalCompat` (private/tidak di-export)
- Tujuan: backward compatibility agar kode yang referensi `config.General.X` tidak perlu diubah
- Field yang dipertahankan: `SecretKey`, `TokenExpire`, `FormatTime`, `FormatDate`, `MaxTimeoutGracefulShutdown`

| Field | Lama | Baru | Alasan Dihapus |
|---|---|---|---|
| `SecretKey` | Ada | Ada | |
| `SecretCost` | Ada | **Dihapus** | Bcrypt cost — fitur hashing belum diimplementasi |
| `RateLimiterExp` | Ada | **Dihapus** | Fitur rate limiter tidak ada di project POS |
| `LogsPathprefix` | Ada | **Dihapus** | Diganti `LogPath` di struct `Config` |
| `CacheSeperator` | Ada | **Dihapus** | Fitur cache (Redis/Janitor) tidak dipakai |
| `TokenExpire` | Ada | Ada | |
| `FormatTime` | Ada | Ada (hardcoded) | Nilai tetap, dijadikan konstanta `DefaultFormatTime` |
| `FormatDate` | Ada | Ada (hardcoded) | Nilai tetap, dijadikan konstanta `DefaultFormatDate` |
| `Timezone` | Ada | **Dipindah** | Pindah ke `Config.Timezone` |
| `MaxTimeoutGracefulShutdown` | Ada | Ada (hardcoded=5) | Hardcoded, tidak perlu dari JSON |
| `Branch` | Ada | **Dihapus** | Spesifik BRI (kode cabang bank), tidak relevan POS |
| `Kostl` | Ada | **Dihapus** | Spesifik BRI (kode cost center SAP), tidak relevan POS |

### Struct Baru
- `Config` — struct utama gabungan semua konfigurasi aplikasi POS

### Struct Dihapus
| Struct | Alasan |
|---|---|
| `ESBConfig` | Integrasi ESB BRI dihapus |
| `ESBMonolithConfig` | Integrasi ESB Monolith BRI dihapus |
| `EMaterai` | Integrasi eMeterai dihapus |
| `BRIGateConfig` | Integrasi BRIGate dihapus |
| `RestClientConfig` | Tidak ada REST client eksternal |
| `BristarsConfig` | Layanan Bristars tidak dipakai |
| `MinioConfig` | Object storage tidak dipakai |

### Struct Dipertahankan (tidak dipakai aktif)
- `RedisConfig` — struct tetap ada untuk kompatibilitas, tapi tidak diisi dari JSON

---

## 2. Global Variables

| Variable | Lama | Baru | Keterangan |
|---|---|---|---|
| `ENV` | Ada | Ada | |
| `Cfg` | Tidak ada | **Baru** | Struct Config utama |
| `Db` | Ada | Ada | |
| `Redis` | Ada | **Dihapus** | Redis tidak dipakai |
| `General` | `*GeneralConfig` | `*generalCompat` | Backward compat |
| `Minio` | Ada | **Dihapus** | |
| `RestClient` | Ada | **Dihapus** | |
| `EsbConf` | Ada | **Dihapus** | |
| `EsbMonolithConf` | Ada | **Dihapus** | |
| `BRIGateConf` | Ada | **Dihapus** | |
| `RestMode` | Ada | **Dihapus** | |
| `Location` | Ada | Ada | |
| `FormatTime` | Ada | Ada | |
| `DefaultFormatTime` | Tidak ada | **Baru** | Konstanta `"2006-01-02 15:04:05"` |
| `DefaultFormatDate` | Tidak ada | **Baru** | Konstanta `"2006-01-02"` |

---

## 3. Environment Map

| Lama | Baru | Keputusan |
|---|---|---|
| `dev` | Ada | Dipertahankan |
| `prod` | Ada | Dipertahankan |
| `uat` | Ada | **Dipertahankan** (keputusan diskusi) |
| `local` | Ada | **Dipertahankan** (keputusan diskusi) |
| `bors` | Ada | **Dipertahankan** (keputusan diskusi) |

---

## 4. config_dev.json — Field yang DIHAPUS

| Field / Block | Alasan Dihapus |
|---|---|
| `TLSMode` | Tidak ada implementasi TLS |
| `SecretCost` | Fitur bcrypt hashing belum diimplementasi |
| `RateLimiterExp` | Fitur rate limiter tidak ada |
| `logsPathPrefix` | Diganti `LogPath` |
| `CacheSeperator` | Fitur cache tidak dipakai |
| `FormatTime` | Dijadikan konstanta hardcoded |
| `FormatDate` | Dijadikan konstanta hardcoded |
| `MaxTimeoutGracefulShutdown` | Hardcoded nilai 5 |
| `Branch` | Spesifik BRI |
| `Kostl` | Spesifik BRI |
| `BrigateBaseUrl` | Integrasi BRIGate dihapus |
| `EsbBaseUrl` | Integrasi ESB dihapus |
| `EsbMonolithBaseUrl` | Integrasi ESB Monolith dihapus |
| `RestClientTO` | Tidak ada REST client eksternal |
| Block `ESBConf` | Seluruh integrasi ESB dihapus |
| Block `BRIgateConf` + `EMaterai` | Seluruh integrasi BRIGate + eMeterai dihapus |
| Block `EsbMonolithConf` | Seluruh integrasi ESB Monolith dihapus |
| Block `RedisLocal` | Redis tidak dipakai |
| Block `Janitor` | In-memory cache tidak dipakai |
| Block `Minio` | Object storage tidak dipakai |
| `RestMode` | Toggle REST mode tidak relevan |

---

## 5. config_dev.json — Field yang DITAMBAH

| Field Baru | Alasan |
|---|---|
| `AppName` | Dipindah dari `.env` ke JSON |
| `AppVersion` | Dipindah dari `.env` ke JSON |
| `AppPort` | Dipindah dari `.env` ke JSON |
| `ReleaseMode` | Dipindah dari `.env` ke JSON |
| `Database.Type` | Eksplisit menyebutkan driver (`mysql`) |
| `Database.MaxLifetime` | Menggantikan `MaxLifeTime` (typo fix) |
| `RefreshTokenExpire` | Project POS butuh refresh token terpisah |
| `CorsAllowOrigins` | Dulu hardcoded di middleware, sekarang per environment |
| `LogPath` | Menggantikan `logsPathPrefix` |
| `MaxLogAge` | Berapa hari log disimpan sebelum dihapus |

---

## 6. Nilai yang Berubah di config_dev.json

| Key | Lama | Baru |
|---|---|---|
| `SecretKey` | `"My secret key"` | `"pos-retail-secret-key-dev"` |
| `TokenExpire` | `7200` (2 jam) | `28800` (8 jam) |
| `Database.Host` | `172.18.135.223` | `127.0.0.1` |
| `Database.Database` | `briperment` | `pos_retail_db` |
| `Database.Password` | `P@ssw0rd` | *(kosong)* |
| `Database.MaxOpenConns` | `100` | `50` |
| `Database.MaxIdleConns` | `5` | `10` |
| `Database.MaxLifeTime` | `3600` | `300` (key baru: `MaxLifetime`) |
