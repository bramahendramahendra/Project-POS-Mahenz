package logmask

import (
	"strings"
	"testing"
)

func TestMask_MasksSecretsAtAnyDepth(t *testing.T) {
	in := `{"username":"budi","password":"rahasia1","data":{"app_key":"pak_abc","list":[{"new_password":"x"}]}}`
	out := Mask(in)

	for _, secret := range []string{"rahasia1", "pak_abc", `"x"`} {
		if strings.Contains(out, secret) {
			t.Errorf("rahasia %q masih muncul di log: %s", secret, out)
		}
	}
	if !strings.Contains(out, `"username":"budi"`) {
		t.Errorf("field biasa ikut berubah: %s", out)
	}
}

func TestMask_SSOCodeMaskedButStatusCodeKept(t *testing.T) {
	ssoCode := strings.Repeat("ab", 32)
	in := `{"code":"00","data":{"redirect_url":"/toko/sso?code=` + ssoCode + `"}}`
	out := Mask(`{"code":"` + ssoCode + `"}`) + Mask(in)

	if strings.Contains(out, ssoCode) {
		t.Errorf("kode SSO masih muncul di log: %s", out)
	}
	if !strings.Contains(out, `"code":"00"`) {
		t.Errorf("kode status amplop respons ikut disamarkan: %s", out)
	}
}

func TestMask_NonJSONUnchanged(t *testing.T) {
	for _, in := range []string{"", "bukan json", "a=1&b=2"} {
		if got := Mask(in); got != in {
			t.Errorf("Mask(%q) = %q, seharusnya tidak berubah", in, got)
		}
	}
}
