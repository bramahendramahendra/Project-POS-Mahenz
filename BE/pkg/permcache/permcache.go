package permcache

import (
	"fmt"
	"pos_api/pkg/janitor"
	"time"
)

const (
	ttl         = 10 * time.Minute
	cleanupEvery = 5 * time.Minute
	maxBytes    = 4 * 1024 * 1024 // 4 MB
	numShards   = 16
)

// Permission menyimpan flag akses sebuah menu untuk satu role.
type Permission struct {
	CanView   bool
	CanCreate bool
	CanEdit   bool
	CanDelete bool
}

var cache = janitor.NewCache(maxBytes, numShards, ttl, cleanupEvery)

func cacheKey(roleName, menuKey string) string {
	return fmt.Sprintf("perm:%s:%s", roleName, menuKey)
}

func rolePrefix(roleName string) string {
	return fmt.Sprintf("perm:%s:", roleName)
}

// Set menyimpan permission satu menu ke cache.
func Set(roleName, menuKey string, perm Permission) {
	cache.Set(cacheKey(roleName, menuKey), perm, int64(len(roleName)+len(menuKey)+64), janitor.DefaultExpiration)
}

// Get mengambil permission dari cache. ok=false jika tidak ada / sudah expired.
func Get(roleName, menuKey string) (Permission, bool) {
	val, ok := cache.Get(cacheKey(roleName, menuKey))
	if !ok {
		return Permission{}, false
	}
	perm, ok := val.(Permission)
	return perm, ok
}

// InvalidateRole menghapus semua cache permission untuk satu role.
// Dipanggil saat SetRoleAccess berhasil.
func InvalidateRole(roleName string, menuKeys []string) {
	prefix := rolePrefix(roleName)
	_ = prefix
	for _, key := range menuKeys {
		cache.Delete(cacheKey(roleName, key))
	}
}
