package segment

import (
	dining_table_handler "coffee_api/domain/dining_table/handler"
	dining_table_repo "coffee_api/domain/dining_table/repo"
	dining_table_service "coffee_api/domain/dining_table/service"
	menu_category_handler "coffee_api/domain/menu_category/handler"
	menu_category_repo "coffee_api/domain/menu_category/repo"
	menu_category_service "coffee_api/domain/menu_category/service"
	menu_item_handler "coffee_api/domain/menu_item/handler"
	menu_item_repo "coffee_api/domain/menu_item/repo"
	menu_item_service "coffee_api/domain/menu_item/service"
	menu_size_handler "coffee_api/domain/menu_size/handler"
	menu_size_repo "coffee_api/domain/menu_size/repo"
	menu_size_service "coffee_api/domain/menu_size/service"
	modifier_group_handler "coffee_api/domain/modifier_group/handler"
	modifier_group_repo "coffee_api/domain/modifier_group/repo"
	modifier_group_service "coffee_api/domain/modifier_group/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MenuItemRoutes: menu jualan Coffee (tahap 4.4, keputusan 53-55, 85-89).
// Izin per halaman: menu_jual.daftar, menu_jual.kategori, menu_jual.tambahan, menu_jual.ukuran,
// menu_jual.meja (tahap 4.6, keputusan 98).
// Pilihan (options) kategori, ukuran, dan kelompok tambahan dipakai form menu, jadi cukup
// izin lihat Daftar Menu.
func MenuItemRoutes(r *gin.RouterGroup) {
	svc := newAccessService()
	perm := func(menuKey string) func(string) gin.HandlerFunc {
		return func(action string) gin.HandlerFunc { return middleware.PermissionMiddleware(svc, menuKey, action) }
	}
	daftar, kategori, tambahan, ukuran := perm("menu_jual.daftar"), perm("menu_jual.kategori"), perm("menu_jual.tambahan"), perm("menu_jual.ukuran")
	meja := perm("menu_jual.meja")
	kasir := perm("penjualan.kasir")

	cat := menu_category_handler.NewMenuCategoryHandler(menu_category_service.NewMenuCategoryService(menu_category_repo.NewMenuCategoryRepo(pkgdatabase.DB)))
	cg := r.Group("/menu-categories")
	{
		cg.POST("/list", kategori("can_view"), cat.GetAll)
		cg.POST("/options", daftar("can_view"), cat.GetOptions)
		cg.POST("/create", kategori("can_create"), cat.Create)
		cg.POST("/update/:id", kategori("can_edit"), cat.Update)
		cg.POST("/toggle-status/:id", kategori("can_edit"), cat.ToggleStatus)
		cg.POST("/delete/:id", kategori("can_delete"), cat.Delete)
	}

	size := menu_size_handler.NewMenuSizeHandler(menu_size_service.NewMenuSizeService(menu_size_repo.NewMenuSizeRepo(pkgdatabase.DB)))
	sg := r.Group("/menu-sizes")
	{
		sg.POST("/list", ukuran("can_view"), size.GetAll)
		sg.POST("/options", daftar("can_view"), size.GetOptions)
		sg.POST("/create", ukuran("can_create"), size.Create)
		sg.POST("/update/:id", ukuran("can_edit"), size.Update)
		sg.POST("/toggle-status/:id", ukuran("can_edit"), size.ToggleStatus)
		sg.POST("/delete/:id", ukuran("can_delete"), size.Delete)
	}

	// Meja: dikelola admin/owner, kasir melihat dan memakai pilihan meja di kasir (4.7).
	tbl := dining_table_handler.NewDiningTableHandler(dining_table_service.NewDiningTableService(dining_table_repo.NewDiningTableRepo(pkgdatabase.DB)))
	tg := r.Group("/dining-tables")
	{
		tg.POST("/list", meja("can_view"), tbl.GetAll)
		tg.POST("/options", meja("can_view"), tbl.GetOptions)
		tg.POST("/create", meja("can_create"), tbl.Create)
		tg.POST("/update/:id", meja("can_edit"), tbl.Update)
		tg.POST("/toggle-status/:id", meja("can_edit"), tbl.ToggleStatus)
		tg.POST("/delete/:id", meja("can_delete"), tbl.Delete)
	}

	mod := modifier_group_handler.NewModifierGroupHandler(modifier_group_service.NewModifierGroupService(modifier_group_repo.NewModifierGroupRepo(pkgdatabase.DB)))
	mg := r.Group("/modifier-groups")
	{
		mg.POST("/list", tambahan("can_view"), mod.GetAll)
		mg.POST("/options", daftar("can_view"), mod.GetOptions)
		mg.POST("/create", tambahan("can_create"), mod.Create)
		mg.POST("/update/:id", tambahan("can_edit"), mod.Update)
		mg.POST("/toggle-status/:id", tambahan("can_edit"), mod.ToggleStatus)
		mg.POST("/delete/:id", tambahan("can_delete"), mod.Delete)
	}

	item := menu_item_handler.NewMenuItemHandler(menu_item_service.NewMenuItemService(menu_item_repo.NewMenuItemRepo(pkgdatabase.DB)))
	ig := r.Group("/menu-items")
	{
		// Katalog kasir: cukup izin lihat halaman Kasir (tahap 4.7b).
		ig.POST("/catalog", kasir("can_view"), item.Catalog)
		ig.POST("/list", daftar("can_view"), item.GetAll)
		ig.POST("/detail/:id", daftar("can_view"), item.GetByID)
		ig.POST("/create", daftar("can_create"), item.Create)
		ig.POST("/update/:id", daftar("can_edit"), item.Update)
		ig.POST("/toggle-status/:id", daftar("can_edit"), item.ToggleStatus)
		// Tandai habis/tersedia: kasir boleh (keputusan 88), cukup izin lihat.
		ig.POST("/sold-out/:id", daftar("can_view"), item.SetSoldOut)
		ig.POST("/delete/:id", daftar("can_delete"), item.Delete)
	}
}
