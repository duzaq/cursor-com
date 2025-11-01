package models

// NginxConfig representa a configura??o do Nginx para um servi?o
type NginxConfig struct {
	Subdomain     string   `json:"subdomain"`
	UpstreamHost  string   `json:"upstream_host"`  // Host do container (ex: container_name:port)
	UpstreamPort  int      `json:"upstream_port"`   // Porta do container
	SSLEnabled    bool     `json:"ssl_enabled"`     // SSL via Cloudflare
	RateLimit     *RateLimit `json:"rate_limit,omitempty"`
	ProxySettings ProxySettings `json:"proxy_settings"`
}

// RateLimit configura??o de rate limiting
type RateLimit struct {
	Zone        string `json:"zone"`         // Nome da zona de rate limit
	Requests    int    `json:"requests"`     // N?mero de requests
	Period      string `json:"period"`       // Per?odo (ex: "1m", "1h")
	Burst       int    `json:"burst"`        // Burst permitido
}

// ProxySettings configura??es do proxy
type ProxySettings struct {
	ConnectTimeout    string `json:"connect_timeout"`    // Ex: "30s"
	SendTimeout       string `json:"send_timeout"`       // Ex: "30s"
	ReadTimeout       string `json:"read_timeout"`       // Ex: "30s"
	ProxyBuffering    bool   `json:"proxy_buffering"`   // Ativar buffering
	ProxyBufferSize   string `json:"proxy_buffer_size"` // Ex: "4k"
	ProxyBuffers      string `json:"proxy_buffers"`    // Ex: "8 4k"
	MaxBodySize       string `json:"max_body_size"`     // Ex: "10m"
}
