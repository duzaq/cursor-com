#!/bin/bash

# Script de instala??o e setup inicial

set -e

echo "?? Configurando Sistema de Gerenciamento de Servi?os..."

# Verificar se Go est? instalado
if ! command -v go &> /dev/null; then
    echo "? Go n?o est? instalado. Por favor, instale Go 1.21+"
    exit 1
fi

echo "? Go encontrado: $(go version)"

# Instalar depend?ncias
echo "?? Instalando depend?ncias..."
go mod download

# Verificar se Docker est? instalado
if ! command -v docker &> /dev/null; then
    echo "??  Docker n?o est? instalado. Algumas funcionalidades podem n?o funcionar."
else
    echo "? Docker encontrado: $(docker --version)"
fi

# Verificar se Docker Compose est? instalado
if ! command -v docker-compose &> /dev/null; then
    echo "??  Docker Compose n?o est? instalado."
else
    echo "? Docker Compose encontrado: $(docker-compose --version)"
fi

# Criar arquivo de configura??o se n?o existir
if [ ! -f "configs/config.local.yaml" ]; then
    echo "?? Criando arquivo de configura??o local..."
    cp configs/config.yaml configs/config.local.yaml
    echo "??  Por favor, edite configs/config.local.yaml com suas configura??es"
fi

# Iniciar banco de dados
if command -v docker-compose &> /dev/null; then
    echo "?? Iniciando PostgreSQL e Redis..."
    docker-compose up -d postgres redis
    
    echo "? Aguardando banco de dados ficar pronto..."
    sleep 5
fi

# Build do projeto
echo "?? Compilando projeto..."
go build -o bin/server cmd/server/main.go

echo ""
echo "? Setup conclu?do!"
echo ""
echo "Pr?ximos passos:"
echo "1. Edite configs/config.local.yaml com suas configura??es"
echo "2. Execute: ./bin/server"
echo "3. Ou use: go run cmd/server/main.go"
echo ""
echo "Para parar os containers: docker-compose down"
