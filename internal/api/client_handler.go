package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/services-manager/internal/models"
)

// createClient cria um novo cliente
func (s *Server) createClient(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &models.Client{
		Name:     req.Name,
		Email:    req.Email,
		Status:   "active",
		Limits:   req.Limits,
		Metadata: req.Metadata,
	}

	if client.Limits.MaxServices == 0 {
		client.Limits.MaxServices = 10
	}
	if client.Limits.MaxCPU == 0 {
		client.Limits.MaxCPU = 4.0
	}
	if client.Limits.MaxMemory == "" {
		client.Limits.MaxMemory = "8g"
	}
	if client.Limits.MaxDisk == "" {
		client.Limits.MaxDisk = "100g"
	}

	if err := s.clientRepo.Create(client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// listClients lista todos os clientes
func (s *Server) listClients(c *gin.Context) {
	limit := 50
	offset := 0

	clients, err := s.clientRepo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// getClient retorna um cliente por ID
func (s *Server) getClient(c *gin.Context) {
	id := c.Param("id")

	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente n?o encontrado"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// updateClient atualiza um cliente
func (s *Server) updateClient(c *gin.Context) {
	id := c.Param("id")

	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente n?o encontrado"})
		return
	}

	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Email != nil {
		client.Email = *req.Email
	}
	if req.Status != nil {
		client.Status = *req.Status
	}
	if req.Limits != nil {
		client.Limits = *req.Limits
	}
	if req.Metadata != nil {
		client.Metadata = req.Metadata
	}

	if err := s.clientRepo.Update(client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// deleteClient remove um cliente
func (s *Server) deleteClient(c *gin.Context) {
	id := c.Param("id")

	if err := s.clientRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente removido com sucesso"})
}
