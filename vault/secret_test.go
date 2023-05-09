package vault

import (
	"testing"
)

// test secret constructor
func TestNewVaultSecret(test *testing.T) {
	dbVaultSecret := &VaultSecret{
		Engine: database,
		Path:   "foo/bar",
	}
	dbVaultSecret.New()

	if dbVaultSecret.Engine != database || dbVaultSecret.Path != "foo/bar" || dbVaultSecret.Mount != "database" {
		test.Error("the database Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", dbVaultSecret.Engine, database)
		test.Errorf("expected path: %s, actual: %s", dbVaultSecret.Path, "foo/bar")
		test.Errorf("expected mount: %s, actual: %s", dbVaultSecret.Mount, "database")
	}

	awsVaultSecret := &VaultSecret{
		Engine: aws,
		Path:   "foo/bar",
		Mount:  "gcp",
	}
	awsVaultSecret.New()

	if awsVaultSecret.Engine != aws || awsVaultSecret.Path != "foo/bar" || awsVaultSecret.Mount != "gcp" {
		test.Error("the AWS Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", awsVaultSecret.Engine, aws)
		test.Errorf("expected path: %s, actual: %s", awsVaultSecret.Path, "foo/bar")
		test.Errorf("expected mount: %s, actual: %s", awsVaultSecret.Mount, "gcp")
	}
}

// test secret retrieve secret

// test secret key value secret
func TestRetrieveKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.authClient()

	kv1VaultSecret := &VaultSecret{
		Engine: keyvalue1,
		Path:   "foo/bar",
	}
	kv1VaultSecret.New()
	kv1VaultSecret.retrieveKVSecret(basicVaultClient)

	if kv1VaultSecret.Value["password"] != "supersecret" {
		test.Error("the retrieved kv1 secret value was incorrect")
		test.Errorf("secret map value: %v", kv1VaultSecret.Value)
	}

	kv2VaultSecret := &VaultSecret{
		Engine: keyvalue2,
		Path:   "foo/bar",
		Mount:  "secret",
	}
	kv2VaultSecret.New()
	kv2VaultSecret.retrieveKVSecret(basicVaultClient)

	if kv2VaultSecret.Value["password"] != "supersecret" {
		test.Error("the retrieved kv2 secret value was incorrect")
		test.Errorf("secret map value: %v", kv2VaultSecret.Value)
	}
}
