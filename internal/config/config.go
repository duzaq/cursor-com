package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config representa a configura??o da aplica??o
type Config struct {
	Database   DatabaseConfig `yaml:"database"`
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
	Docker     DockerConfig   `yaml:"docker"`
	API        APIConfig      `yaml:"api"`
	Redis      RedisConfig    `yaml:"redis"`
}

// DatabaseConfig configura??o do banco de dados
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// CloudflareConfig configura??o do Cloudflare
type CloudflareConfig struct {
	APIKey      string `yaml:"api_key"`
	Email       string `yaml:"email"`
	ZoneID      string `yaml:"zone_id"`
	BaseDomain  string `yaml:"base_domain"`
	APIBaseURL  string `yaml:"api_base_url"`
}

// DockerConfig configura??o do Docker
type DockerConfig struct {
	Host string `yaml:"host"`
}

// APIConfig configura??o da API
type APIConfig struct {
	Port       string `yaml:"port"`
	SecretKey  string `yaml:"secret_key"`
	JWTSecret  string `yaml:"jwt_secret"`
}

// RedisConfig configura??o do Redis
type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

// Load carrega a configura??o do arquivo YAML e vari?veis de ambiente
func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg := &Config{}

	// Ler arquivo YAML se existir
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler arquivo de configura??o: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("erro ao fazer parse do arquivo de configura??o: %w", err)
		}
	}

	// Sobrescrever com vari?veis de ambiente
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Database.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}

	if key := os.Getenv("CLOUDFLARE_API_KEY"); key != "" {
		cfg.Cloudflare.APIKey = key
	}
	if email := os.Getenv("CLOUDFLARE_EMAIL"); email != "" {
		cfg.Cloudflare.Email = email
	}
	if zoneID := os.Getenv("CLOUDFLARE_ZONE_ID"); zoneID != "" {
		cfg.Cloudflare.ZoneID = zoneID
	}
	if domain := os.Getenv("CLOUDFLARE_BASE_DOMAIN"); domain != "" {
		cfg.Cloudflare.BaseDomain = domain
	}
	if apiURL := os.Getenv("CLOUDFLARE_API_BASE_URL"); apiURL != "" {
		cfg.Cloudflare.APIBaseURL = apiURL
	} else {
		cfg.Cloudflare.APIBaseURL = "https://api.cloudflare.com/client/v4"
	}

	if host := os.Getenv("DOCKER_HOST"); host != "" {
		cfg.Docker.Host = host
	} else {
		cfg.Docker.Host = "unix:///var/run/docker.sock"
	}

	if port := os.Getenv("API_PORT"); port != "" {
		cfg.API.Port = port
	} else {
		cfg.API.Port = "8080"
	}
	if secret := os.Getenv("API_SECRET_KEY"); secret != "" {
		cfg.API.SecretKey = secret
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.API.JWTSecret = jwtSecret
	}

	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	} else {
		cfg.Redis.Host = "localhost"
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Redis.Port)
	} else {
		cfg.Redis.Port = 6379
	}

	return cfg, nil
}
