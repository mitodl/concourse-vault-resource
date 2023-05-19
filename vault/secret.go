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

// secret defines a composite Vault secret configuration; TODO: convert value into kv or cred value and then use in functions (maybe re-add to VaultSecret and then re-populate instead of return?)
type VaultSecret struct {
	Engine SecretEngine
	Path   string
	Mount  string
}

// secret constructor; TODO does not need model type, so could be proper constructor
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
		return map[string]interface{}{}, nil // unreachable code, but compile error otherwise
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
		// read kv2 secret
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

// populate key-value pair secrets TODO: enable value merges with current value and propagate upwards
func (secret *VaultSecret) PopulateKVSecret(client *vault.Client, secretValue map[string]interface{}) error {
	// declare error for later reporting
	var err error

	switch secret.Engine {
	case keyvalue1:
		// put kv secret
		err = client.KVv1(secret.Mount).Put(
			context.Background(),
			secret.Path,
			secretValue,
		)
	case keyvalue2:
		// put kv2 secret TODO validate kvSecret.Data return == secretValue without screwing up err scope
		_, err = client.KVv2(secret.Mount).Put(
			context.Background(),
			secret.Path,
			secretValue,
		)
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
	}

	// verify secret put
	if err != nil {
		log.Printf("failed to put secret %s into %s secrets Engine", secret.Path, secret.Engine)
		log.Print(err)
		return err
	}

	// return no error
	return nil
}
