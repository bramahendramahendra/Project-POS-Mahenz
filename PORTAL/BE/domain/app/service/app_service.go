package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	dto "portal_api/domain/app/dto"
	model "portal_api/domain/app/model"
	"portal_api/errors"
)

// code aplikasi dipakai sebagai app_id pada SSO: huruf kecil, angka, dan strip.
var codePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,49}$`)

func (s *appService) GetMyApps(userID int) ([]*dto.MyAppResponse, error) {
	apps, err := s.repo.GetMyApps(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.MyAppResponse, 0, len(apps))
	for _, a := range apps {
		result = append(result, &dto.MyAppResponse{
			ID:          a.ID,
			Code:        a.Code,
			Name:        a.Name,
			Description: a.Description,
			Path:        a.Path,
			Icon:        a.Icon,
		})
	}
	return result, nil
}

func (s *appService) GetAll() ([]*dto.AppAdminResponse, error) {
	apps, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AppAdminResponse, 0, len(apps))
	for _, a := range apps {
		result = append(result, toAdminResponse(a))
	}
	return result, nil
}

func (s *appService) Create(req *dto.CreateRequest) (*dto.CreateAppResponse, error) {
	code := strings.TrimSpace(req.Code)
	if !codePattern.MatchString(code) {
		return nil, &errors.BadRequestError{Message: "Kode aplikasi hanya boleh huruf kecil, angka, dan strip (diawali huruf), 2-50 karakter"}
	}
	if err := validatePath(req.Path); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Kode aplikasi sudah digunakan"}
	}

	key, err := generateAppKey()
	if err != nil {
		return nil, err
	}

	newID, err := s.repo.Create(&model.App{
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Path:        strings.TrimSpace(req.Path),
		Icon:        iconOrDefault(req.Icon),
		SortOrder:   req.SortOrder,
	}, hashHex(key))
	if err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(int(newID))
	if err != nil || created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data aplikasi baru"}
	}

	return &dto.CreateAppResponse{App: *toAdminResponse(created), AppKey: key}, nil
}

func (s *appService) Update(id int, req *dto.UpdateRequest) error {
	if err := s.mustExist(id); err != nil {
		return err
	}
	if err := validatePath(req.Path); err != nil {
		return err
	}

	return s.repo.Update(id, &model.App{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Path:        strings.TrimSpace(req.Path),
		Icon:        iconOrDefault(req.Icon),
		SortOrder:   req.SortOrder,
	})
}

func (s *appService) ToggleStatus(id int) error {
	if err := s.mustExist(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// RotateKey membuat kunci baru. Kunci lama langsung tidak berlaku.
func (s *appService) RotateKey(id int) (*dto.RotateKeyResponse, error) {
	if err := s.mustExist(id); err != nil {
		return nil, err
	}

	key, err := generateAppKey()
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetKeyHash(id, hashHex(key)); err != nil {
		return nil, err
	}
	return &dto.RotateKeyResponse{AppKey: key}, nil
}

func (s *appService) GetUserAccess(userID int) ([]*dto.UserAccessItem, error) {
	if err := s.mustUserExist(userID); err != nil {
		return nil, err
	}

	rows, err := s.repo.GetUserAccess(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.UserAccessItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, &dto.UserAccessItem{
			AppID:    r.AppID,
			Code:     r.Code,
			Name:     r.Name,
			IsActive: r.IsActive,
			Granted:  r.Granted,
		})
	}
	return result, nil
}

func (s *appService) SetUserAccess(userID int, req *dto.SetUserAccessRequest) error {
	if err := s.mustUserExist(userID); err != nil {
		return err
	}

	// Buang duplikat agar tidak melanggar primary key.
	seen := make(map[int]bool, len(req.AppIDs))
	unique := make([]int, 0, len(req.AppIDs))
	for _, id := range req.AppIDs {
		if !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}

	return s.repo.SetUserAccess(userID, unique)
}

func (s *appService) mustExist(id int) error {
	app, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if app == nil {
		return &errors.NotFoundError{Message: "Aplikasi tidak ditemukan"}
	}
	return nil
}

func (s *appService) mustUserExist(userID int) error {
	ok, err := s.repo.UserExists(userID)
	if err != nil {
		return err
	}
	if !ok {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	return nil
}

// validatePath: path di domain yang sama ("/toko/") atau URL penuh (http/https).
func validatePath(path string) error {
	p := strings.TrimSpace(path)
	if strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "//") {
		return nil
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return nil
	}
	return &errors.BadRequestError{Message: "Alamat aplikasi harus berupa path diawali '/' (mis. /toko/) atau URL http/https"}
}

func iconOrDefault(icon string) string {
	icon = strings.TrimSpace(icon)
	if icon == "" {
		return "app"
	}
	return icon
}

// generateAppKey: 32 byte acak dalam hex, berawalan "pak_" agar mudah dikenali.
func generateAppKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "pak_" + hex.EncodeToString(b), nil
}

func hashHex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func toAdminResponse(a *model.AppAdmin) *dto.AppAdminResponse {
	return &dto.AppAdminResponse{
		ID:          a.ID,
		Code:        a.Code,
		Name:        a.Name,
		Description: a.Description,
		Path:        a.Path,
		Icon:        a.Icon,
		SortOrder:   a.SortOrder,
		IsActive:    a.IsActive,
		HasKey:      a.HasKey,
		UserCount:   a.UserCount,
		CreatedAt:   a.CreatedAt,
	}
}
