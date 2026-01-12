package api

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/services-manager/internal/config"
	"github.com/services-manager/internal/database"
	"github.com/services-manager/internal/services"
)

// Server representa o servidor HTTP
type Server struct {
	router            *gin.Engine
	cfg               *config.Config
	db                *sql.DB
	clientRepo        *database.ClientRepository
	serviceRepo       *database.ServiceRepository
	installationRepo *database.InstallationRepository
	nginxManager      *services.NginxManager
}

// NewServer cria uma nova inst?ncia do servidor
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	nginxConfigDir := os.Getenv("NGINX_CONFIG_DIR")
	if nginxConfigDir == "" {
		nginxConfigDir = "./nginx/sites-enabled"
	}

	nginxManager := services.NewNginxManager(
		nginxConfigDir,
		"./nginx/templates",
		"nginx",
	)

	server := &Server{
		cfg:        cfg,
		db:         db,
		clientRepo: database.NewClientRepository(db),
		serviceRepo: database.NewServiceRepository(db),
		installationRepo: database.NewInstallationRepository(db),
		nginxManager: nginxManager,
	}

	server.setupRouter()
	return server
}

// setupRouter configura as rotas da API
func (s *Server) setupRouter() {
	if s.cfg.API.Port == "8080" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.Default()

	// Middleware
	s.router.Use(gin.Recovery())
	s.router.Use(corsMiddleware())

	// Health check
	s.router.GET("/health", s.healthCheck)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Clientes
		v1.POST("/clients", s.createClient)
		v1.GET("/clients", s.listClients)
		v1.GET("/clients/:id", s.getClient)
		v1.PUT("/clients/:id", s.updateClient)
		v1.DELETE("/clients/:id", s.deleteClient)

		// Servi?os
		v1.POST("/services", s.createService)
		v1.GET("/services", s.listServices)
		v1.GET("/services/:id", s.getService)
		v1.PUT("/services/:id", s.updateService)
		v1.DELETE("/services/:id", s.deleteService)
		v1.PUT("/services/:id/resources", s.updateServiceResources)
		v1.GET("/services/:id/resources/usage", s.getServiceResourceUsage)
		v1.POST("/services/:id/start", s.startService)
		v1.POST("/services/:id/stop", s.stopService)
		v1.POST("/services/:id/restart", s.restartService)

		// Instala??es
		v1.POST("/installations", s.createInstallation)
		v1.GET("/installations/:id", s.getInstallation)
		v1.GET("/installations/:id/logs", s.getInstallationLogs)

		// Recursos
		v1.GET("/resources/available", s.getAvailableResources)
		v1.GET("/resources/usage", s.getResourceUsage)
		v1.GET("/clients/:id/resources", s.getClientResources)
	}
}

// Start inicia o servidor HTTP
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// healthCheck retorna o status de sa?de da API
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"service": "services-manager",
	})
}

// corsMiddleware configura CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
