package global_dto

// LogEntry adalah satu-satunya struct untuk semua level log.
// Gunakan helper/log.FromContext atau helper/log.FromBackground untuk membangunnya.
type LogEntry struct {
	Message    string
	Context    string // kategori luas: "Internal Error", "Auth", "Permission"
	Scope      string // operasi spesifik: "Bearer Token", "permcache"
	RequestId  string
	Method     string // HTTP method atau label operasi (e.g. "CACHE") untuk non-HTTP
	Endpoint   string // HTTP path atau operation path untuk non-HTTP
	StartTime  string
	EndTime    string
	Stacktrace string // diisi hanya untuk Warn/Error
	Data       any
}
