package vault

import (
	"context"
	"log"
	"strconv"
	"time"

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

// secret metadata
type Metadata struct {
	LeaseID       string
	LeaseDuration string
	Renewable     string
}

// secret defines a composite Vault secret configuration; TODO: split value into kv or cred value and then use in functions (maybe re-add to vaultSecret and then re-populate instead of return?); consider making these private and defining getters
type vaultSecret struct {
	Engine   SecretEngine
	Metadata Metadata
	Mount    string
	Path     string
}

// secret constructor
func NewVaultSecret(engineString string, mount string, path string) *vaultSecret {
	// validate mandatory fields specified
	if len(engineString) == 0 || len(path) == 0 {
		log.Fatal("the secret engine and path parameters are mandatory")
	}

	// validate engine parameter
	engine := SecretEngine(engineString)
	if len(engine) == 0 {
		log.Fatalf("an invalid secrets engine was specified: %s", engineString)
	}

	// initialize vault secret
	vaultSecret := &vaultSecret{
		Engine: engine,
		Path:   path,
		Mount:  mount,
	}

	// determine default mount path if not specified
	// note current schema renders this pointless, but it would ensure safety to retain
	if len(mount) == 0 {
		switch engine {
		case database:
			vaultSecret.Mount = "database"
		case aws:
			vaultSecret.Mount = "aws"
		case keyvalue1:
			vaultSecret.Mount = "kv"
		case keyvalue2:
			vaultSecret.Mount = "secret"
		default:
			log.Fatalf("an invalid secret engine %s was selected", engine)
		}
	}

	return vaultSecret
}

// return secret value, version, metadata, and possible error
func (secret *vaultSecret) SecretValue(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	switch secret.Engine {
	case database, aws:
		return secret.generateCredentials(client)
	case keyvalue1, keyvalue2:
		return secret.retrieveKVSecret(client)
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
		return map[string]interface{}{}, "0", Metadata{}, nil // unreachable code, but compile error otherwise
	}
}

// generate credentials
func (secret *vaultSecret) generateCredentials(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	// initialize api endpoint for cred generation
	endpoint := secret.Mount + "/creds/" + secret.Path
	// GET the secret from the API endpoint
	response, err := client.Logical().Read(endpoint)
	if err != nil {
		log.Printf("failed to generate credentials for %s with %s secrets engine", secret.Path, secret.Engine)
		log.Print(err)
		return map[string]interface{}{}, "0", Metadata{}, err
	}
	// calculate the expiration time for version
	expirationTime := time.Now().Local().Add(time.Second * time.Duration(response.LeaseDuration))

	// return secret value implicitly coerced to map[string]interface{}, expiration time as version, and metadata TODO expiration time to string conversion looks terrible
	return response.Data, expirationTime.String(), rawSecretToMetadata(response), nil
}

// retrieve key-value pair secrets
func (secret *vaultSecret) retrieveKVSecret(client *vault.Client) (map[string]interface{}, string, Metadata, error) {
	// declare error for return to cmd, and kvSecret for metadata.version and raw secret assignments and returns
	var err error
	var kvSecret *vault.KVSecret

	switch secret.Engine {
	case keyvalue1:
		// read kv secret
		kvSecret, err = client.KVv1(secret.Mount).Get(
			context.Background(),
			secret.Path,
		)
		// instantiate dummy metadata if secret successfully retrieved
		if err == nil && kvSecret != nil {
			kvSecret.VersionMetadata = &vault.KVVersionMetadata{Version: 0}
		}
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
	if err != nil || kvSecret == nil {
		log.Printf("failed to read secret at mount %s and path %s from %s secrets engine", secret.Mount, secret.Path, secret.Engine)
		log.Print(err)
		// return empty values since error triggers at end of execution
		return map[string]interface{}{}, "0", Metadata{}, err
	}

	// return secret value and implicitly coerce type to map[string]interface{}
	return kvSecret.Data, strconv.Itoa(kvSecret.VersionMetadata.Version), rawSecretToMetadata(kvSecret.Raw), nil
}

// populate key-value pair secrets and return version, metadata, and error
func (secret *vaultSecret) PopulateKVSecret(client *vault.Client, secretValue map[string]interface{}, patch bool) (string, Metadata, error) {
	// declare error for return to cmd, and kvSecret for metadata.version and raw secret assignments and returns TODO compare/contrast below init with declare in get
	var err error
	kvSecret := &vault.KVSecret{
		VersionMetadata: &vault.KVVersionMetadata{Version: 0},
		Raw:             &vault.Secret{},
	}

	switch secret.Engine {
	case keyvalue1:
		// put kv1 secret
		err = client.KVv1(secret.Mount).Put(
			context.Background(),
			secret.Path,
			secretValue,
		)
	case keyvalue2:
		if patch {
			// patch kv2 secret
			kvSecret, err = client.KVv2(secret.Mount).Patch(
				context.Background(),
				secret.Path,
				secretValue,
			)
		} else {
			// put kv2 secret
			kvSecret, err = client.KVv2(secret.Mount).Put(
				context.Background(),
				secret.Path,
				secretValue,
			)
		}
	default:
		log.Fatalf("an invalid secret engine %s was selected", secret.Engine)
	}

	// verify secret put
	if err != nil {
		log.Printf("failed to update secret %s into %s secrets Engine", secret.Path, secret.Engine)
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
