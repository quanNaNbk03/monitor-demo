package router

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/vm"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/quanNaNbk03/monitor-demo/internal/server/handler"
	"github.com/quanNaNbk03/monitor-demo/internal/server/middleware"
	"github.com/quanNaNbk03/monitor-demo/internal/server/repository"
	"github.com/quanNaNbk03/monitor-demo/internal/server/service"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/health"
	org "github.com/quanNaNbk03/monitor-demo/pkg/api/organization"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/zone"
	zoneclient "github.com/quanNaNbk03/monitor-demo/pkg/api/zoneclient"

	"git.ocn.com.vn/ocn/common/appcfg"
	"git.ocn.com.vn/ocn/common/storage"
)

// Setup creates and configures the Gin router with all routes and middleware
func Setup(mysqlInstance *storage.MySQL) *gin.Engine {
	// Set Gin mode based on configuration
	var ginMode string
	if appcfg.IsProductionMode() {
		ginMode = gin.ReleaseMode
	} else {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// Create a Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: middleware.LogFormatterMiddleware,
		SkipPaths: viper.GetStringSlice("log.accessLog.ignorePaths"),
	}))

	if viper.GetBool("http.swagger.isEnabled") {
		// Swagger endpoints - serve OpenAPI spec at a different path to avoid conflicts
		router.GET("/openapi.json", handler.GetSwaggerJSON)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.json")))
	}
	// Init Repository
	zoneRepo := repository.NewZoneRepository(mysqlInstance.DB)
	orgRepo := repository.NewOrgRepository(mysqlInstance.DB)
	prometheusRepo, prometheusRepoErr := repository.NewPrometheusClient()
	if prometheusRepoErr != nil {
		logger.Fatalf("Error creating prometheus client: %v", prometheusRepoErr)
	}
	victoriaMetricsRepo, victoriaMetricsRepoErr := repository.NewVictoriaMetricsClient()
	if victoriaMetricsRepoErr != nil {
		logger.Fatalf("Error creating victoriametrics client: %v", victoriaMetricsRepoErr)
	}
	// Init services
	zoneService := service.NewZoneService(zoneRepo, orgRepo)
	orgService := service.NewOrgService(orgRepo)
	monitorVMService := service.NewMonitorVMService(prometheusRepo)
	victoriaMetricsService := service.NewVictoriaMetricsService(victoriaMetricsRepo)
	// Create separate handlers for each domain
	healthHandler := handler.NewHealthHandler()

	// Register routes with the API version prefix
	apiGroup := router.Group("")

	appInstance := viper.GetString("app.instance")

	if appInstance == "admin" {
		zoneHandler := handler.NewZoneHandler(zoneService)
		orgHandler := handler.NewOrgHandler(orgService)
		vmHandler := handler.NewVMHandler(monitorVMService, victoriaMetricsService)

		zone.RegisterHandlers(apiGroup, zoneHandler)
		org.RegisterHandlers(apiGroup, orgHandler)
		vm.RegisterHandlers(apiGroup, vmHandler)
	} else if appInstance == "client" {
		zoneClientHandler := handler.NewZoneClientHandler(zoneService)

		zoneclient.RegisterHandlers(apiGroup, zoneClientHandler)
	}

	// Register each handler to its routes
	health.RegisterHandlers(apiGroup, healthHandler)

	return router
}
