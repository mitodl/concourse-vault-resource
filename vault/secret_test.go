package vault

import (
	"testing"
)

const (
	KVPath  = "foo/bar"
	KVKey   = "password"
	KVValue = "supersecret"
)

// test secret constructor
func TestNewVaultSecret(test *testing.T) {
	dbVaultSecret := &VaultSecret{
		Engine: database,
		Path:   KVPath,
	}
	dbVaultSecret.New()

	if dbVaultSecret.Engine != database || dbVaultSecret.Path != KVPath || dbVaultSecret.Mount != "database" {
		test.Error("the database Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", dbVaultSecret.Engine, database)
		test.Errorf("expected path: %s, actual: %s", dbVaultSecret.Path, KVPath)
		test.Errorf("expected mount: %s, actual: %s", dbVaultSecret.Mount, "database")
	}

	awsVaultSecret := &VaultSecret{
		Engine: aws,
		Path:   KVPath,
		Mount:  "gcp",
	}
	awsVaultSecret.New()

	if awsVaultSecret.Engine != aws || awsVaultSecret.Path != KVPath || awsVaultSecret.Mount != "gcp" {
		test.Error("the AWS Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", awsVaultSecret.Engine, aws)
		test.Errorf("expected path: %s, actual: %s", awsVaultSecret.Path, KVPath)
		test.Errorf("expected mount: %s, actual: %s", awsVaultSecret.Mount, "gcp")
	}
}

// test secret retrieve secret

// test secret key value secret
func TestRetrieveKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.AuthClient()

	kv1VaultSecret := &VaultSecret{
		Engine: keyvalue1,
		Path:   KVPath,
	}
	kv1VaultSecret.New()
	kv1Value := kv1VaultSecret.retrieveKVSecret(basicVaultClient)

	if kv1Value[KVKey] != KVValue {
		test.Error("the retrieved kv1 secret value was incorrect")
		test.Errorf("secret map value: %v", kv1Value)
	}

	kv2VaultSecret := &VaultSecret{
		Engine: keyvalue2,
		Path:   KVPath,
		Mount:  "secret",
	}
	kv2VaultSecret.New()
	kv2Value := kv2VaultSecret.retrieveKVSecret(basicVaultClient)

	if kv2Value[KVKey] != KVValue {
		test.Error("the retrieved kv2 secret value was incorrect")
		test.Errorf("secret map value: %v", kv2Value)
	}
}

// test populate secret
