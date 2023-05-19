package vault

import (
	"testing"
)

// globals for vault package testing
const (
	KVPath   = "foo/bar"
	KVKey    = "password"
	KVValue  = "supersecret"
	KV1Mount = "kv"
	KV2Mount = "secret"
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

// test secret Read operation

// test secret generate credential

// test secret key value secret
func TestRetrieveKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.AuthClient()

	kv1VaultSecret := &VaultSecret{
		Engine: keyvalue1,
		Path:   KVPath,
	}
	kv1VaultSecret.New()
	kv1Value, err := kv1VaultSecret.retrieveKVSecret(basicVaultClient)

	if err != nil {
		test.Error("kv1 secret retrieval failed")
		test.Error(err)
	}
	if kv1Value[KVKey] != KVValue {
		test.Error("the retrieved kv1 secret value was incorrect")
		test.Errorf("secret map value: %v", kv1Value)
	}

	kv2VaultSecret := &VaultSecret{
		Engine: keyvalue2,
		Path:   KVPath,
		Mount:  KV2Mount,
	}
	kv2VaultSecret.New()
	kv2Value, err := kv2VaultSecret.retrieveKVSecret(basicVaultClient)

	if err != nil {
		test.Error("kv2 secret retrieval failed")
		test.Error(err)
	}
	if kv2Value[KVKey] != KVValue {
		test.Error("the retrieved kv2 secret value was incorrect")
		test.Errorf("secret map value: %v", kv2Value)
	}
}

// test populate secret
func TestPopulateKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.AuthClient()

	kv1VaultSecret := &VaultSecret{
		Engine: keyvalue1,
		Path:   KVPath,
	}
	kv1VaultSecret.New()
	err := kv1VaultSecret.PopulateKVSecret(basicVaultClient, map[string]interface{}{KVKey: KVValue})
	if err != nil {
		test.Error("the kv1 secret was not successfully put")
		test.Error(err)
	}

	kv2VaultSecret := &VaultSecret{
		Engine: keyvalue2,
		Path:   KVPath,
	}
	kv2VaultSecret.New()
	err = kv2VaultSecret.PopulateKVSecret(basicVaultClient, map[string]interface{}{KVKey: KVValue})
	if err != nil {
		test.Error("the kv2 secret was not successfully put")
		test.Error(err)
	}
}
