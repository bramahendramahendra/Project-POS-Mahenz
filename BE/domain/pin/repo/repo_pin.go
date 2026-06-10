package repo

const (
	getPinHashQuery   = `SELECT pin_hash FROM users WHERE id = ? LIMIT 1`
	setPinHashQuery   = `UPDATE users SET pin_hash = ? WHERE id = ?`
	clearPinHashQuery = `UPDATE users SET pin_hash = NULL WHERE id = ?`
)

func (r *pinRepo) GetPinHash(userID int) (string, error) {
	var result struct {
		PinHash *string `gorm:"column:pin_hash"`
	}
	res := r.db.Raw(getPinHashQuery, userID).Scan(&result)
	if res.Error != nil {
		return "", res.Error
	}
	if result.PinHash == nil {
		return "", nil
	}
	return *result.PinHash, nil
}

func (r *pinRepo) SetPinHash(userID int, pinHash string) error {
	return r.db.Exec(setPinHashQuery, pinHash, userID).Error
}

func (r *pinRepo) ClearPinHash(userID int) error {
	return r.db.Exec(clearPinHashQuery, userID).Error
}
