package model

type RoleMenuAccessItem struct {
	MenuID    int    `gorm:"column:menu_id"`
	KeyName   string `gorm:"column:key_name"`
	Label     string `gorm:"column:label"`
	ParentID  *int   `gorm:"column:parent_id"`
	CanView   bool   `gorm:"column:can_view"`
	CanCreate bool   `gorm:"column:can_create"`
	CanEdit   bool   `gorm:"column:can_edit"`
	CanDelete bool   `gorm:"column:can_delete"`
}
