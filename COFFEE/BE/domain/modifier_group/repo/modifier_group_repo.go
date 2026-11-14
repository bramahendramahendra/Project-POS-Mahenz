package repo

import (
	dto "coffee_api/domain/modifier_group/dto"
	model "coffee_api/domain/modifier_group/model"
	time_helper "coffee_api/helper/time"
)

const (
	groupSelect = `SELECT g.id, g.name, g.selection_type, g.is_required, g.is_active,
		(SELECT COUNT(*) FROM modifier_options o WHERE o.group_id = g.id) AS option_count,
		(SELECT COUNT(*) FROM menu_item_modifier_groups mg WHERE mg.group_id = g.id) AS menu_count
		FROM modifier_groups g`
	groupOrder            = ` ORDER BY g.name ASC`
	getGroupByIDQuery     = groupSelect + ` WHERE g.id = ? LIMIT 1`
	groupNameExistsQuery  = `SELECT COUNT(*) FROM modifier_groups WHERE name = ? AND id <> ?`
	createGroupQuery      = `INSERT INTO modifier_groups (name, selection_type, is_required, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	lastInsertIDQuery     = `SELECT LAST_INSERT_ID()`
	updateGroupQuery      = `UPDATE modifier_groups SET name = ?, selection_type = ?, is_required = ?, updated_at = ? WHERE id = ?`
	toggleGroupQuery      = `UPDATE modifier_groups SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteGroupQuery      = `DELETE FROM modifier_groups WHERE id = ?`
	optionsByGroupsQuery  = `SELECT id, group_id, name, price, is_default, is_active, sort_order FROM modifier_options WHERE group_id IN ?`
	optionsOrder          = ` ORDER BY group_id, sort_order, id`
	insertOptionQuery     = `INSERT INTO modifier_options (group_id, name, price, is_default, is_active, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	updateOptionQuery     = `UPDATE modifier_options SET name = ?, price = ?, is_default = ?, is_active = ?, sort_order = ?, updated_at = ? WHERE id = ? AND group_id = ?`
	deleteOptionsQuery    = `DELETE FROM modifier_options WHERE group_id = ?`
	deleteOptionsNotInQry = `DELETE FROM modifier_options WHERE group_id = ? AND id NOT IN ?`
)

func (r *modifierGroupRepo) GetAll(req *dto.GetAllRequest, onlyActive bool) ([]*model.ModifierGroup, error) {
	query := groupSelect + ` WHERE 1=1`
	var args []any
	if req.Search != "" {
		query += ` AND g.name LIKE ?`
		args = append(args, "%"+req.Search+"%")
	}
	if onlyActive {
		query += ` AND g.is_active = 1`
	} else if req.IsActive != nil {
		query += ` AND g.is_active = ?`
		args = append(args, *req.IsActive)
	}
	var data []*model.ModifierGroup
	if err := r.db.Raw(query+groupOrder, args...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *modifierGroupRepo) GetByID(id int) (*model.ModifierGroup, error) {
	var data model.ModifierGroup
	if err := r.db.Raw(getGroupByIDQuery, id).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *modifierGroupRepo) GetOptionsByGroupIDs(groupIDs []int, onlyActive bool) ([]*model.ModifierOption, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	query := optionsByGroupsQuery
	if onlyActive {
		query += ` AND is_active = 1`
	}
	var data []*model.ModifierOption
	if err := r.db.Raw(query+optionsOrder, groupIDs).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *modifierGroupRepo) NameExists(name string, excludeID int) (bool, error) {
	var count int64
	if err := r.db.Raw(groupNameExistsQuery, name, excludeID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *modifierGroupRepo) Create(req *dto.SaveRequest) (int64, error) {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(createGroupQuery, req.Name, req.SelectionType, req.IsRequired, now, now).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(lastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *modifierGroupRepo) Update(req *dto.SaveRequest) error {
	return r.db.Exec(updateGroupQuery, req.Name, req.SelectionType, req.IsRequired, time_helper.GetTimeNow(), req.ID).Error
}

func (r *modifierGroupRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleGroupQuery, time_helper.GetTimeNow(), id).Error
}

func (r *modifierGroupRepo) Delete(id int) error {
	if err := r.db.Exec(deleteOptionsQuery, id).Error; err != nil {
		return err
	}
	return r.db.Exec(deleteGroupQuery, id).Error
}

func (r *modifierGroupRepo) InsertOption(groupID int, o *dto.OptionInput, sortOrder int) error {
	now := time_helper.GetTimeNow()
	return r.db.Exec(insertOptionQuery, groupID, o.Name, o.Price, o.IsDefault, o.IsActive, sortOrder, now, now).Error
}

// UpdateOption: kepemilikan id pilihan sudah dicek service.
func (r *modifierGroupRepo) UpdateOption(groupID int, o *dto.OptionInput, sortOrder int) error {
	return r.db.Exec(updateOptionQuery, o.Name, o.Price, o.IsDefault, o.IsActive, sortOrder, time_helper.GetTimeNow(), o.ID, groupID).Error
}

func (r *modifierGroupRepo) DeleteOptionsNotIn(groupID int, keepIDs []int) error {
	if len(keepIDs) == 0 {
		return r.db.Exec(deleteOptionsQuery, groupID).Error
	}
	return r.db.Exec(deleteOptionsNotInQry, groupID, keepIDs).Error
}
