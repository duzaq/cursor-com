package models

import (
	"time"
)

// Client representa um cliente do sistema
type Client struct {
	ID        string            `json:"id" db:"id"`
	Name      string            `json:"name" db:"name"`
	Email     string            `json:"email" db:"email"`
	Status    string            `json:"status" db:"status"` // active, inactive, suspended
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
	Limits    ClientLimits      `json:"limits" db:"limits"`
	Metadata  map[string]string `json:"metadata" db:"metadata"`
}

// ClientLimits define os limites de recursos por cliente
type ClientLimits struct {
	MaxServices    int     `json:"max_services"`
	MaxCPU         float64 `json:"max_cpu"`
	MaxMemory      string  `json:"max_memory"`      // Ex: "10g"
	MaxDisk        string  `json:"max_disk"`        // Ex: "100g"
	MaxSubdomains  int     `json:"max_subdomains"`
}

// CreateClientRequest representa a requisi??o de cria??o de cliente
type CreateClientRequest struct {
	Name     string            `json:"name" binding:"required"`
	Email    string            `json:"email" binding:"required,email"`
	Limits   ClientLimits      `json:"limits"`
	Metadata map[string]string `json:"metadata"`
}

// UpdateClientRequest representa a requisi??o de atualiza??o de cliente
type UpdateClientRequest struct {
	Name     *string            `json:"name"`
	Email    *string            `json:"email"`
	Status   *string            `json:"status"`
	Limits   *ClientLimits      `json:"limits"`
	Metadata map[string]string `json:"metadata"`
}
