package models

import (
	"time"
)

// Service representa um servi?o gerenciado pelo sistema
type Service struct {
	ID          string                 `json:"id" db:"id"`
	ClientID    string                 `json:"client_id" db:"client_id"`
	Type        string                 `json:"type" db:"type"` // n8n, wordpress, custom
	Name        string                 `json:"name" db:"name"`
	Subdomain   string                 `json:"subdomain" db:"subdomain"`
	Status      string                 `json:"status" db:"status"` // created, running, stopped, error
	Resources   Resources              `json:"resources" db:"resources"`
	Config      map[string]interface{} `json:"config" db:"config"`
	ContainerID string                 `json:"container_id" db:"container_id"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// Resources define os recursos alocados para um servi?o
type Resources struct {
	CPU     CPUResources     `json:"cpu"`
	Memory  MemoryResources   `json:"memory"`
	Disk    DiskResources     `json:"disk"`
	Network NetworkResources  `json:"network"`
	Limits  ResourceLimits    `json:"limits"`
}

// CPUResources define recursos de CPU
type CPUResources struct {
	Cores    float64 `json:"cores"`      // Quantidade de cores CPU (ex: 0.5, 1.0, 2.0)
	MaxCores float64 `json:"max_cores"`  // Limite m?ximo de cores
	Priority int     `json:"priority"`    // Prioridade CPU (0-1024)
}

// MemoryResources define recursos de mem?ria
type MemoryResources struct {
	Limit       string `json:"limit"`        // Limite de mem?ria (ex: "512m", "1g", "2g")
	Reservation string `json:"reservation"`   // Reserva garantida (ex: "256m")
	Swap        string `json:"swap"`         // Limite de swap (ex: "512m", "-1" para ilimitado)
}

// DiskResources define recursos de disco
type DiskResources struct {
	Size       string `json:"size"`        // Tamanho do volume (ex: "10g", "50g", "100g")
	VolumeType string `json:"volume_type"`  // Tipo de volume (standard, ssd, nvme)
	IOPS       int    `json:"iops"`         // IOPS para discos SSD/NVMe
	Throughput int    `json:"throughput"`   // Throughput em MB/s
}

// NetworkResources define recursos de rede
type NetworkResources struct {
	Bandwidth      string `json:"bandwidth"`        // Largura de banda (ex: "100m", "1g")
	MaxConnections int    `json:"max_connections"`   // M?ximo de conex?es simult?neas
}

// ResourceLimits define limites adicionais
type ResourceLimits struct {
	Pids    int      `json:"pids"`     // M?ximo de processos
	Ulimits []ULimit `json:"ulimits"`  // Limites do sistema (nofile, nproc, etc.)
}

// ULimit define um limite do sistema
type ULimit struct {
	Name string `json:"name"`
	Soft int    `json:"soft"`
	Hard int    `json:"hard"`
}

// CreateServiceRequest representa a requisi??o de cria??o de servi?o
type CreateServiceRequest struct {
	ClientID  string                 `json:"client_id" binding:"required"`
	Type      string                 `json:"type" binding:"required,oneof=n8n wordpress custom"`
	Name      string                 `json:"name" binding:"required"`
	Subdomain string                 `json:"subdomain" binding:"required"`
	Resources Resources              `json:"resources" binding:"required"`
	Config    map[string]interface{} `json:"config"`
}

// UpdateServiceRequest representa a requisi??o de atualiza??o de servi?o
type UpdateServiceRequest struct {
	Name      *string                `json:"name"`
	Subdomain *string                `json:"subdomain"`
	Status    *string                `json:"status"`
	Resources *Resources             `json:"resources"`
	Config    map[string]interface{} `json:"config"`
}

// UpdateServiceResourcesRequest representa a requisi??o de atualiza??o de recursos
type UpdateServiceResourcesRequest struct {
	Resources Resources `json:"resources" binding:"required"`
}
