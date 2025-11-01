package models

// NginxConfig representa a configura??o do Nginx para um servi?o
type NginxConfig struct {
	Subdomain     string `json:"subdomain"`
	UpstreamHost  string `json:"upstream_host"`  // IP do container Docker
	UpstreamPort  int    `json:"upstream_port"`   // Porta do container
	SSLEnabled    bool   `json:"ssl_enabled"`     // SSL via Cloudflare
	ForceHTTPS    bool   `json:"force_https"`     // Redirect HTTP -> HTTPS
	RateLimit     string `json:"rate_limit"`       // Rate limiting (ex: "100r/m")
	ProxyTimeout  int    `json:"proxy_timeout"`   // Timeout em segundos
}

// GenerateNginxConfig gera a configura??o do Nginx para um servi?o
func (nc *NginxConfig) GenerateNginxConfig() string {
	config := `
# Configura??o para ` + nc.Subdomain + `
server {
    listen 80;
    server_name ` + nc.Subdomain + `.` + `example.com` + `;
    
    # Rate limiting
    limit_req zone=api_limit burst=10 nodelay;
    
    # Headers de seguran?a
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Logs
    access_log /var/log/nginx/` + nc.Subdomain + `.access.log;
    error_log /var/log/nginx/` + nc.Subdomain + `.error.log;
    
    # Proxy settings
    proxy_connect_timeout ` + string(rune(nc.ProxyTimeout)) + `s;
    proxy_send_timeout ` + string(rune(nc.ProxyTimeout)) + `s;
    proxy_read_timeout ` + string(rune(nc.ProxyTimeout)) + `s;
    
    # Headers para proxy
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;
    
    # Cloudflare headers
    proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
    proxy_set_header CF-Ray $http_cf_ray;
    proxy_set_header CF-Visitor $http_cf_visitor;
    
    location / {
        proxy_pass http://` + nc.UpstreamHost + `:` + string(rune(nc.UpstreamPort)) + `;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://` + nc.UpstreamHost + `:` + string(rune(nc.UpstreamPort)) + `/health;
    }
`

	if nc.ForceHTTPS {
		config += `
    # Redirect HTTP to HTTPS
    if ($scheme != "https") {
        return 301 https://$server_name$request_uri;
    }
`
	}

	if nc.SSLEnabled {
		config += `
    # SSL configuration (Cloudflare SSL)
    listen 443 ssl http2;
    ssl_certificate /etc/nginx/ssl/cloudflare.crt;
    ssl_certificate_key /etc/nginx/ssl/cloudflare.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
`
	}

	config += `
}
`

	return config
}
