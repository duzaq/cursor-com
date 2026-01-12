package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/services-manager/internal/models"
)

// InstallationRepository gerencia opera??es de instala??es no banco de dados
type InstallationRepository struct {
	db *sql.DB
}

// NewInstallationRepository cria uma nova inst?ncia do InstallationRepository
func NewInstallationRepository(db *sql.DB) *InstallationRepository {
	return &InstallationRepository{db: db}
}

// Create cria uma nova instala??o
func (r *InstallationRepository) Create(installation *models.Installation) error {
	installation.ID = uuid.New().String()
	installation.StartedAt = time.Now()

	resourcesJSON, _ := json.Marshal(installation.Resources)
	logsJSON, _ := json.Marshal(installation.Logs)

	query := `
		INSERT INTO installations (id, service_id, version, status, resources, logs, error, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		installation.ID,
		installation.ServiceID,
		installation.Version,
		installation.Status,
		resourcesJSON,
		logsJSON,
		installation.Error,
		installation.StartedAt,
	)

	return err
}

// GetByID retorna uma instala??o por ID
func (r *InstallationRepository) GetByID(id string) (*models.Installation, error) {
	var installation models.Installation
	var resourcesJSON, logsJSON []byte

	query := `
		SELECT id, service_id, version, status, resources, logs, error, started_at, completed_at
		FROM installations
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&installation.ID,
		&installation.ServiceID,
		&installation.Version,
		&installation.Status,
		&resourcesJSON,
		&logsJSON,
		&installation.Error,
		&installation.StartedAt,
		&installation.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resourcesJSON, &installation.Resources)
	json.Unmarshal(logsJSON, &installation.Logs)

	return &installation, nil
}

// Update atualiza uma instala??o
func (r *InstallationRepository) Update(installation *models.Installation) error {
	resourcesJSON, _ := json.Marshal(installation.Resources)
	logsJSON, _ := json.Marshal(installation.Logs)

	query := `
		UPDATE installations
		SET status = $2, resources = $3, logs = $4, error = $5, completed_at = $6
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		installation.ID,
		installation.Status,
		resourcesJSON,
		logsJSON,
		installation.Error,
		installation.CompletedAt,
	)

	return err
}
