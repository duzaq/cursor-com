package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/services-manager/internal/models"
)

// NginxManager gerencia configura??es do Nginx
type NginxManager struct {
	ConfigDir      string // Diret?rio de configura??es do Nginx
	ConfigFileName string // Nome do arquivo de configura??o por servi?o
	ReloadCommand  string // Comando para recarregar Nginx
}

// NewNginxManager cria uma nova inst?ncia do NginxManager
func NewNginxManager(configDir string) *NginxManager {
	return &NginxManager{
		ConfigDir:      configDir,
		ConfigFileName: "services.conf",
		ReloadCommand:  "nginx -s reload",
	}
}

// AddService adiciona configura??o do Nginx para um servi?o
func (nm *NginxManager) AddService(service *models.Service, upstreamHost string, upstreamPort int) error {
	config := &models.NginxConfig{
		Subdomain:     service.Subdomain,
		UpstreamHost:  upstreamHost,
		UpstreamPort:  upstreamPort,
		SSLEnabled:    true,
		ForceHTTPS:    true,
		RateLimit:     "100r/m",
		ProxyTimeout:  60,
	}

	configContent := config.GenerateNginxConfig()

	// Criar diret?rio se n?o existir
	if err := os.MkdirAll(nm.ConfigDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diret?rio de configura??o: %w", err)
	}

	// Arquivo de configura??o do servi?o
	configPath := filepath.Join(nm.ConfigDir, fmt.Sprintf("%s.conf", service.Subdomain))

	// Escrever configura??o
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever configura??o: %w", err)
	}

	// Recarregar Nginx
	if err := nm.Reload(); err != nil {
		return fmt.Errorf("erro ao recarregar Nginx: %w", err)
	}

	return nil
}

// RemoveService remove configura??o do Nginx para um servi?o
func (nm *NginxManager) RemoveService(subdomain string) error {
	configPath := filepath.Join(nm.ConfigDir, fmt.Sprintf("%s.conf", subdomain))

	// Remover arquivo de configura??o
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao remover configura??o: %w", err)
	}

	// Recarregar Nginx
	if err := nm.Reload(); err != nil {
		return fmt.Errorf("erro ao recarregar Nginx: %w", err)
	}

	return nil
}

// UpdateService atualiza configura??o do Nginx para um servi?o
func (nm *NginxManager) UpdateService(service *models.Service, upstreamHost string, upstreamPort int) error {
	return nm.AddService(service, upstreamHost, upstreamPort)
}

// Reload recarrega a configura??o do Nginx
func (nm *NginxManager) Reload() error {
	// Verificar se Nginx est? rodando em container
	if os.Getenv("NGINX_CONTAINER") != "" {
		// Executar reload dentro do container
		cmd := exec.Command("docker", "exec", os.Getenv("NGINX_CONTAINER"), "nginx", "-s", "reload")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao recarregar Nginx no container: %w", err)
		}
		return nil
	}

	// Executar reload local
	cmd := exec.Command("sh", "-c", nm.ReloadCommand)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao recarregar Nginx: %w", err)
	}

	return nil
}

// GenerateMainConfig gera a configura??o principal do Nginx
func (nm *NginxManager) GenerateMainConfig() string {
	return `
# Configura??o principal do Nginx para gerenciamento de servi?os

# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

# Upstream para API do gerenciador
upstream api_backend {
    server localhost:8080;
}

# Servidor principal para API de gerenciamento
server {
    listen 80;
    server_name api.example.com;  # Substituir pelo dom?nio real
    
    location / {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Incluir configura??es din?micas dos servi?os
include ` + nm.ConfigDir + `/*.conf;
`
}

// ValidateConfig valida a configura??o do Nginx
func (nm *NginxManager) ValidateConfig() error {
	var cmd *exec.Cmd

	if os.Getenv("NGINX_CONTAINER") != "" {
		cmd = exec.Command("docker", "exec", os.Getenv("NGINX_CONTAINER"), "nginx", "-t")
	} else {
		cmd = exec.Command("nginx", "-t")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro na valida??o do Nginx: %s", string(output))
	}

	return nil
}

// GenerateUpstreamConfig gera configura??o de upstream para um servi?o
func (nm *NginxManager) GenerateUpstreamConfig(subdomain string, hosts []string) string {
	var upstream strings.Builder
	upstream.WriteString(fmt.Sprintf("upstream %s_backend {\n", subdomain))

	for _, host := range hosts {
		upstream.WriteString(fmt.Sprintf("    server %s;\n", host))
	}

	upstream.WriteString("    least_conn;\n") // Load balancing
	upstream.WriteString("}\n\n")

	return upstream.String()
}
