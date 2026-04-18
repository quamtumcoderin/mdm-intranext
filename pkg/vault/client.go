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
		return nil, fmt.Errorf("error creating vault client: %s", err)
	}

	client.SetToken(cfg.Token)

	return &SecretManager{client: client}, nil
}

func (s *SecretManager) GetDatabaseCredentials(ctx context.Context) (string, string, string, error) {
	secret, err := s.client.KVv2("secret").Get(ctx, "mdm/database")
	if err != nil {
		return "", "", "", fmt.Errorf("error retrieving database credentials: %s", err)
	}

	user := secret.Data["username"].(string)
	pass := secret.Data["password"].(string)
	dbname := secret.Data["dbname"].(string)

	return user, pass, dbname, nil
}
