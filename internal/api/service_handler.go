package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/services-manager/internal/models"
)

// createService cria um novo servi?o
func (s *Server) createService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar cliente existe
	_, err := s.clientRepo.GetByID(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cliente n?o encontrado"})
		return
	}

	// Validar subdom?nio ?nico
	_, err = s.serviceRepo.GetBySubdomain(req.Subdomain)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Subdom?nio j? est? em uso"})
		return
	}

	service := &models.Service{
		ClientID:  req.ClientID,
		Type:      req.Type,
		Name:      req.Name,
		Subdomain: req.Subdomain,
		Status:    "created",
		Resources: req.Resources,
		Config:    req.Config,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// listServices lista todos os servi?os
func (s *Server) listServices(c *gin.Context) {
	clientID := c.Query("client_id")
	limit := 50
	offset := 0

	services, err := s.serviceRepo.List(clientID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// getService retorna um servi?o por ID
func (s *Server) getService(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// updateService atualiza um servi?o
func (s *Server) updateService(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		service.Name = *req.Name
	}
	if req.Subdomain != nil {
		service.Subdomain = *req.Subdomain
	}
	if req.Status != nil {
		service.Status = *req.Status
	}
	if req.Resources != nil {
		service.Resources = *req.Resources
	}
	if req.Config != nil {
		service.Config = req.Config
	}

	if err := s.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// updateServiceResources atualiza apenas os recursos de um servi?o
func (s *Server) updateServiceResources(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	var req models.UpdateServiceResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Resources = req.Resources

	if err := s.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// deleteService remove um servi?o
func (s *Server) deleteService(c *gin.Context) {
	id := c.Param("id")

	if err := s.serviceRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Servi?o removido com sucesso"})
}

// startService inicia um servi?o
func (s *Server) startService(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	// TODO: Implementar l?gica de iniciar container Docker
	service.Status = "running"

	if err := s.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// stopService para um servi?o
func (s *Server) stopService(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	// TODO: Implementar l?gica de parar container Docker
	service.Status = "stopped"

	if err := s.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// restartService reinicia um servi?o
func (s *Server) restartService(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	// TODO: Implementar l?gica de reiniciar container Docker
	service.Status = "running"

	if err := s.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// getServiceResourceUsage retorna o uso de recursos de um servi?o
func (s *Server) getServiceResourceUsage(c *gin.Context) {
	id := c.Param("id")

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Servi?o n?o encontrado"})
		return
	}

	// TODO: Implementar coleta de m?tricas reais do container
	c.JSON(http.StatusOK, gin.H{
		"service_id": service.ID,
		"resources":  service.Resources,
		"usage": gin.H{
			"cpu_percent":    0.0,
			"memory_used":    "0m",
			"memory_percent": 0.0,
			"disk_used":      "0g",
			"disk_percent":   0.0,
		},
	})
}
