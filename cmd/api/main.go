package main

import (
	"context"
	"errors"
	"log"
	httphandler "mdm-intranext/internal/adapters/handler/http"
	"mdm-intranext/internal/adapters/repository/postgres"
	"mdm-intranext/internal/core/services"
	"net/http"
	"os"
	"time"

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

	deviceRepo := postgres.NewDeviceRepository(dbPool)
	log.Println("Connected to database")

	deviceSvc := services.NewDeviceService(deviceRepo)
	log.Println("Connected to device service")

	mux := http.NewServeMux()
	deviceHandler := httphandler.NewDeviceHandler(deviceSvc)
	deviceHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    		":" + os.Getenv("PORT"),
		Handler: 		mux,
		ReadTimeout:  	5 * time.Second,
		WriteTimeout: 	10 * time.Second,
		IdleTimeout:  	120 * time.Second,
	}

	log.Println("Listening on port " + os.Getenv("PORT"))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Backend inicializado: Conexión segura a Vault y PostgreSQL establecida.")
}
