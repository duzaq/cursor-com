package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/services-manager/internal/models"
)

// NginxManager gerencia configura??es do Nginx
type NginxManager struct {
	configDir    string // Diret?rio de configura??es do Nginx
	templateDir  string // Diret?rio de templates
	nginxCmd     string // Comando para reload do Nginx
}

// NewNginxManager cria uma nova inst?ncia do NginxManager
func NewNginxManager(configDir, templateDir, nginxCmd string) *NginxManager {
	return &NginxManager{
		configDir:   configDir,
		templateDir: templateDir,
		nginxCmd:    nginxCmd,
	}
}

// CreateServiceConfig cria configura??o do Nginx para um servi?o
func (nm *NginxManager) CreateServiceConfig(service *models.Service, nginxConfig models.NginxConfig) error {
	// Criar diret?rio se n?o existir
	if err := os.MkdirAll(nm.configDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diret?rio de configura??o: %w", err)
	}

	// Configurar upstream host se n?o especificado
	if nginxConfig.UpstreamHost == "" {
		nginxConfig.UpstreamHost = service.ContainerID
	}
	if nginxConfig.UpstreamPort == 0 {
		nginxConfig.UpstreamPort = getDefaultPort(service.Type)
	}

	// Criar configura??o de upstream
	upstreamConfig := nm.generateUpstreamConfig(service, nginxConfig)
	upstreamFile := filepath.Join(nm.configDir, fmt.Sprintf("upstream_%s.conf", service.ID))
	if err := os.WriteFile(upstreamFile, []byte(upstreamConfig), 0644); err != nil {
		return fmt.Errorf("erro ao criar configura??o de upstream: %w", err)
	}

	// Criar configura??o de server block
	serverConfig := nm.generateServerConfig(service, nginxConfig)
	serverFile := filepath.Join(nm.configDir, fmt.Sprintf("server_%s.conf", service.ID))
	if err := os.WriteFile(serverFile, []byte(serverConfig), 0644); err != nil {
		return fmt.Errorf("erro ao criar configura??o de server: %w", err)
	}

	// Reload do Nginx
	return nm.Reload()
}

// UpdateServiceConfig atualiza configura??o do Nginx para um servi?o
func (nm *NginxManager) UpdateServiceConfig(service *models.Service, nginxConfig models.NginxConfig) error {
	return nm.CreateServiceConfig(service, nginxConfig)
}

// DeleteServiceConfig remove configura??o do Nginx de um servi?o
func (nm *NginxManager) DeleteServiceConfig(serviceID string) error {
	upstreamFile := filepath.Join(nm.configDir, fmt.Sprintf("upstream_%s.conf", serviceID))
	serverFile := filepath.Join(nm.configDir, fmt.Sprintf("server_%s.conf", serviceID))

	if err := os.Remove(upstreamFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao remover configura??o de upstream: %w", err)
	}

	if err := os.Remove(serverFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao remover configura??o de server: %w", err)
	}

	// Reload do Nginx
	return nm.Reload()
}

// Reload recarrega configura??o do Nginx
func (nm *NginxManager) Reload() error {
	if nm.nginxCmd == "" {
		nm.nginxCmd = "nginx"
	}

	// Testar configura??o antes de recarregar
	testCmd := exec.Command(nm.nginxCmd, "-t")
	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("erro na configura??o do Nginx: %w", err)
	}

	// Recarregar Nginx
	reloadCmd := exec.Command(nm.nginxCmd, "-s", "reload")
	if err := reloadCmd.Run(); err != nil {
		return fmt.Errorf("erro ao recarregar Nginx: %w", err)
	}

	return nil
}

// generateUpstreamConfig gera configura??o de upstream
func (nm *NginxManager) generateUpstreamConfig(service *models.Service, nginxConfig models.NginxConfig) string {
	return fmt.Sprintf(`upstream %s {
    server %s:%d;
    keepalive 32;
}
`, service.ID, nginxConfig.UpstreamHost, nginxConfig.UpstreamPort)
}

// generateServerConfig gera configura??o de server block
func (nm *NginxManager) generateServerConfig(service *models.Service, nginxConfig models.NginxConfig) string {
	baseDomain := os.Getenv("CLOUDFLARE_BASE_DOMAIN")
	if baseDomain == "" {
		baseDomain = "example.com"
	}

	serverName := fmt.Sprintf("%s.%s", nginxConfig.Subdomain, baseDomain)

	// Template de configura??o
	tmpl := `server {
    listen 80;
    server_name {{ .ServerName }};

    {{ if .SSLEnabled }}
    # Redirecionar HTTP para HTTPS (via Cloudflare)
    return 301 https://$server_name$request_uri;
    {{ else }}
    # Configura??o HTTP
    {{ end }}

    {{ if .SSLEnabled }}
    listen 443 ssl http2;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    # SSL via Cloudflare (proxy mode)
    {{ end }}

    {{ if .RateLimit }}
    limit_req_zone $binary_remote_addr zone={{ .RateLimitZone }}:{{ .RateLimitPeriod }} rate={{ .RateLimitRate }}r/s;
    limit_req zone={{ .RateLimitZone }} burst={{ .RateLimitBurst }} nodelay;
    {{ end }}

    # Logs
    access_log /var/log/nginx/{{ .ServiceID }}_access.log;
    error_log /var/log/nginx/{{ .ServiceID }}_error.log;

    # Configura??es de proxy
    proxy_connect_timeout {{ .ConnectTimeout }};
    proxy_send_timeout {{ .SendTimeout }};
    proxy_read_timeout {{ .ReadTimeout }};
    {{ if .ProxyBuffering }}
    proxy_buffering on;
    proxy_buffer_size {{ .ProxyBufferSize }};
    proxy_buffers {{ .ProxyBuffers }};
    {{ else }}
    proxy_buffering off;
    {{ end }}
    client_max_body_size {{ .MaxBodySize }};

    # Headers
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;

    # Upstream
    location / {
        proxy_pass http://{{ .ServiceID }};
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }

    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://{{ .ServiceID }};
    }
}
`

	// Calcular rate limit rate e per?odo
	rateLimitRate := ""
	rateLimitPeriod := "10m"
	rateLimitBurst := 10
	if nginxConfig.RateLimit != nil {
		rateLimitRate = fmt.Sprintf("%d/%s", nginxConfig.RateLimit.Requests, nginxConfig.RateLimit.Period)
		rateLimitPeriod = nginxConfig.RateLimit.Period
		rateLimitBurst = nginxConfig.RateLimit.Burst
		if rateLimitBurst == 0 {
			rateLimitBurst = 10
		}
	}

	// Valores padr?o para proxy settings
	connectTimeout := nginxConfig.ProxySettings.ConnectTimeout
	if connectTimeout == "" {
		connectTimeout = "30s"
	}
	sendTimeout := nginxConfig.ProxySettings.SendTimeout
	if sendTimeout == "" {
		sendTimeout = "30s"
	}
	readTimeout := nginxConfig.ProxySettings.ReadTimeout
	if readTimeout == "" {
		readTimeout = "30s"
	}
	proxyBufferSize := nginxConfig.ProxySettings.ProxyBufferSize
	if proxyBufferSize == "" {
		proxyBufferSize = "4k"
	}
	proxyBuffers := nginxConfig.ProxySettings.ProxyBuffers
	if proxyBuffers == "" {
		proxyBuffers = "8 4k"
	}
	maxBodySize := nginxConfig.ProxySettings.MaxBodySize
	if maxBodySize == "" {
		maxBodySize = "10m"
	}

	data := struct {
		ServerName       string
		SSLEnabled       bool
		RateLimitZone    string
		RateLimitPeriod  string
		RateLimitRate    string
		RateLimitBurst   int
		ServiceID        string
		ConnectTimeout   string
		SendTimeout      string
		ReadTimeout      string
		ProxyBuffering   bool
		ProxyBufferSize  string
		ProxyBuffers     string
		MaxBodySize      string
	}{
		ServerName:      serverName,
		SSLEnabled:      nginxConfig.SSLEnabled,
		RateLimitZone:   fmt.Sprintf("limit_%s", service.ID),
		RateLimitPeriod: rateLimitPeriod,
		RateLimitRate:   rateLimitRate,
		RateLimitBurst:  rateLimitBurst,
		ServiceID:       service.ID,
		ConnectTimeout:  connectTimeout,
		SendTimeout:     sendTimeout,
		ReadTimeout:     readTimeout,
		ProxyBuffering:  nginxConfig.ProxySettings.ProxyBuffering,
		ProxyBufferSize: proxyBufferSize,
		ProxyBuffers:    proxyBuffers,
		MaxBodySize:     maxBodySize,
	}

	t := template.Must(template.New("server").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Sprintf("# Erro ao gerar template: %v", err)
	}

	return buf.String()
}

// getDefaultPort retorna a porta padr?o para um tipo de servi?o
func getDefaultPort(serviceType string) int {
	ports := map[string]int{
		"n8n":       5678,
		"wordpress": 80,
		"custom":    8080,
	}
	if port, ok := ports[serviceType]; ok {
		return port
	}
	return 8080
}
