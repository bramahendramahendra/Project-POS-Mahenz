package repo

import model "coffee_api/domain/cash_drawer/model"

// Aturan kas dan batas hari Coffee (tahap 4.9, keputusan 110-112).

const (
	countOpenBillsOnDateQuery    = `SELECT COUNT(*) FROM transactions WHERE status = 'open' AND DATE(transaction_date) = ?`
	countOtherOpenDrawersOnDateQ = `SELECT COUNT(*) FROM cash_drawer WHERE status = 'open' AND DATE(open_time) = ? AND id <> ? AND (is_backdate = 0 OR is_backdate IS NULL)`
)

// CountOpenBillsOnDate: jumlah pesanan terbuka yang dibuat pada tanggal (YYYY-MM-DD) itu.
func (r *cashDrawerRepo) CountOpenBillsOnDate(date string) (int64, error) {
	var n int64
	err := r.db.Raw(countOpenBillsOnDateQuery, date).Scan(&n).Error
	return n, err
}

// CountOtherOpenDrawersOnDate: jumlah kas lain yang masih terbuka dan dibuka pada tanggal itu.
func (r *cashDrawerRepo) CountOtherOpenDrawersOnDate(date string, excludeID int) (int64, error) {
	var n int64
	err := r.db.Raw(countOtherOpenDrawersOnDateQ, date, excludeID).Scan(&n).Error
	return n, err
}

const findPendingDrawerQuery = `SELECT id, user_id, open_time FROM cash_drawer
	WHERE status = 'open' AND DATE(open_time) = ? AND (is_backdate = 0 OR is_backdate IS NULL)
	ORDER BY (user_id = ?) DESC, (user_id = ?) DESC, open_time ASC LIMIT 1`

// FindDrawerForDate: kas yang masih terbuka pada tanggal itu untuk menampung pembayaran pesanan tertunda
// (keputusan 113): utamakan kas pembuat, lalu kas kasir penerima, lalu kas lain hari itu. nil bila tidak ada.
func (r *cashDrawerRepo) FindDrawerForDate(date string, creatorID, payerID int) (*model.CashDrawer, error) {
	var cd model.CashDrawer
	res := r.db.Raw(findPendingDrawerQuery, date, creatorID, payerID).Scan(&cd)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &cd, nil
}
