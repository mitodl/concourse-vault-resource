package vault

import (
  "testing"

  "github.com/mitodl/concourse-vault-resource/concourse"
)

// test secret constructor
func TestNewVaultSecret(test *testing.T) {
  basicConcourseParams := &concourse.Params{
    SecretEngine: database,
    SecretPath:   "foo/bar",
    MountPath:    "database",
  }
  newVaultSecret := NewSecret(basicConcourseParams)

  if newVaultSecret.engine != basicConcourseParams.SecretEngine || newVaultSecret.path != basicConcourseParams.SecretPath || newVaultSecret.mount != basicConcourseParams.MountPath {
    test.Error("the Vault secret constructor returned unexpected values")
  }
}

// test secret retrieve secret

// test secret key value secret
func TestRetrieveKVSecret(test *testing.T) {
  kv1ConcourseParams := &concourse.Params{
    SecretEngine: keyvalue1,
    SecretPath:   "foo/bar",
    MountPath:    "kv",
  }
  kv1VaultSecret := NewSecret(kv1ConcourseParams)

  if kv1VaultSecret["password"] != "supersecret" {
    test.Error("the retrieved kv1 secret value was incorrect")
    test.Errorf("secret map value: %v", kv1VaultSecret)
  }

  kv2ConcourseParams := &concourse.Params{
    SecretEngine: keyvalue2,
    SecretPath:   "foo/bar",
    MountPath:    "secret",
  }
  kv2VaultSecret := NewSecret(kv2ConcourseParams)

  if kv2VaultSecret["password"] != "supersecret" {
    test.Error("the retrieved kv1 secret value was incorrect")
    test.Errorf("secret map value: %v", kv2VaultSecret)
  }
}
