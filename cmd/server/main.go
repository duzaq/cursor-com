package main

import (
	"flag"
	"log"
	"os"

	"github.com/services-manager/internal/api"
	"github.com/services-manager/internal/config"
	"github.com/services-manager/internal/database"
)

func main() {
	// Parse de argumentos da linha de comando
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = os.Args[1:]
	}

	flag.Parse()

	// Carregar configura??o
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configura??o: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Executar comando
	switch command {
	case "migrate":
		log.Println("Executando migrations...")
		if err := database.Migrate(db); err != nil {
			log.Fatalf("Erro ao executar migrations: %v", err)
		}
		log.Println("Migrations executadas com sucesso!")
		return

	case "serve", "server", "":
		// Executar migrations antes de iniciar servidor
		if err := database.Migrate(db); err != nil {
			log.Printf("Aviso: Erro ao executar migrations: %v", err)
		}

		// Criar servidor HTTP
		server := api.NewServer(cfg, db)

		// Iniciar servidor
		port := cfg.API.Port
		if port == "" {
			port = "8080"
		}

		log.Printf("Servidor iniciado na porta %s", port)
		if err := server.Start(":" + port); err != nil {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}

	default:
		log.Fatalf("Comando desconhecido: %s. Use 'migrate' ou 'serve'", command)
	}
}
