package request_helper

// NormalizePagination menormalkan page/limit dari request lalu menghitung offset.
// maxLimit=0 berarti tidak ada batas atas.
// Kalau limit <= 0 atau melebihi maxLimit (bila diset), limit direset ke defaultLimit.
func NormalizePagination(page, limit, defaultLimit, maxLimit int) (normPage, normLimit, offset int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || (maxLimit > 0 && limit > maxLimit) {
		limit = defaultLimit
	}
	return page, limit, (page - 1) * limit
}
