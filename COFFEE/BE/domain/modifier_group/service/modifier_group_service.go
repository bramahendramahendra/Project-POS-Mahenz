package service

import (
	"strings"

	dto "coffee_api/domain/modifier_group/dto"
	model "coffee_api/domain/modifier_group/model"
	"coffee_api/errors"

	"gorm.io/gorm"
)

func (s *modifierGroupService) GetAll(req *dto.GetAllRequest) ([]dto.GroupResponse, error) {
	return s.list(req, false)
}

func (s *modifierGroupService) GetOptions() ([]dto.GroupResponse, error) {
	return s.list(&dto.GetAllRequest{}, true)
}

func (s *modifierGroupService) list(req *dto.GetAllRequest, onlyActive bool) ([]dto.GroupResponse, error) {
	groups, err := s.repo.GetAll(req, onlyActive)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(groups))
	for _, g := range groups {
		ids = append(ids, g.ID)
	}
	options, err := s.repo.GetOptionsByGroupIDs(ids, onlyActive)
	if err != nil {
		return nil, err
	}
	byGroup := map[int][]dto.OptionResponse{}
	for _, o := range options {
		byGroup[o.GroupID] = append(byGroup[o.GroupID], toOptionResponse(o))
	}
	data := make([]dto.GroupResponse, 0, len(groups))
	for _, g := range groups {
		data = append(data, toGroupResponse(g, byGroup[g.ID]))
	}
	return data, nil
}

func (s *modifierGroupService) GetByID(id int) (*dto.GroupResponse, error) {
	g, err := s.mustGet(id)
	if err != nil {
		return nil, err
	}
	options, err := s.repo.GetOptionsByGroupIDs([]int{id}, false)
	if err != nil {
		return nil, err
	}
	res := make([]dto.OptionResponse, 0, len(options))
	for _, o := range options {
		res = append(res, toOptionResponse(o))
	}
	out := toGroupResponse(g, res)
	return &out, nil
}

func (s *modifierGroupService) Create(req *dto.SaveRequest) (*dto.GroupResponse, error) {
	if err := s.normalizeAndValidate(req, nil); err != nil {
		return nil, err
	}
	var newID int64
	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		id, err := r.Create(req)
		if err != nil {
			return err
		}
		newID = id
		for i := range req.Options {
			if err := r.InsertOption(int(id), &req.Options[i], i); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(int(newID))
}

// Update: pilihan dengan id dipertahankan dan diperbarui, id 0 ditambah, yang tidak dikirim dihapus.
func (s *modifierGroupService) Update(req *dto.SaveRequest) error {
	if _, err := s.mustGet(req.ID); err != nil {
		return err
	}
	existing, err := s.repo.GetOptionsByGroupIDs([]int{req.ID}, false)
	if err != nil {
		return err
	}
	owned := map[int]bool{}
	for _, o := range existing {
		owned[o.ID] = true
	}
	if err := s.normalizeAndValidate(req, owned); err != nil {
		return err
	}

	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		if err := r.Update(req); err != nil {
			return err
		}
		keep := []int{}
		for i := range req.Options {
			if req.Options[i].ID > 0 {
				keep = append(keep, req.Options[i].ID)
			}
		}
		if err := r.DeleteOptionsNotIn(req.ID, keep); err != nil {
			return err
		}
		for i := range req.Options {
			o := &req.Options[i]
			if o.ID > 0 {
				if err := r.UpdateOption(req.ID, o, i); err != nil {
					return err
				}
			} else if err := r.InsertOption(req.ID, o, i); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *modifierGroupService) ToggleStatus(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// Delete: kelompok yang masih dipasang ke menu tidak bisa dihapus (lepas dari menu dulu atau nonaktifkan).
func (s *modifierGroupService) Delete(id int) error {
	g, err := s.mustGet(id)
	if err != nil {
		return err
	}
	if g.MenuCount > 0 {
		return &errors.BadRequestError{Message: "Kelompok tambahan masih dipasang di menu. Lepas dari menu dulu atau nonaktifkan"}
	}
	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		return s.repo.WithTx(tx).Delete(id)
	})
}

// normalizeAndValidate: nama unik, nama pilihan unik dalam kelompok, paling banyak satu pilihan
// bawaan dan harus aktif, minimal satu pilihan aktif, id pilihan harus milik kelompok ini.
func (s *modifierGroupService) normalizeAndValidate(req *dto.SaveRequest, owned map[int]bool) error {
	req.Name = strings.TrimSpace(req.Name)
	exists, err := s.repo.NameExists(req.Name, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama kelompok tambahan sudah digunakan"}
	}

	seen := map[string]bool{}
	defaults, active := 0, 0
	for i := range req.Options {
		o := &req.Options[i]
		o.Name = strings.TrimSpace(o.Name)
		key := strings.ToLower(o.Name)
		if seen[key] {
			return &errors.BadRequestError{Message: "Nama pilihan \"" + o.Name + "\" ganda dalam satu kelompok"}
		}
		seen[key] = true
		if o.ID > 0 && !owned[o.ID] {
			return &errors.BadRequestError{Message: "Pilihan tidak ditemukan di kelompok ini"}
		}
		if o.IsDefault {
			if !o.IsActive {
				return &errors.BadRequestError{Message: "Pilihan bawaan harus aktif"}
			}
			defaults++
		}
		if o.IsActive {
			active++
		}
	}
	if defaults > 1 {
		return &errors.BadRequestError{Message: "Pilihan bawaan paling banyak satu"}
	}
	if active == 0 {
		return &errors.BadRequestError{Message: "Minimal satu pilihan aktif"}
	}
	return nil
}

func (s *modifierGroupService) mustGet(id int) (*model.ModifierGroup, error) {
	g, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, &errors.NotFoundError{Message: "Kelompok tambahan tidak ditemukan"}
	}
	return g, nil
}

func toOptionResponse(o *model.ModifierOption) dto.OptionResponse {
	return dto.OptionResponse{ID: o.ID, Name: o.Name, Price: o.Price, IsDefault: o.IsDefault, IsActive: o.IsActive}
}

func toGroupResponse(g *model.ModifierGroup, options []dto.OptionResponse) dto.GroupResponse {
	if options == nil {
		options = []dto.OptionResponse{}
	}
	return dto.GroupResponse{
		ID:            g.ID,
		Name:          g.Name,
		SelectionType: g.SelectionType,
		IsRequired:    g.IsRequired,
		IsActive:      g.IsActive,
		MenuCount:     g.MenuCount,
		Options:       options,
	}
}
