package repo

import (
	"pos_api/domain/access/dto"
	"pos_api/domain/access/model"

	"gorm.io/gorm"
)

const getByRoleIDQuery = `
SELECT
    m.id         AS menu_id,
    m.key_name,
    m.label,
    m.parent_id,
    COALESCE(rma.can_view,   0) AS can_view,
    COALESCE(rma.can_create, 0) AS can_create,
    COALESCE(rma.can_edit,   0) AS can_edit,
    COALESCE(rma.can_delete, 0) AS can_delete
FROM menus m
LEFT JOIN role_menu_access rma ON rma.menu_id = m.id AND rma.role_id = ?
WHERE m.is_active = 1
ORDER BY m.order_index ASC, m.id ASC
`

const upsertAccessQuery = `
INSERT INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    can_view   = VALUES(can_view),
    can_create = VALUES(can_create),
    can_edit   = VALUES(can_edit),
    can_delete = VALUES(can_delete),
    updated_at = NOW()
`

const deleteRoleAccessQuery = `DELETE FROM role_menu_access WHERE role_id = ?`

func (r *accessRepo) GetByRoleID(roleID int) ([]*model.RoleMenuAccessItem, error) {
	var dataDB []*model.RoleMenuAccessItem
	err := r.db.Raw(getByRoleIDQuery, roleID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *accessRepo) SetRoleAccess(roleID int, accesses []dto.SetAccessItem) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(deleteRoleAccessQuery, roleID).Error
		if err != nil {
			return err
		}
		for _, a := range accesses {
			err := tx.Exec(upsertAccessQuery,
				roleID, a.MenuID, a.CanView, a.CanCreate, a.CanEdit, a.CanDelete,
			).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
