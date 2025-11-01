package models

import (
	"time"
)

// Installation representa uma instala??o de servi?o
type Installation struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"service_id" db:"service_id"`
	Version     string    `json:"version" db:"version"`
	Status      string    `json:"status" db:"status"` // pending, installing, completed, failed
	Resources   Resources `json:"resources" db:"resources"`
	Logs        []string  `json:"logs" db:"logs"`
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Error       string    `json:"error,omitempty" db:"error"`
}

// CreateInstallationRequest representa a requisi??o de cria??o de instala??o
type CreateInstallationRequest struct {
	ServiceID string    `json:"service_id" binding:"required"`
	Version   string    `json:"version"`
	Resources Resources `json:"resources" binding:"required"`
}
