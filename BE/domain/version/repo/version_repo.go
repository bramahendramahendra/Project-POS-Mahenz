package repo

import "pos_api/domain/version/model"

const (
	getAllVersionsQuery    = `SELECT id, platform, version, download_url, release_notes, is_mandatory, is_latest, created_at FROM app_versions ORDER BY created_at DESC`
	getLatestVersionQuery = `SELECT version, download_url, release_notes, is_mandatory FROM app_versions WHERE platform = 'android' AND is_latest = 1 LIMIT 1`
	setAllNotLatestQuery  = `UPDATE app_versions SET is_latest = 0 WHERE platform = 'android'`
	createVersionQuery    = `INSERT INTO app_versions (platform, version, download_url, release_notes, is_mandatory, is_latest) VALUES ('android', ?, ?, ?, ?, 1)`
)

func (r *versionRepo) GetAll() ([]model.AppVersion, error) {
	var versions []model.AppVersion
	if err := r.db.Raw(getAllVersionsQuery).Scan(&versions).Error; err != nil {
		return nil, err
	}
	if versions == nil {
		versions = []model.AppVersion{}
	}
	return versions, nil
}

func (r *versionRepo) GetLatestAndroid() (*model.AppVersion, error) {
	v := &model.AppVersion{}
	err := r.db.Raw(getLatestVersionQuery).Scan(v).Error
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *versionRepo) SetAllNotLatest() error {
	return r.db.Exec(setAllNotLatestQuery).Error
}

func (r *versionRepo) CreateVersion(version, downloadURL, releaseNotes string, isMandatory bool) error {
	return r.db.Exec(createVersionQuery, version, downloadURL, releaseNotes, isMandatory).Error
}
