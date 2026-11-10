package model

type AppVersion struct {
	ID           int    `gorm:"column:id"`
	Platform     string `gorm:"column:platform"`
	Version      string `gorm:"column:version"`
	DownloadURL  string `gorm:"column:download_url"`
	ReleaseNotes string `gorm:"column:release_notes"`
	IsMandatory  bool   `gorm:"column:is_mandatory"`
	IsLatest     bool   `gorm:"column:is_latest"`
	CreatedAt    string `gorm:"column:created_at"`
}
