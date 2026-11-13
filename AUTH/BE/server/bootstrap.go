package bootstrap

import (
	"auth_api/database"
	error_helper "auth_api/helper/error"
	pkgdatabase "auth_api/pkg/database"
	"auth_api/pkg/logger"
	"auth_api/routes"
	_ "auth_api/validation"

	"auth_api/config"

	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
)

func Initialized() *gin.Engine {
	initializedLogger()
	InitializedDB()

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.Default()

	trustedProxies := []string{
		"127.0.0.1",
		"::1",
	}

	if err := ginEngine.SetTrustedProxies(trustedProxies); err != nil {
		errData := error_helper.SetError(nil, "Proxy Configuration", err.Error(), error_helper.GetStackTrace(1), nil)
		panic(errData)
	}

	routes.Router(ginEngine)
	return ginEngine
}

func InitializedDB() {
	dbManager := pkgdatabase.New()
	pkgdatabase.DbManager = dbManager
	err := dbManager.Register(config.Db.Database, config.Db)
	if err != nil {
		errData := error_helper.SetError(nil, "DB Initialization", err.Error(), error_helper.GetStackTrace(1), nil)
		panic(errData)
	}

	pkgdatabase.DB = dbManager.GetDatabase(config.Db.Database)

	if err := database.RunMigrations(pkgdatabase.DB); err != nil {
		errData := error_helper.SetError(nil, "Migration", err.Error(), error_helper.GetStackTrace(1), nil)
		panic(errData)
	}
}

func initializedLogger() {
	logger.Log = logger.New()
	go logger.StartRotationWatcher()
}
