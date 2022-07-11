package router

import (
	"fmt"
	"github.com/NubeIO/edge/controller"
	dbase "github.com/NubeIO/edge/database"
	"github.com/NubeIO/edge/pkg/config"
	"github.com/NubeIO/edge/pkg/logger"
	"github.com/spf13/viper"
	"io"

	"github.com/NubeIO/edge/service/apps/installer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"time"
)

func Setup(db *gorm.DB) *gin.Engine {
	engine := gin.New()

	// Set gin access logs
	fileLocation := fmt.Sprintf("%s/edge.access.log", config.Config.GetAbsDataDir())
	f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		logger.Logger.Errorf("Failed to create access log file: %v", err)
	} else {
		gin.SetMode(viper.GetString("gin.loglevel"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders: []string{
			"X-FLOW-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host",
		},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowAllOrigins:        true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	}))

	appDB := &dbase.DB{
		DB: db,
	}

	install := installer.New(&installer.Installer{
		DB: appDB,
	})

	api := controller.Controller{DB: appDB, Installer: install}

	apiRoutes := engine.Group("/api")

	store := apiRoutes.Group("/stores")
	{
		store.GET("/", api.GetAppStores)
		store.POST("/", api.CreateAppStore)
		store.GET("/:uuid", api.GetAppStore)
		store.PATCH("/:uuid", api.UpdateAppStore)
		store.DELETE("/:uuid", api.DeleteAppStore)
		store.DELETE("/drop", api.DropAppStores)
	}

	app := apiRoutes.Group("/apps")
	{
		app.GET("/", api.GetApps)
		app.POST("/", api.InstallApp)
		app.GET("/:uuid", api.GetApp)
		app.PATCH("/:uuid", api.UpdateApp)
		app.DELETE("/", api.UnInstallApp)
		app.DELETE("/drop", api.DropApps)

		// stats
		app.POST("/progress/install", api.GetInstallProgress)
		app.POST("/progress/uninstall", api.GetUnInstallProgress)
		app.POST("/stats", api.AppStats)
	}

	device := apiRoutes.Group("/device")
	{
		device.GET("/", api.GetDeviceInfo)
		device.POST("/", api.AddDeviceInfo)
		device.PATCH("/", api.UpdateDeviceInfo)
	}

	system := apiRoutes.Group("/system")
	{
		system.GET("/ping", api.Ping)
		system.GET("/time", api.HostTime)
		system.GET("/product", api.GetProduct)
		system.POST("/scanner", api.RunScanner)
	}

	networking := apiRoutes.Group("/networking")
	{
		networking.GET("/networks", api.Networking)
		networking.GET("/interfaces", api.GetInterfacesNames)
		networking.GET("/internet", api.InternetIP)
		networking.GET("/update/schema", api.GetIpSchema)
		networking.POST("/update/dhcp", api.SetDHCP)
		networking.POST("/update/static", api.SetStaticIP)
	}

	files := apiRoutes.Group("/files")
	{
		files.GET("/list", api.ListFiles)
		files.POST("/rename", api.RenameFile)
		files.POST("/copy", api.CopyFile)
		files.POST("/move", api.MoveFile)
		files.POST("/upload", api.UploadFile)
		files.POST("/download", api.DownloadFile)
		files.DELETE("/delete", api.DeleteFile)
	}

	dirs := apiRoutes.Group("/dirs")
	{
		dirs.POST("/create", api.CreateDir)
		dirs.POST("/copy", api.CopyDir)
		dirs.DELETE("/delete", api.DeleteDir)
	}

	zip := apiRoutes.Group("/zip")
	{
		zip.POST("/unzip", api.Unzip)
		zip.POST("/zip", api.ZipDir)
	}

	systemctl := apiRoutes.Group("/systemctl")
	{
		systemctl.POST("/", api.SystemCtl)
	}

	return engine
}
