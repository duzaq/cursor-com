# Roadmap: Sistema de Gerenciamento Dinâmico de Serviços

## Visão Geral
Sistema desenvolvido em Go para gerenciar dinamicamente múltiplos serviços (n8n, WordPress, outras aplicações) com integração Cloudflare para criação automática de subdomínios.

---

## 1. Arquitetura e Estrutura do Projeto

### 1.1 Estrutura de Diretórios
```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes.go
│   ├── services/
│   │   ├── n8n/
│   │   ├── wordpress/
│   │   ├── docker/
│   │   └── cloudflare/
│   ├── models/
│   ├── database/
│   └── config/
├── pkg/
│   ├── docker/
│   ├── cloudflare/
│   └── monitoring/
├── deployments/
│   ├── docker-compose.yml
│   └── kubernetes/
├── scripts/
│   ├── install.sh
│   └── deploy.sh
├── configs/
│   └── config.yaml
├── go.mod
├── go.sum
└── README.md
```

### 1.2 Tecnologias e Dependências
- **Go 1.21+**
- **Docker & Docker Compose** (para containers)
- **PostgreSQL** (banco de dados)
- **Redis** (cache e filas)
- **Nginx Reverse Proxy** (roteamento por subdomínio)
- **Cloudflare API v4**
- **Prometheus** (métricas)
- **Grafana** (dashboards)

---

## 2. Funcionalidades Principais

### 2.1 Gestão de Clientes
- CRUD completo de clientes
- Autenticação e autorização por cliente
- Quotas e limites por cliente
- Histórico de atividades
- Faturamento e cobrança

### 2.2 Instalação de Serviços
- Instalação automática via Docker
- Templates configuráveis por tipo de serviço
- Provisionamento de recursos (CPU, memória, disco)
- Configuração inicial automática
- Validação pós-instalação

### 2.3 Criação e Edição de Serviços
- API REST para criação de instâncias
- **Definição obrigatória de recursos** (CPU, Memória, Disco, Network)
- Validação de recursos disponíveis antes da criação
- Configuração dinâmica de serviços
- Atualização de recursos sem downtime (quando possível)
- Atualização sem downtime
- Versionamento de configurações
- Rollback automático
- **Ajuste de recursos** em tempo de execução (upscale/downscale)

### 2.4 Gerenciamento de Serviços
- Start/Stop/Restart de serviços
- **Gerenciamento de recursos**:
  - Ajuste de CPU (cores, shares, priority)
  - Ajuste de memória (limit, reservation, swap)
  - Ajuste de disco (expandir volumes)
  - Ajuste de bandwidth
- Escalonamento automático baseado em métricas
- Escalonamento manual de recursos
- Backup e restore
- Logs centralizados
- Health checks
- **Monitoramento de uso de recursos** em tempo real

### 2.5 Monitoramento
- Métricas em tempo real
- Alertas configuráveis
- Dashboard de status
- Logs estruturados
- Performance tracking

### 2.6 Integração Cloudflare
- Criação automática de subdomínios
- Gerenciamento de DNS
- SSL/TLS automático
- Cache e CDN
- Firewall rules

---

## 3. Modelos de Dados

### 3.1 Cliente
```go
type Client struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    Limits      Limits    `json:"limits"`
    Metadata    map[string]string `json:"metadata"`
}
```

### 3.2 Serviço
```go
type Service struct {
    ID          string    `json:"id"`
    ClientID    string    `json:"client_id"`
    Type        string    `json:"type"` // n8n, wordpress, custom
    Name        string    `json:"name"`
    Subdomain   string    `json:"subdomain"`
    Status      string    `json:"status"`
    Resources   Resources `json:"resources"` // CPU, Memória, Disco, etc.
    Config      map[string]interface{} `json:"config"`
    ContainerID string    `json:"container_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Recursos alocados para o serviço
type Resources struct {
    CPU       CPUResources    `json:"cpu"`
    Memory    MemoryResources `json:"memory"`
    Disk      DiskResources   `json:"disk"`
    Network   NetworkResources `json:"network"`
    Limits    Limits          `json:"limits"`
}

// CPU: cores e limites
type CPUResources struct {
    Cores      float64 `json:"cores"`       // Quantidade de cores CPU (ex: 0.5, 1.0, 2.0)
    MaxCores   float64 `json:"max_cores"`   // Limite máximo de cores
    Priority   int     `json:"priority"`    // Prioridade CPU (0-1024)
}

// Memória: RAM e swap
type MemoryResources struct {
    Limit      string  `json:"limit"`       // Limite de memória (ex: "512m", "1g", "2g")
    Reservation string `json:"reservation"` // Reserva garantida (ex: "256m")
    Swap        string  `json:"swap"`       // Limite de swap (ex: "512m", "-1" para ilimitado)
}

// Disco: armazenamento persistente
type DiskResources struct {
    Size        string  `json:"size"`        // Tamanho do volume (ex: "10g", "50g", "100g")
    VolumeType  string  `json:"volume_type"` // Tipo de volume (standard, ssd, nvme)
    IOPS        int     `json:"iops"`        // IOPS para discos SSD/NVMe
    Throughput  int     `json:"throughput"`  // Throughput em MB/s
}

// Rede: bandwidth e conexões
type NetworkResources struct {
    Bandwidth   string  `json:"bandwidth"`   // Largura de banda (ex: "100m", "1g")
    MaxConnections int  `json:"max_connections"` // Máximo de conexões simultâneas
}

// Limites adicionais
type Limits struct {
    Pids        int     `json:"pids"`        // Máximo de processos
    Ulimits     []ULimit `json:"ulimits"`   // Limites do sistema (nofile, nproc, etc.)
}

type ULimit struct {
    Name  string `json:"name"`
    Soft  int    `json:"soft"`
    Hard  int    `json:"hard"`
}
```

### 3.3 Instalação
```go
type Installation struct {
    ID          string    `json:"id"`
    ServiceID   string    `json:"service_id"`
    Version     string    `json:"version"`
    Status      string    `json:"status"`
    Resources   Resources `json:"resources"` // Recursos solicitados na instalação
    Logs        []string  `json:"logs"`
    StartedAt   time.Time `json:"started_at"`
    CompletedAt *time.Time `json:"completed_at"`
}
```

### 3.4 Recursos e Limites do Sistema
```go
// Recursos disponíveis no host
type HostResources struct {
    TotalCPU    float64 `json:"total_cpu"`    // Total de CPUs disponíveis
    TotalMemory int64   `json:"total_memory"`  // Total de memória em bytes
    TotalDisk   int64   `json:"total_disk"`    // Total de disco em bytes
    UsedCPU     float64 `json:"used_cpu"`      // CPU em uso
    UsedMemory  int64   `json:"used_memory"`   // Memória em uso
    UsedDisk    int64   `json:"used_disk"`     // Disco em uso
    AvailableCPU float64 `json:"available_cpu"` // CPU disponível
    AvailableMemory int64 `json:"available_memory"` // Memória disponível
    AvailableDisk int64 `json:"available_disk"` // Disco disponível
}
```

---

## 4. APIs REST

### 4.1 Clientes
```
POST   /api/v1/clients              - Criar cliente
GET    /api/v1/clients              - Listar clientes
GET    /api/v1/clients/{id}         - Obter cliente
PUT    /api/v1/clients/{id}         - Atualizar cliente
DELETE /api/v1/clients/{id}         - Deletar cliente
```

### 4.2 Serviços
```
POST   /api/v1/services             - Criar serviço (com recursos)
GET    /api/v1/services             - Listar serviços
GET    /api/v1/services/{id}        - Obter serviço
PUT    /api/v1/services/{id}        - Atualizar serviço (incluindo recursos)
DELETE /api/v1/services/{id}        - Deletar serviço
POST   /api/v1/services/{id}/start  - Iniciar serviço
POST   /api/v1/services/{id}/stop   - Parar serviço
POST   /api/v1/services/{id}/restart - Reiniciar serviço
PUT    /api/v1/services/{id}/resources - Atualizar recursos do serviço
GET    /api/v1/services/{id}/resources/usage - Uso atual de recursos
```

#### Exemplo de Payload para Criação de Serviço
```json
{
  "client_id": "client-123",
  "type": "n8n",
  "name": "meu-n8n",
  "subdomain": "automation",
  "resources": {
    "cpu": {
      "cores": 1.0,
      "max_cores": 2.0,
      "priority": 512
    },
    "memory": {
      "limit": "2g",
      "reservation": "1g",
      "swap": "512m"
    },
    "disk": {
      "size": "20g",
      "volume_type": "ssd",
      "iops": 3000,
      "throughput": 500
    },
    "network": {
      "bandwidth": "100m",
      "max_connections": 1000
    },
    "limits": {
      "pids": 100,
      "ulimits": [
        {
          "name": "nofile",
          "soft": 65536,
          "hard": 65536
        }
      ]
    }
  },
  "config": {
    "username": "admin",
    "password": "secure_password"
  }
}
```

### 4.3 Instalações
```
POST   /api/v1/installations        - Iniciar instalação
GET    /api/v1/installations/{id}   - Status da instalação
GET    /api/v1/installations/{id}/logs - Logs da instalação
```

### 4.4 Monitoramento
```
GET    /api/v1/metrics              - Métricas agregadas
GET    /api/v1/services/{id}/metrics - Métricas do serviço
GET    /api/v1/services/{id}/logs   - Logs do serviço
GET    /api/v1/health               - Health check
GET    /api/v1/resources/available - Recursos disponíveis no host
GET    /api/v1/resources/usage     - Uso de recursos por cliente/serviço
```

### 4.5 Recursos
```
GET    /api/v1/resources/available - Recursos disponíveis no sistema
GET    /api/v1/resources/usage     - Uso agregado de recursos
GET    /api/v1/clients/{id}/resources - Recursos utilizados por cliente
```

---

## 5. Implementação Cloudflare

### 5.1 Criação de Subdomínios
```go
// Estrutura para integração Cloudflare
type CloudflareService struct {
    APIKey    string
    Email     string
    ZoneID    string
    BaseURL   string
}

func (c *CloudflareService) CreateSubdomain(subdomain string, target string) error {
    // Implementar chamada à API Cloudflare
    // POST /zones/{zone_id}/dns_records
}
```

### 5.2 Endpoints Cloudflare
```
POST   /api/v1/cloudflare/subdomains  - Criar subdomínio
GET    /api/v1/cloudflare/subdomains  - Listar subdomínios
DELETE /api/v1/cloudflare/subdomains/{id} - Deletar subdomínio
PUT    /api/v1/cloudflare/subdomains/{id} - Atualizar subdomínio
```

---

## 6. Gerenciamento Docker

### 6.1 Definição de Recursos para Serviços

Cada serviço criado **DEVE** ter recursos definidos obrigatoriamente:

#### Recursos Obrigatórios
- **CPU**: Número de cores (podem ser frações, ex: 0.5, 1.0, 2.0)
- **Memória**: Limite de RAM (ex: "512m", "1g", "2g")
- **Disco**: Tamanho do volume persistente (ex: "10g", "50g", "100g")

#### Recursos Opcionais mas Recomendados
- **CPU Priority**: Prioridade de CPU (shares)
- **Memory Reservation**: Memória garantida mínima
- **Memory Swap**: Limite de swap
- **Disk Type**: Tipo de disco (standard, ssd, nvme)
- **Disk IOPS**: IOPS para discos SSD/NVMe
- **Network Bandwidth**: Largura de banda de rede
- **Max Connections**: Máximo de conexões simultâneas
- **Process Limits**: Máximo de processos (pids)
- **Ulimits**: Limites do sistema (nofile, nproc, etc.)

### 6.2 Templates de Serviços

#### n8n
```yaml
version: '3.8'
services:
  n8n:
    image: n8nio/n8n:latest
    container_name: n8n_{service_id}
    environment:
      - N8N_BASIC_AUTH_ACTIVE=true
      - N8N_BASIC_AUTH_USER={username}
      - N8N_BASIC_AUTH_PASSWORD={password}
      - N8N_HOST={subdomain}
      - N8N_PROTOCOL=https
    # Recursos de CPU
    cpus: '{cpu_cores}'           # Ex: 1.0, 2.0
    cpu_count: '{cpu_count}'      # Ex: 1, 2
    cpu_percent: '{cpu_percent}'  # Ex: 50 (50% de 1 core)
    cpu_quota: '{cpu_quota}'      # Ex: 50000 (microseconds)
    cpu_period: '{cpu_period}'    # Ex: 100000
    cpu_shares: '{cpu_priority}'  # Ex: 512
    
    # Recursos de Memória
    mem_limit: '{memory_limit}'        # Ex: 2g
    mem_reservation: '{memory_reservation}' # Ex: 1g
    memswap_limit: '{memory_swap}'     # Ex: 512m
    
    # Recursos de Disco
    storage_opt:
      size: '{disk_size}'              # Ex: 20g
    
    # Limites de Processos
    pids_limit: '{pids_limit}'         # Ex: 100
    
    # Ulimits
    ulimits:
      nofile:
        soft: {ulimit_nofile_soft}     # Ex: 65536
        hard: {ulimit_nofile_hard}     # Ex: 65536
      nproc:
        soft: {ulimit_nproc_soft}      # Ex: 32768
        hard: {ulimit_nproc_hard}      # Ex: 32768
    
    volumes:
      - n8n_data_{service_id}:/home/node/.n8n
    networks:
      - service_network
    deploy:
      resources:
        limits:
          cpus: '{cpu_cores}'
          memory: '{memory_limit}'
        reservations:
          cpus: '{cpu_reservation}'
          memory: '{memory_reservation}'
```

#### WordPress
```yaml
version: '3.8'
services:
  wordpress:
    image: wordpress:latest
    container_name: wordpress_{service_id}
    environment:
      WORDPRESS_DB_HOST: db_{service_id}
      WORDPRESS_DB_USER: {db_user}
      WORDPRESS_DB_PASSWORD: {db_password}
      WORDPRESS_DB_NAME: {db_name}
    
    # Recursos de CPU
    cpus: '{wp_cpu_cores}'
    cpu_shares: '{wp_cpu_priority}'
    
    # Recursos de Memória
    mem_limit: '{wp_memory_limit}'
    mem_reservation: '{wp_memory_reservation}'
    
    # Recursos de Disco
    storage_opt:
      size: '{wp_disk_size}'
    
    volumes:
      - wordpress_data_{service_id}:/var/www/html
    networks:
      - service_network
    depends_on:
      - db_{service_id}
    deploy:
      resources:
        limits:
          cpus: '{wp_cpu_cores}'
          memory: '{wp_memory_limit}'
        reservations:
          memory: '{wp_memory_reservation}'
  
  db_{service_id}:
    image: mysql:8.0
    container_name: db_{service_id}
    environment:
      MYSQL_DATABASE: {db_name}
      MYSQL_USER: {db_user}
      MYSQL_PASSWORD: {db_password}
    
    # Recursos de CPU para MySQL
    cpus: '{db_cpu_cores}'
    cpu_shares: '{db_cpu_priority}'
    
    # Recursos de Memória para MySQL
    mem_limit: '{db_memory_limit}'
    mem_reservation: '{db_memory_reservation}'
    
    # Recursos de Disco para MySQL
    storage_opt:
      size: '{db_disk_size}'
    
    # IOPS para banco de dados (se suportado)
    storage_opt:
      size: '{db_disk_size}'
    
    volumes:
      - db_data_{service_id}:/var/lib/mysql
    deploy:
      resources:
        limits:
          cpus: '{db_cpu_cores}'
          memory: '{db_memory_limit}'
        reservations:
          memory: '{db_memory_reservation}'
```

### 6.3 Implementação Go para Gerenciamento de Recursos

```go
// Service para gerenciar recursos de containers
type ResourceManager struct {
    dockerClient *docker.Client
    db          *database.DB
}

// Criar container com recursos definidos
func (rm *ResourceManager) CreateContainerWithResources(
    serviceID string,
    resources Resources,
    image string,
    env []string,
) (*container.ContainerCreateCreatedBody, error) {
    
    // Converter recursos para configuração Docker
    hostConfig := &container.HostConfig{
        // CPU
        CPUCount:    int64(resources.CPU.Cores),
        CPUShares:   int64(resources.CPU.Priority),
        NanoCPUs:    int64(resources.CPU.Cores * 1e9),
        
        // Memória
        Memory:      parseMemoryToBytes(resources.Memory.Limit),
        MemoryReservation: parseMemoryToBytes(resources.Memory.Reservation),
        MemorySwap:  parseMemoryToBytes(resources.Memory.Swap),
        
        // Limites de processos
        PidsLimit:   int64(resources.Limits.Pids),
        
        // Ulimits
        Ulimits: convertUlimits(resources.Limits.Ulimits),
    }
    
    // Criar volume com tamanho específico
    volumeName := fmt.Sprintf("%s_data_%s", serviceID, generateID())
    volume, err := rm.dockerClient.VolumeCreate(context.Background(), volume.VolumeCreateBody{
        Name: volumeName,
        DriverOpts: map[string]string{
            "size": resources.Disk.Size,
            "type": resources.Disk.VolumeType,
        },
    })
    if err != nil {
        return nil, err
    }
    
    // Configuração do container
    config := &container.Config{
        Image: image,
        Env:   env,
    }
    
    // Criar container
    createResp, err := rm.dockerClient.ContainerCreate(
        context.Background(),
        config,
        hostConfig,
        nil,
        nil,
        fmt.Sprintf("%s_%s", serviceID, generateID()),
    )
    
    return &createResp, err
}

// Atualizar recursos de um container em execução
func (rm *ResourceManager) UpdateContainerResources(
    containerID string,
    resources Resources,
) error {
    
    updateConfig := container.UpdateConfig{
        CPUCount:    int64(resources.CPU.Cores),
        CPUShares:   int64(resources.CPU.Priority),
        NanoCPUs:    int64(resources.CPU.Cores * 1e9),
        Memory:      parseMemoryToBytes(resources.Memory.Limit),
        MemoryReservation: parseMemoryToBytes(resources.Memory.Reservation),
        MemorySwap:  parseMemoryToBytes(resources.Memory.Swap),
        PidsLimit:   int64(resources.Limits.Pids),
        Ulimits:     convertUlimits(resources.Limits.Ulimits),
    }
    
    _, err := rm.dockerClient.ContainerUpdate(
        context.Background(),
        containerID,
        updateConfig,
    )
    
    return err
}

// Expandir volume de disco
func (rm *ResourceManager) ExpandVolume(volumeName string, newSize string) error {
    // Implementar expansão de volume conforme driver de storage
    // Para Docker volumes locais, pode ser necessário usar comandos do sistema
    // Para volumes em cloud, usar APIs específicas
    return nil
}

// Verificar recursos disponíveis no host
func (rm *ResourceManager) GetAvailableResources() (*HostResources, error) {
    info, err := rm.dockerClient.Info(context.Background())
    if err != nil {
        return nil, err
    }
    
    // Obter recursos utilizados por todos os containers
    containers, err := rm.dockerClient.ContainerList(
        context.Background(),
        types.ContainerListOptions{All: true},
    )
    
    var usedCPU float64
    var usedMemory int64
    
    for _, c := range containers {
        stats, err := rm.dockerClient.ContainerStats(
            context.Background(),
            c.ID,
            false,
        )
        if err != nil {
            continue
        }
        
        // Calcular uso de recursos (simplificado)
        // Implementar lógica completa de cálculo
    }
    
    return &HostResources{
        TotalCPU:    float64(info.NCPU),
        TotalMemory: info.MemTotal,
        TotalDisk:   getTotalDisk(),
        UsedCPU:     usedCPU,
        UsedMemory:  usedMemory,
        UsedDisk:    getUsedDisk(),
        AvailableCPU: float64(info.NCPU) - usedCPU,
        AvailableMemory: info.MemTotal - usedMemory,
        AvailableDisk: getTotalDisk() - getUsedDisk(),
    }, nil
}
```

---

## 7. Fluxo de Criação de Serviço

1. **Receber requisição** via API REST com recursos especificados
2. **Validar recursos solicitados**:
   - Verificar disponibilidade de CPU, Memória e Disco no host
   - Validar limites do cliente (quota de recursos)
   - Verificar se recursos atendem requisitos mínimos do tipo de serviço
   - Validar formato e valores dos recursos (cores, memória, disco)
3. **Validar cliente** e limites totais
4. **Reservar recursos** no sistema (marcar como alocados)
5. **Gerar subdomínio** único
6. **Criar registro DNS** via Cloudflare API
7. **Gerar configuração** Docker Compose com recursos aplicados
8. **Criar volumes** com tamanho especificado
9. **Iniciar containers** via Docker API com limites de recursos
10. **Aplicar configurações de recursos** via Docker API:
    - CPU limits e shares
    - Memory limits e reservations
    - Disk quotas
    - Network bandwidth
    - Process limits (pids)
    - Ulimits
11. **Aguardar saúde** do serviço
12. **Validar recursos aplicados** (verificar se containers respeitam limites)
13. **Configurar SSL** automático (Cloudflare)
14. **Registrar no banco** de dados com recursos alocados
15. **Iniciar monitoramento** de recursos (CPU, memória, disco)
16. **Retornar informações** do serviço incluindo recursos configurados

### 7.1 Validação de Recursos

#### Requisitos Mínimos por Tipo de Serviço
```go
var ServiceMinRequirements = map[string]Resources{
    "n8n": {
        CPU: CPUResources{Cores: 0.5, MaxCores: 1.0},
        Memory: MemoryResources{Limit: "512m", Reservation: "256m"},
        Disk: DiskResources{Size: "5g", VolumeType: "standard"},
    },
    "wordpress": {
        CPU: CPUResources{Cores: 0.5, MaxCores: 1.0},
        Memory: MemoryResources{Limit: "512m", Reservation: "256m"},
        Disk: DiskResources{Size: "10g", VolumeType: "standard"},
    },
    "custom": {
        CPU: CPUResources{Cores: 0.25, MaxCores: 0.5},
        Memory: MemoryResources{Limit: "256m", Reservation: "128m"},
        Disk: DiskResources{Size: "1g", VolumeType: "standard"},
    },
}
```

#### Validação de Disponibilidade
```go
func ValidateResourceAvailability(requested Resources, available HostResources) error {
    // Validar CPU
    if requested.CPU.Cores > available.AvailableCPU {
        return fmt.Errorf("CPU insuficiente: solicitado %.2f, disponível %.2f", 
            requested.CPU.Cores, available.AvailableCPU)
    }
    
    // Validar Memória
    requestedMemBytes := parseMemory(requested.Memory.Limit)
    if requestedMemBytes > available.AvailableMemory {
        return fmt.Errorf("Memória insuficiente: solicitado %d, disponível %d", 
            requestedMemBytes, available.AvailableMemory)
    }
    
    // Validar Disco
    requestedDiskBytes := parseDisk(requested.Disk.Size)
    if requestedDiskBytes > available.AvailableDisk {
        return fmt.Errorf("Disco insuficiente: solicitado %d, disponível %d", 
            requestedDiskBytes, available.AvailableDisk)
    }
    
    return nil
}
```

---

## 8. Monitoramento e Alertas

### 8.1 Métricas Coletadas
- **CPU**:
  - CPU usage por serviço (percentual e cores)
  - CPU throttling
  - CPU shares utilizados
- **Memória**:
  - Memória utilizada vs. limite
  - Memória reservada vs. utilizada
  - Swap utilizado
  - OOM (Out of Memory) events
- **Disco**:
  - Espaço utilizado vs. alocado
  - IOPS utilizados
  - Throughput de leitura/escrita
  - Taxa de crescimento de disco
- **Rede**:
  - Bandwidth utilizado
  - Conexões ativas vs. máximo
  - Latência de rede
- **Outras**:
  - Requests por segundo
  - Tempo de resposta
  - Taxa de erro
  - Uptime/Downtime
  - Número de processos ativos

### 8.2 Alertas Configuráveis
- **CPU**:
  - CPU > 80% do limite alocado
  - CPU throttling detectado
  - CPU > 100% por período prolongado
- **Memória**:
  - Memória > 90% do limite
  - Memória > 95% do limite (crítico)
  - OOM (Out of Memory) events
  - Swap usage > 50%
- **Disco**:
  - Disco > 85% do espaço alocado
  - Disco > 90% do espaço alocado (crítico)
  - IOPS próximo ao limite
  - Throughput próximo ao limite
- **Rede**:
  - Bandwidth > 80% do limite
  - Conexões > 90% do máximo
- **Serviço**:
  - Serviço offline
  - Erros > threshold
  - Latência alta
  - Health check falhando
- **Recursos**:
  - Recursos insuficientes disponíveis no host
  - Cliente próximo ao limite de quota

---

## 9. Segurança

- Autenticação JWT
- Autorização baseada em roles
- Isolamento de containers por cliente
- Network isolation
- Secrets management (Vault)
- Rate limiting
- HTTPS obrigatório
- Validação de inputs

---

## 10. Backup e Disaster Recovery

- Backup automático de dados
- Backup de configurações
- Snapshot de volumes Docker
- Restore point-in-time
- Backup off-site
- Testes de restore periódicos

---

## 11. Escalabilidade

- Horizontal scaling de serviços
- Load balancing
- Auto-scaling baseado em métricas
- Queue system para tarefas assíncronas
- Distributed locking
- Service discovery

---

## 12. Roadmap de Implementação

### Fase 1: Fundação (Semanas 1-2)
- [ ] Estrutura do projeto Go
- [ ] Configuração de banco de dados
- [ ] Modelos básicos (Cliente, Serviço)
- [ ] API REST básica
- [ ] Autenticação JWT

### Fase 2: Docker Integration (Semanas 3-4)
- [ ] Integração com Docker API
- [ ] Templates de serviços (n8n, WordPress)
- [ ] **Sistema de recursos**:
  - [ ] Modelos de recursos (CPU, Memória, Disco)
  - [ ] Validação de recursos disponíveis
  - [ ] Aplicação de limites de recursos em containers
  - [ ] Reserva de recursos no sistema
- [ ] Criação de containers dinâmica com recursos
- [ ] Gerenciamento de volumes com tamanho específico
- [ ] Network management com bandwidth limits

### Fase 3: Cloudflare Integration (Semanas 5-6)
- [ ] SDK Cloudflare Go
- [ ] Criação de subdomínios
- [ ] Gerenciamento de DNS
- [ ] SSL automático
- [ ] Cache configuration

### Fase 4: Monitoramento (Semanas 7-8)
- [ ] Coleta de métricas
- [ ] **Monitoramento de recursos**:
  - [ ] Métricas de CPU em tempo real
  - [ ] Métricas de memória em tempo real
  - [ ] Métricas de disco em tempo real
  - [ ] Métricas de rede em tempo real
  - [ ] Histórico de uso de recursos
- [ ] Integração Prometheus
- [ ] Dashboards Grafana com recursos
- [ ] Sistema de alertas baseado em recursos
- [ ] Log aggregation
- [ ] **APIs para consulta de recursos disponíveis**

### Fase 5: Gestão Completa (Semanas 9-10)
- [ ] CRUD completo de clientes
- [ ] CRUD completo de serviços
- [ ] Backup automático
- [ ] Restore functionality
- [ ] Documentação API

### Fase 6: Otimização (Semanas 11-12)
- [ ] Performance tuning
- [ ] Cache implementation
- [ ] Queue system
- [ ] Auto-scaling
- [ ] Load testing

---

## 13. Comandos Principais

### Instalação
```bash
# Clone e build
git clone <repo>
cd project
go mod download
go build -o bin/server cmd/server/main.go

# Configurar variáveis de ambiente
cp configs/config.example.yaml configs/config.yaml
# Editar config.yaml com suas credenciais

# Executar migrations
./bin/server migrate

# Iniciar servidor
./bin/server serve
```

### Deploy
```bash
# Via Docker Compose
docker-compose up -d

# Via Kubernetes
kubectl apply -f deployments/kubernetes/
```

---

## 14. Variáveis de Ambiente

```yaml
# Database
DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: password
DB_NAME: services_db

# Cloudflare
CLOUDFLARE_API_KEY: your_api_key
CLOUDFLARE_EMAIL: your_email
CLOUDFLARE_ZONE_ID: your_zone_id
CLOUDFLARE_BASE_DOMAIN: example.com

# Docker
DOCKER_HOST: unix:///var/run/docker.sock

# API
API_PORT: 8080
API_SECRET_KEY: your_secret_key

# Redis
REDIS_HOST: localhost
REDIS_PORT: 6379
```

---

## 15. Próximos Passos

1. **Setup inicial**: Configurar ambiente de desenvolvimento
2. **Protótipo**: Implementar POC com um serviço (n8n)
3. **Testes**: Testes unitários e de integração
4. **Documentação**: Documentar APIs e processos
5. **Deploy**: Deploy em ambiente de staging
6. **Produção**: Deploy em produção com monitoramento

---

## 16. Recursos Adicionais

- **Documentação API**: Swagger/OpenAPI
- **Client SDK**: SDK em Go para consumir a API
- **Webhooks**: Notificações de eventos
- **Audit Log**: Log de todas as operações
- **Multi-tenancy**: Suporte completo multi-tenant
- **Quotas**: Limites por cliente/serviço
- **Billing**: Integração com sistema de cobrança

---

## Conclusão

Este roadmap fornece uma visão completa do sistema de gerenciamento dinâmico de serviços. A implementação deve ser feita de forma incremental, começando pelas funcionalidades core e expandindo gradualmente conforme a necessidade.
