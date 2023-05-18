package vault

import (
	"context"
	"log"

	vault "github.com/hashicorp/vault/api"
)

// secret engine with pseudo-enum
type SecretEngine string

const (
	database  SecretEngine = "database"
	aws       SecretEngine = "aws"
	keyvalue1 SecretEngine = "kv1"
	keyvalue2 SecretEngine = "kv2"
)

// secret defines a composite Vault secret configuration; TODO: convert value into kv or cred value and propagate to concourse models metadata.values?
type VaultSecret struct {
	Engine SecretEngine
	Path   string
	Mount  string
}

// secret constructor
func (secret *VaultSecret) New() {
	// validate mandatory fields specified
	if len(secret.Engine) == 0 || len(secret.Path) == 0 {
		log.Fatal("the secret engine and path parameters are mandatory")
	}

	// determine default mount path if not specified
	// note current schema renders this pointless, but it would ensure safety to retain
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

// populate secret type struct with value
func (secret *VaultSecret) SecretValue(client *vault.Client) (map[string]interface{}, error) {
	switch secret.Engine {
	case database, aws:
		return secret.generateCredentials(client)
	case keyvalue1, keyvalue2:
		return secret.retrieveKVSecret(client)
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
		return map[string]interface{}{}, nil // unreachable, but compile error otherwise
	}
}

// generate credentials
func (secret *VaultSecret) generateCredentials(client *vault.Client) (map[string]interface{}, error) {
	// initialize api endpoint for cred generation
	endpoint := secret.Mount + "/creds/" + secret.Path
	// GET the secret from the API endpoint
	response, err := client.Logical().Read(endpoint)
	if err != nil {
		log.Printf("failed to generate credentials for %s with %s secrets engine", secret.Path, secret.Engine)
		log.Print(err)
		return map[string]interface{}{}, err
	}

	// return secret value and implicitly coerce type to map[string]interface{}
	// TODO: return data key?
	return response.Data, nil
}

// retrieve key-value pair secrets
func (secret *VaultSecret) retrieveKVSecret(client *vault.Client) (map[string]interface{}, error) {
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
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
	}

	// verify secret read
	if err != nil {
		log.Printf("failed to read secret %s from %s secrets Engine", secret.Path, secret.Engine)
		log.Print(err)
		return map[string]interface{}{}, err
	}

	// return secret value and implicitly coerce type to map[string]interface{}
	return kvSecret.Data, nil
}
