// Package logmask menyamarkan isi rahasia sebelum body permintaan/respons disimpan ke log.
// Toko meneruskan password (login) ke BE User dan mengeluarkan token sesi, sehingga isinya
// tidak boleh tersimpan apa adanya di log_requests (keputusan 74).
package logmask

import "encoding/json"

// Kunci JSON yang nilainya rahasia. redirect_url ikut disamarkan karena memuat kode SSO.
var sensitiveKeys = map[string]bool{
	"password": true, "old_password": true, "new_password": true,
	"app_key": true, "token": true, "refresh_token": true, "redirect_url": true,
}

// "code" dipakai dua arti: kode SSO (64 karakter, rahasia) dan kode status amplop respons
// ("00", "41", ...). Hanya nilai yang lebih panjang dari kode status yang disamarkan.
const maxStatusCodeLen = 8

func isSecret(key string, val any) bool {
	if sensitiveKeys[key] {
		return true
	}
	s, ok := val.(string)
	return key == "code" && ok && len(s) > maxStatusCodeLen
}

// Mask mengganti nilai kunci rahasia dengan "***" (bersarang di objek/array).
// Body yang bukan JSON dikembalikan apa adanya.
func Mask(body string) string {
	if body == "" {
		return body
	}
	var data any
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return body
	}
	masked, err := json.Marshal(maskValue(data))
	if err != nil {
		return body
	}
	return string(masked)
}

func maskValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if isSecret(k, val) {
				t[k] = "***"
			} else {
				t[k] = maskValue(val)
			}
		}
		return t
	case []any:
		for i := range t {
			t[i] = maskValue(t[i])
		}
		return t
	default:
		return v
	}
}
