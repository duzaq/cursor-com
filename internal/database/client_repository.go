package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/services-manager/internal/models"
)

// ClientRepository gerencia opera??es de clientes no banco de dados
type ClientRepository struct {
	db *sql.DB
}

// NewClientRepository cria uma nova inst?ncia do ClientRepository
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// Create cria um novo cliente
func (r *ClientRepository) Create(client *models.Client) error {
	client.ID = uuid.New().String()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	limitsJSON, _ := json.Marshal(client.Limits)
	metadataJSON, _ := json.Marshal(client.Metadata)

	query := `
		INSERT INTO clients (id, name, email, status, limits, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		client.ID,
		client.Name,
		client.Email,
		client.Status,
		limitsJSON,
		metadataJSON,
		client.CreatedAt,
		client.UpdatedAt,
	)

	return err
}

// GetByID retorna um cliente por ID
func (r *ClientRepository) GetByID(id string) (*models.Client, error) {
	var client models.Client
	var limitsJSON, metadataJSON []byte

	query := `
		SELECT id, name, email, status, limits, metadata, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.Status,
		&limitsJSON,
		&metadataJSON,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(limitsJSON, &client.Limits)
	json.Unmarshal(metadataJSON, &client.Metadata)

	return &client, nil
}

// GetByEmail retorna um cliente por email
func (r *ClientRepository) GetByEmail(email string) (*models.Client, error) {
	var client models.Client
	var limitsJSON, metadataJSON []byte

	query := `
		SELECT id, name, email, status, limits, metadata, created_at, updated_at
		FROM clients
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.Status,
		&limitsJSON,
		&metadataJSON,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(limitsJSON, &client.Limits)
	json.Unmarshal(metadataJSON, &client.Metadata)

	return &client, nil
}

// List retorna todos os clientes
func (r *ClientRepository) List(limit, offset int) ([]*models.Client, error) {
	query := `
		SELECT id, name, email, status, limits, metadata, created_at, updated_at
		FROM clients
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		var client models.Client
		var limitsJSON, metadataJSON []byte

		err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Email,
			&client.Status,
			&limitsJSON,
			&metadataJSON,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(limitsJSON, &client.Limits)
		json.Unmarshal(metadataJSON, &client.Metadata)

		clients = append(clients, &client)
	}

	return clients, nil
}

// Update atualiza um cliente
func (r *ClientRepository) Update(client *models.Client) error {
	client.UpdatedAt = time.Now()

	limitsJSON, _ := json.Marshal(client.Limits)
	metadataJSON, _ := json.Marshal(client.Metadata)

	query := `
		UPDATE clients
		SET name = $2, email = $3, status = $4, limits = $5, metadata = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		client.ID,
		client.Name,
		client.Email,
		client.Status,
		limitsJSON,
		metadataJSON,
		client.UpdatedAt,
	)

	return err
}

// Delete remove um cliente
func (r *ClientRepository) Delete(id string) error {
	query := `DELETE FROM clients WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
