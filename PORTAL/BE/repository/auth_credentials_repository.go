package repository

import (
	"portal_api/pkg/authclient"

	"gorm.io/gorm"
)

const (
	getAuthCredentialsQuery = `SELECT app_id, app_key FROM auth_credentials WHERE id = 1 LIMIT 1`
)

// authCredentialsRepository membaca kunci aplikasi Portal untuk BE User (tabel auth_credentials).
// Dibaca setiap panggilan agar kunci baru langsung berlaku tanpa restart.
type authCredentialsRepository struct {
	Db *gorm.DB
}

func NewAuthCredentialsRepository(db *gorm.DB) *authCredentialsRepository {
	return &authCredentialsRepository{Db: db}
}

func (r *authCredentialsRepository) GetAuthCredentials() (*authclient.Credentials, error) {
	var row struct {
		AppID  string `gorm:"column:app_id"`
		AppKey string `gorm:"column:app_key"`
	}
	if err := r.Db.Raw(getAuthCredentialsQuery).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.AppKey == "" {
		return nil, nil
	}
	return &authclient.Credentials{AppID: row.AppID, AppKey: row.AppKey}, nil
}
