package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/services-manager/internal/models"
)

// createInstallation cria uma nova instala??o
func (s *Server) createInstallation(c *gin.Context) {
	var req models.CreateInstallationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar servi?o existe
	_, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	if req.Version == "" {
		req.Version = "latest"
	}

	installation := &models.Installation{
		ServiceID: req.ServiceID,
		Version:   req.Version,
		Status:   "pending",
		Resources: req.Resources,
		Logs:     []string{},
	}

	if err := s.installationRepo.Create(installation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Iniciar processo ass?ncrono de instala??o

	c.JSON(http.StatusCreated, installation)
}

// getInstallation retorna uma instala??o por ID
func (s *Server) getInstallation(c *gin.Context) {
	id := c.Param("id")

	installation, err := s.installationRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instala??o n?o encontrada"})
		return
	}

	c.JSON(http.StatusOK, installation)
}

// getInstallationLogs retorna os logs de uma instala??o
func (s *Server) getInstallationLogs(c *gin.Context) {
	id := c.Param("id")

	installation, err := s.installationRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instala??o n?o encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"installation_id": installation.ID,
		"logs":            installation.Logs,
	})
}

// getAvailableResources retorna recursos dispon?veis no host
func (s *Server) getAvailableResources(c *gin.Context) {
	// TODO: Implementar coleta real de recursos do host
	c.JSON(http.StatusOK, models.HostResources{
		TotalCPU:        8.0,
		TotalMemory:     16 * 1024 * 1024 * 1024, // 16GB
		TotalDisk:       500 * 1024 * 1024 * 1024, // 500GB
		UsedCPU:         2.0,
		UsedMemory:      4 * 1024 * 1024 * 1024, // 4GB
		UsedDisk:        50 * 1024 * 1024 * 1024, // 50GB
		AvailableCPU:    6.0,
		AvailableMemory: 12 * 1024 * 1024 * 1024, // 12GB
		AvailableDisk:   450 * 1024 * 1024 * 1024, // 450GB
	})
}

// getResourceUsage retorna uso de recursos agregado
func (s *Server) getResourceUsage(c *gin.Context) {
	// TODO: Implementar c?lculo real de uso de recursos
	c.JSON(http.StatusOK, gin.H{
		"total_services": 0,
		"total_cpu":      0.0,
		"total_memory":   "0g",
		"total_disk":     "0g",
	})
}

// getClientResources retorna recursos utilizados por um cliente
func (s *Server) getClientResources(c *gin.Context) {
	clientID := c.Param("id")

	services, err := s.serviceRepo.List(clientID, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Calcular recursos agregados dos servi?os do cliente
	c.JSON(http.StatusOK, gin.H{
		"client_id":      clientID,
		"total_services": len(services),
		"total_cpu":      0.0,
		"total_memory":   "0g",
		"total_disk":     "0g",
	})
}
