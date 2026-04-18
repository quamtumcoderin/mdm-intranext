package main

import (
	"context"
	"log"
	"os"

	"mdm-intranext/pkg/database"
	"mdm-intranext/pkg/vault"
)

func main() {
	ctx := context.Background()

	vaultCfg := vault.Config{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}

	secretManager, err := vault.NewSecretManager(vaultCfg)
	if err != nil {
		log.Fatalf("Failed to create vault secret manager: %v", err)
	}

	dbUser, dbPass, dbName, err := secretManager.GetDatabaseCredentials(ctx)
	if err != nil {
		log.Fatalf("Failed to get database credentials: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbPool, err := database.NewPostgresPool(ctx, dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Backend inicializado: Conexión segura a Vault y PostgreSQL establecida.")
}
