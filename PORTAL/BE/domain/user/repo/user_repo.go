package repo

const (
	deleteSessionByAuthUserIDQuery = `DELETE s FROM sessions s INNER JOIN users u ON u.id = s.user_id WHERE u.auth_user_id = ?`
	deleteProfileByAuthUserIDQuery = `DELETE FROM users WHERE auth_user_id = ?`
)

// DeleteSessionByAuthUserID mengeluarkan user dari Portal (sesi dihapus).
func (r *userRepo) DeleteSessionByAuthUserID(authUserID int) error {
	return r.db.Exec(deleteSessionByAuthUserIDQuery, authUserID).Error
}

// DeleteProfileByAuthUserID menghapus profil lokal (sesi ikut terhapus lewat foreign key).
func (r *userRepo) DeleteProfileByAuthUserID(authUserID int) error {
	return r.db.Exec(deleteProfileByAuthUserIDQuery, authUserID).Error
}
