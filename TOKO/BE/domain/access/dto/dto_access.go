package dto

type GetByRoleIDRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type SetRoleAccessUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type RoleMenuAccessItem struct {
	MenuID    int    `json:"menu_id"`
	KeyName   string `json:"key_name"`
	Label     string `json:"label"`
	ParentID  *int   `json:"parent_id"`
	CanView   bool   `json:"can_view"`
	CanCreate bool   `json:"can_create"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
}

type SetAccessItem struct {
	MenuID    int  `json:"menu_id"    validate:"required,gt=0"`
	CanView   bool `json:"can_view"`
	CanCreate bool `json:"can_create"`
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}

type SetRoleAccessRequest struct {
	RoleID          int             `json:"-"`
	CurrentRoleName string          `json:"-"`
	Accesses        []SetAccessItem `json:"accesses" validate:"required,min=1,dive"`
}
