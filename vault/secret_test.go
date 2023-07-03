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
	dbVaultSecret := NewVaultSecret("database", "", KVPath)
	if dbVaultSecret.Engine != database || dbVaultSecret.Path != KVPath || dbVaultSecret.Mount != "database" || dbVaultSecret.Metadata != (Metadata{}) {
		test.Error("the database Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", dbVaultSecret.Engine, database)
		test.Errorf("expected path: %s, actual: %s", dbVaultSecret.Path, KVPath)
		test.Errorf("expected mount: %s, actual: %s", dbVaultSecret.Mount, "database")
		test.Errorf("expected empty metadata, actual: %v", dbVaultSecret.Metadata)
	}

	awsVaultSecret := NewVaultSecret("aws", "gcp", KVPath)
	if awsVaultSecret.Engine != aws || awsVaultSecret.Path != KVPath || awsVaultSecret.Mount != "gcp" || awsVaultSecret.Metadata != (Metadata{}) {
		test.Error("the AWS Vault secret constructor returned unexpected values")
		test.Errorf("expected engine: %s, actual: %s", awsVaultSecret.Engine, aws)
		test.Errorf("expected path: %s, actual: %s", awsVaultSecret.Path, KVPath)
		test.Errorf("expected mount: gcp, actual: %s", awsVaultSecret.Mount)
		test.Errorf("expected empty metadata, actual: %v", awsVaultSecret.Metadata)
	}
}

// test secret Read operation

// test secret generate credential

// test secret key value secret
func TestRetrieveKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.AuthClient()

	kv1VaultSecret := NewVaultSecret("kv1", "", KVPath)
	kv1Value, version, secretMetadata, err := kv1VaultSecret.retrieveKVSecret(basicVaultClient)

	if err != nil {
		test.Error("kv1 secret retrieval failed")
		test.Error(err)
	}
	if secretMetadata == (Metadata{}) {
		test.Error("the kv2 secret retrieval returned empty metadata")
	}
	if version != "0" {
		test.Errorf("the kv1 secret retrieval returned non-zero version: %s", version)
	}
	if kv1Value[KVKey] != KVValue {
		test.Error("the retrieved kv1 secret value was incorrect")
		test.Errorf("secret map value: %v", kv1Value)
	}

	kv2VaultSecret := NewVaultSecret("kv2", KV2Mount, KVPath)
	kv2Value, version, secretMetadata, err := kv2VaultSecret.retrieveKVSecret(basicVaultClient)

	if err != nil {
		test.Error("kv2 secret retrieval failed")
		test.Error(err)
	}
	if secretMetadata == (Metadata{}) {
		test.Error("the kv2 secret retrieval returned empty metadata")
	}
	if version == "0" {
		test.Errorf("the kv2 secret retrieval returned an invalid version: %s", version)
	}
	if kv2Value[KVKey] != KVValue {
		test.Error("the retrieved kv2 secret value was incorrect")
		test.Errorf("secret map value: %v", kv2Value)
	}
}

// test populate secret
func TestPopulateKVSecret(test *testing.T) {
	basicVaultClient := basicVaultConfig.AuthClient()

	kv1VaultSecret := NewVaultSecret("kv1", "", KVPath)
	version, secretMetadata, err := kv1VaultSecret.PopulateKVSecret(
		basicVaultClient,
		map[string]interface{}{KVKey: KVValue},
		false,
	)
	if err != nil {
		test.Error("the kv1 secret was not successfully put")
		test.Error(err)
	}
	if secretMetadata == (Metadata{}) {
		test.Error("the kv2 secret retrieval returned empty metadata")
	}
	if version != "0" {
		test.Errorf("the kv1 secret put returned non-zero version: %s", version)
	}

	kv2VaultSecret := NewVaultSecret("kv2", "", KVPath)
	version, secretMetadata, err = kv2VaultSecret.PopulateKVSecret(
		basicVaultClient,
		map[string]interface{}{KVKey: KVValue},
		false,
	)
	if err != nil {
		test.Error("the kv2 secret was not successfully put")
		test.Error(err)
	}
	if secretMetadata == (Metadata{}) {
		test.Error("the kv2 secret put returned empty metadata")
	}
	if version == "0" {
		test.Errorf("the kv2 secret put returned an invalid version: %s", version)
	}
	version, secretMetadata, err = kv2VaultSecret.PopulateKVSecret(
		basicVaultClient,
		map[string]interface{}{"other_password": "ultrasecret"},
		true,
	)
	if err != nil {
		test.Error("the kv2 secret was not successfully patched")
		test.Error(err)
	}
	if secretMetadata == (Metadata{}) {
		test.Error("the kv2 secret patch returned empty metadata")
	}
	if version == "0" {
		test.Errorf("the kv2 secret patch returned an invalid version: %s", version)
	}
}
