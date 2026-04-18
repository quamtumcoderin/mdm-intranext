# Conexión del Backend a PostgreSQL y Vault

## 1. Introducción

Para iniciar el desarrollo del backend concurrente en Go, se adopta una arquitectura basada en **Clean Architecture (Arquitectura Hexagonal)**. Este enfoque permite una separación estricta de responsabilidades, aislando la lógica de negocio de los detalles de infraestructura como PostgreSQL, MQTT y Vault.

El objetivo principal es garantizar mantenibilidad, testabilidad y bajo acoplamiento entre componentes.

---

## 2. Inicialización del módulo

Desde la raíz del proyecto `mdm-intranext`, se crea el directorio del backend e inicializa el módulo de Go:

```bash
mkdir -p backend
cd backend
go mod init mdm-intranext
```

Esto establece el módulo base de Go para gestionar dependencias y estructura del proyecto.

---

## 3. Estructura de directorios

Se define una estructura modular basada en capas (dominio, puertos y adaptadores), siguiendo principios de arquitectura hexagonal:

```bash
mkdir -p cmd/api \
internal/core/domain \
internal/core/ports \
internal/core/services \
internal/adapters/repository/postgres \
internal/adapters/handler/http \
internal/adapters/handler/mqtt \
pkg/vault \
pkg/database
```

### Descripción de capas

* **cmd/api**: Punto de entrada de la aplicación.
* **internal/core/domain**: Entidades del dominio.
* **internal/core/ports**: Interfaces que definen contratos del sistema.
* **internal/core/services**: Lógica de negocio.
* **internal/adapters/**: Implementaciones concretas (HTTP, MQTT, PostgreSQL).
* **pkg/vault**: Cliente de Vault para gestión de secretos.
* **pkg/database**: Configuración de conexiones a bases de datos.

---

## 4. Dependencias de infraestructura

Se instalan las librerías necesarias para integración con Vault y PostgreSQL:

```bash
go get github.com/hashicorp/vault/api
go get github.com/jackc/pgx/v5/pgxpool
```

* **Vault API**: Para la gestión segura de secretos.
* **pgxpool**: Pool de conexiones eficiente para PostgreSQL.

---

## 5. Cliente de Vault

Archivo: `pkg/vault/client.go`

Este componente encapsula la conexión con Vault y la recuperación de credenciales de base de datos.

```go
package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type Config struct {
	Address string
	Token   string
}

type SecretManager struct {
	client *vault.Client
}

func NewSecretManager(cfg Config) (*SecretManager, error) {
	config := vault.DefaultConfig()
	config.Address = cfg.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error inicializando cliente vault: %w", err)
	}

	client.SetToken(cfg.Token)

	return &SecretManager{client: client}, nil
}

func (s *SecretManager) GetDatabaseCredentials(ctx context.Context) (string, string, string, error) {
	secret, err := s.client.KVv2("secret").Get(ctx, "mdm/database")
	if err != nil {
		return "", "", "", fmt.Errorf("error leyendo secreto de base de datos: %w", err)
	}

	user := secret.Data["username"].(string)
	pass := secret.Data["password"].(string)
	dbName := secret.Data["dbname"].(string)

	return user, pass, dbName, nil
}
```

### Responsabilidad

* Autenticación contra Vault.
* Recuperación de credenciales de base de datos.
* Eliminación de dependencia de variables de entorno para secretos sensibles.

---

## 6. Cliente de PostgreSQL

Archivo: `pkg/database/postgres.go`

Este módulo gestiona el pool de conexiones optimizado a PostgreSQL usando `pgxpool`.

```go
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, host, port, user, pass, dbname string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Configuración de concurrencia
	config.MaxConns = 20
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
```

### Características

* Pool de conexiones reutilizable.
* Configuración de concurrencia controlada.
* Validación de conectividad mediante `Ping`.

---

## 7. Punto de entrada del sistema

Archivo: `cmd/api/main.go`

Este archivo orquesta la inicialización del backend, conectando Vault y PostgreSQL de forma segura.

```go
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

	// 1. Inicializar Vault
	vaultCfg := vault.Config{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}

	secretManager, err := vault.NewSecretManager(vaultCfg)
	if err != nil {
		log.Fatalf("Fallo crítico: no se pudo conectar a Vault: %v", err)
	}

	// 2. Obtener credenciales seguras
	dbUser, dbPass, dbName, err := secretManager.GetDatabaseCredentials(ctx)
	if err != nil {
		log.Fatalf("Fallo crítico: no se pudieron obtener credenciales: %v", err)
	}

	// 3. Inicializar PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := "5432"

	dbPool, err := database.NewPostgresPool(ctx, dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Fallo crítico: no se pudo conectar a PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	log.Println("Backend inicializado: conexión segura a Vault y PostgreSQL establecida.")

	// TODO: Inicializar servicios de dominio y adaptadores (HTTP, MQTT)
}
```

---

## 8. Resumen de arquitectura

Este diseño establece:

* Separación clara entre dominio e infraestructura.
* Gestión centralizada de secretos mediante Vault.
* Conexiones eficientes a PostgreSQL mediante pooling.
* Punto de entrada desacoplado preparado para extensiones (HTTP, MQTT).

El sistema queda preparado para evolucionar hacia un backend concurrente escalable siguiendo Clean Architecture.
