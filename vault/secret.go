package vault

import (
	"context"
	"log"

	vault "github.com/hashicorp/vault/api"
)

// secret engine with pseudo-enum
type SecretEngine string

const (
	database  SecretEngine = "DB"
	aws       SecretEngine = "AWS"
	keyvalue1 SecretEngine = "KV1"
	keyvalue2 SecretEngine = "KV2"
)

// secret defines a composite Vault secret configuration
type VaultSecret struct {
	Engine SecretEngine
	Path   string
	Mount  string
	Value  map[string]interface{}
}

func (secret *VaultSecret) New() {
	// determine default mount path if not specified
  if len(secret.Mount) == 0 {
		switch secret.Engine {
		case database:
			secret.Mount = "database"
		case aws:
			secret.Mount = "aws"
		case keyvalue1:
			secret.Mount = "kv"
		case keyvalue2:
			secret.Mount = "secret"
		default:
			log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
	  }
	}
}

// retrieve and return secrets
func (secret *VaultSecret) generateCredentials(client *vault.Client) {
	// initialize api endpoint for cred generation
	endpoint := secret.Mount + "/creds/" + secret.Path
	// GET the secret from the API endpoint
	response, err := client.Logical().Read(endpoint)
	if err != nil {
		log.Printf("failed to generate credentials for %s with %s secrets engine", secret.Path, secret.Engine)
		log.Fatal(err)
	}

  // assign secret value and implicitly coerce type to map[string]interface{}
	secret.Value = response.Data
}

func (secret *VaultSecret) retrieveKVSecret(client *vault.Client) {
	// declare func scope variable
	var kvSecret *vault.KVSecret
	var err error

	switch secret.Engine {
	case keyvalue1:
		// read kv secret
		kvSecret, err = client.KVv1(secret.Mount).Get(
			context.Background(),
			secret.Path,
		)
	case keyvalue2:
		// read kv secret
		kvSecret, err = client.KVv2(secret.Mount).Get(
			context.Background(),
			secret.Path,
		)
	}

	// verify secret read
	if err != nil {
		log.Printf("failed to read secret %s from %s secrets Engine", secret.Path, secret.Engine)
		log.Fatal(err)
	}

	// assign secret value and implicitly coerce type to map[string]interface{}
	secret.Value = kvSecret.Data
}
