// TODO split helpers into another file as this is becoming unwieldy
package vault

import (
	"context"
	"log"
	"strconv"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// secret engine with pseudo-enum TODO private
type secretEngine string

const (
	database  secretEngine = "database"
	aws       secretEngine = "aws"
	keyvalue1 secretEngine = "kv1"
	keyvalue2 secretEngine = "kv2"
)

// secret metadata
type Metadata struct {
	LeaseID       string
	LeaseDuration string
	Renewable     string
}

// secret defines a composite Vault secret configuration
type vaultSecret struct {
	engine   secretEngine
	metadata Metadata
	mount    string
	path     string
	dynamic  bool
}

// secret constructor
func NewVaultSecret(engineString string, mount string, path string) *vaultSecret {
	// validate mandatory fields specified
	if len(engineString) == 0 || len(path) == 0 {
		log.Fatal("the secret engine and path parameters are mandatory")
	}

	// validate engine parameter
	engine := secretEngine(engineString)
	if len(engine) == 0 {
		log.Fatalf("an invalid secrets engine was specified: %s", engineString)
	}

	// initialize vault secret
	vaultSecret := &vaultSecret{
		engine: engine,
		path:   path,
		mount:  mount,
	}

	// determine if secret is dynamic TODO use this
	switch engine {
	case database, aws:
		vaultSecret.dynamic = true
	case keyvalue1, keyvalue2:
		vaultSecret.dynamic = false
	default:
		log.Fatalf("an invalid secret engine %s was selected", engine)
	}

	// determine default mount path if not specified
	// note current schema renders this pointless, but it would ensure safety to retain
	if len(mount) == 0 {
		switch engine {
		case database:
			vaultSecret.mount = "database"
		case aws:
			vaultSecret.mount = "aws"
		case keyvalue1:
			vaultSecret.mount = "kv"
		case keyvalue2:
			vaultSecret.mount = "secret"
		default:
			log.Fatalf("an invalid secret engine %s was selected", engine)
		}
	}

	return vaultSecret
}

// secret readers
func (secret *vaultSecret) Engine() secretEngine {
	return secret.engine
}

func (secret *vaultSecret) Metadata() Metadata {
	return secret.metadata
}

func (secret *vaultSecret) Mount() string {
	return secret.mount
}

func (secret *vaultSecret) Path() string {
	return secret.path
}

func (secret *vaultSecret) Dynamic() bool {
	return secret.dynamic
}

// return secret value, version, metadata, and possible error
func (secret *vaultSecret) SecretValue(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	switch secret.engine {
	case database, aws:
		return secret.generateCredentials(client)
	case keyvalue1, keyvalue2:
		return secret.retrieveKVSecret(client)
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.engine)
		return map[string]interface{}{}, "0", Metadata{}, nil // unreachable code, but compile error otherwise
	}
}

// generate credentials
func (secret *vaultSecret) generateCredentials(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	// initialize api endpoint for cred generation
	endpoint := secret.mount + "/creds/" + secret.path
	// GET the secret from the API endpoint
	response, err := client.Logical().Read(endpoint)
	if err != nil {
		log.Printf("failed to generate credentials for %s with %s secrets engine", secret.path, secret.engine)
		log.Print(err)
		return map[string]interface{}{}, "0", Metadata{}, err
	}
	// calculate the expiration time for version
	expirationTime := time.Now().Local().Add(time.Second * time.Duration(response.LeaseDuration))

	// return secret value implicitly coerced to map[string]interface{}, expiration time as version, and metadata
	return response.Data, expirationTime.String(), rawSecretToMetadata(response), nil
}

// retrieve key-value pair secrets
func (secret *vaultSecret) retrieveKVSecret(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	// declare error for return to cmd, and kvSecret for metadata.version and raw secret assignments and returns
	var err error
	var kvSecret *vault.KVSecret

	switch secret.engine {
	case keyvalue1:
		// read kv secret
		kvSecret, err = client.KVv1(secret.mount).Get(
			context.Background(),
			secret.path,
		)
		// instantiate dummy metadata if secret successfully retrieved
		if err == nil && kvSecret != nil {
			kvSecret.VersionMetadata = &vault.KVVersionMetadata{Version: 0}
		}
	case keyvalue2:
		// read kv2 secret
		kvSecret, err = client.KVv2(secret.mount).Get(
			context.Background(),
			secret.path,
		)
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.engine)
	}

	// verify secret read
	if err != nil || kvSecret == nil {
		log.Printf("failed to read secret at mount %s and path %s from %s secrets engine", secret.mount, secret.path, secret.engine)
		log.Print(err)
		// return empty values since error triggers at end of execution
		return map[string]interface{}{}, "0", Metadata{}, err
	}

	// return secret value and implicitly coerce type to map[string]interface{}
	return kvSecret.Data, strconv.Itoa(kvSecret.VersionMetadata.Version), rawSecretToMetadata(kvSecret.Raw), nil
}

// populate key-value pair secrets and return version, metadata, and error
func (secret *vaultSecret) PopulateKVSecret(client *vault.Client, secretValue map[string]interface{}, patch bool) (string, Metadata, error) {
	// declare error for return to cmd, and kvSecret for metadata.version and raw secret assignments and returns (with dummies for kv1)
	var err error
	kvSecret := &vault.KVSecret{
		VersionMetadata: &vault.KVVersionMetadata{Version: 0},
		Raw:             &vault.Secret{},
	}

	switch secret.engine {
	case keyvalue1:
		// put kv1 secret
		err = client.KVv1(secret.mount).Put(
			context.Background(),
			secret.path,
			secretValue,
		)
	case keyvalue2:
		if patch {
			// patch kv2 secret
			kvSecret, err = client.KVv2(secret.mount).Patch(
				context.Background(),
				secret.path,
				secretValue,
			)
		} else {
			// put kv2 secret
			kvSecret, err = client.KVv2(secret.mount).Put(
				context.Background(),
				secret.path,
				secretValue,
			)
		}
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.engine)
	}

	// verify secret put
	if err != nil {
		log.Printf("failed to update secret %s into %s secrets Engine", secret.path, secret.engine)
		log.Print(err)
		return "0", Metadata{}, err
	}

	// return no error
	return strconv.Itoa(kvSecret.VersionMetadata.Version), rawSecretToMetadata(kvSecret.Raw), nil
}

// convert *vault.Secret raw secret to secret metadata
func rawSecretToMetadata(rawSecret *vault.Secret) Metadata {
	// returne metadata with fields populated from raw secret
	return Metadata{
		LeaseID:       rawSecret.LeaseID,
		LeaseDuration: strconv.Itoa(rawSecret.LeaseDuration),
		Renewable:     strconv.FormatBool(rawSecret.Renewable),
	}
}
