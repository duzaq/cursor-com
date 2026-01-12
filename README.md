# Servi?o de Gerenciamento Din?mico de Servi?os

Sistema desenvolvido em Go para gerenciar dinamicamente m?ltiplos servi?os (n8n, WordPress, outras aplica??es) com integra??o Cloudflare para cria??o autom?tica de subdom?nios.

## Requisitos

- Go 1.21+
- Docker e Docker Compose
- PostgreSQL 15+
- Redis (opcional, para cache)

## Instala??o

1. Clone o reposit?rio:
```bash
git clone <repo>
cd services-manager
```

2. Instale as depend?ncias:
```bash
go mod download
```

3. Configure as vari?veis de ambiente:
```bash
cp configs/config.yaml configs/config.local.yaml
# Edite configs/config.local.yaml com suas configura??es
```

4. Inicie o banco de dados:
```bash
docker-compose up -d postgres redis
```

5. Execute as migrations:
```bash
go run cmd/server/main.go migrate
```

6. Inicie o servidor:
```bash
go run cmd/server/main.go serve
```

Ou compile e execute:
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## Vari?veis de Ambiente

Veja `configs/config.yaml` para todas as op??es de configura??o.

Principais vari?veis:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `CLOUDFLARE_API_KEY`, `CLOUDFLARE_EMAIL`, `CLOUDFLARE_ZONE_ID`, `CLOUDFLARE_BASE_DOMAIN`
- `DOCKER_HOST`
- `API_PORT`, `API_SECRET_KEY`, `JWT_SECRET`

## APIs

### Clientes

- `POST /api/v1/clients` - Criar cliente
- `GET /api/v1/clients` - Listar clientes
- `GET /api/v1/clients/:id` - Obter cliente
- `PUT /api/v1/clients/:id` - Atualizar cliente
- `DELETE /api/v1/clients/:id` - Deletar cliente

### Servi?os

- `POST /api/v1/services` - Criar servi?o
- `GET /api/v1/services` - Listar servi?os
- `GET /api/v1/services/:id` - Obter servi?o
- `PUT /api/v1/services/:id` - Atualizar servi?o
- `DELETE /api/v1/services/:id` - Deletar servi?o
- `PUT /api/v1/services/:id/resources` - Atualizar recursos
- `GET /api/v1/services/:id/resources/usage` - Uso de recursos
- `POST /api/v1/services/:id/start` - Iniciar servi?o
- `POST /api/v1/services/:id/stop` - Parar servi?o
- `POST /api/v1/services/:id/restart` - Reiniciar servi?o

### Instala??es

- `POST /api/v1/installations` - Criar instala??o
- `GET /api/v1/installations/:id` - Obter instala??o
- `GET /api/v1/installations/:id/logs` - Logs da instala??o

### Recursos

- `GET /api/v1/resources/available` - Recursos dispon?veis
- `GET /api/v1/resources/usage` - Uso de recursos
- `GET /api/v1/clients/:id/resources` - Recursos do cliente

## Desenvolvimento

Para desenvolvimento local com hot reload:
```bash
# Instalar air (hot reload)
go install github.com/cosmtrek/air@latest

# Executar com hot reload
air
```

## Estrutura do Projeto

```
.
??? cmd/
?   ??? server/
?       ??? main.go          # Ponto de entrada
??? internal/
?   ??? api/                 # API REST
?   ?   ??? handlers/        # Handlers HTTP
?   ?   ??? server.go        # Configura??o do servidor
?   ??? config/              # Configura??o
?   ??? database/            # Camada de dados
?   ?   ??? client_repository.go
?   ?   ??? service_repository.go
?   ?   ??? installation_repository.go
?   ??? models/              # Modelos de dados
??? configs/                 # Arquivos de configura??o
??? docker-compose.yml       # Docker Compose para desenvolvimento
??? go.mod                   # Depend?ncias Go
```

## Pr?ximos Passos

Consulte o arquivo `roadmap.md` para ver o roadmap completo de implementa??o.
