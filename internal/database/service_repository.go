package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/services-manager/internal/models"
)

// ServiceRepository gerencia opera??es de servi?os no banco de dados
type ServiceRepository struct {
	db *sql.DB
}

// NewServiceRepository cria uma nova inst?ncia do ServiceRepository
func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create cria um novo servi?o
func (r *ServiceRepository) Create(service *models.Service) error {
	service.ID = uuid.New().String()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	resourcesJSON, _ := json.Marshal(service.Resources)
	configJSON, _ := json.Marshal(service.Config)

	query := `
		INSERT INTO services (id, client_id, type, name, subdomain, status, resources, config, container_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(
		query,
		service.ID,
		service.ClientID,
		service.Type,
		service.Name,
		service.Subdomain,
		service.Status,
		resourcesJSON,
		configJSON,
		service.ContainerID,
		service.CreatedAt,
		service.UpdatedAt,
	)

	return err
}

// GetByID retorna um servi?o por ID
func (r *ServiceRepository) GetByID(id string) (*models.Service, error) {
	var service models.Service
	var resourcesJSON, configJSON []byte

	query := `
		SELECT id, client_id, type, name, subdomain, status, resources, config, container_id, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&service.ID,
		&service.ClientID,
		&service.Type,
		&service.Name,
		&service.Subdomain,
		&service.Status,
		&resourcesJSON,
		&configJSON,
		&service.ContainerID,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resourcesJSON, &service.Resources)
	json.Unmarshal(configJSON, &service.Config)

	return &service, nil
}

// GetBySubdomain retorna um servi?o por subdom?nio
func (r *ServiceRepository) GetBySubdomain(subdomain string) (*models.Service, error) {
	var service models.Service
	var resourcesJSON, configJSON []byte

	query := `
		SELECT id, client_id, type, name, subdomain, status, resources, config, container_id, created_at, updated_at
		FROM services
		WHERE subdomain = $1
	`

	err := r.db.QueryRow(query, subdomain).Scan(
		&service.ID,
		&service.ClientID,
		&service.Type,
		&service.Name,
		&service.Subdomain,
		&service.Status,
		&resourcesJSON,
		&configJSON,
		&service.ContainerID,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resourcesJSON, &service.Resources)
	json.Unmarshal(configJSON, &service.Config)

	return &service, nil
}

// List retorna todos os servi?os
func (r *ServiceRepository) List(clientID string, limit, offset int) ([]*models.Service, error) {
	var query string
	var rows *sql.Rows
	var err error

	if clientID != "" {
		query = `
			SELECT id, client_id, type, name, subdomain, status, resources, config, container_id, created_at, updated_at
			FROM services
			WHERE client_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.db.Query(query, clientID, limit, offset)
	} else {
		query = `
			SELECT id, client_id, type, name, subdomain, status, resources, config, container_id, created_at, updated_at
			FROM services
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = r.db.Query(query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		var service models.Service
		var resourcesJSON, configJSON []byte

		err := rows.Scan(
			&service.ID,
			&service.ClientID,
			&service.Type,
			&service.Name,
			&service.Subdomain,
			&service.Status,
			&resourcesJSON,
			&configJSON,
			&service.ContainerID,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(resourcesJSON, &service.Resources)
		json.Unmarshal(configJSON, &service.Config)

		services = append(services, &service)
	}

	return services, nil
}

// Update atualiza um servi?o
func (r *ServiceRepository) Update(service *models.Service) error {
	service.UpdatedAt = time.Now()

	resourcesJSON, _ := json.Marshal(service.Resources)
	configJSON, _ := json.Marshal(service.Config)

	query := `
		UPDATE services
		SET name = $2, subdomain = $3, status = $4, resources = $5, config = $6, container_id = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		service.ID,
		service.Name,
		service.Subdomain,
		service.Status,
		resourcesJSON,
		configJSON,
		service.ContainerID,
		service.UpdatedAt,
	)

	return err
}

// Delete remove um servi?o
func (r *ServiceRepository) Delete(id string) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
