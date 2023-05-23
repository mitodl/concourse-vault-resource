package helper

import (
	"os"
	"testing"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// minimum coverage testing for helper functions
func TestVaultClientFromSource(test *testing.T) {
	source := concourse.Source{Token: "abcdefghijklmnopqrstuvwxyz09"}
	VaultClientFromSource(source)
}

func TestVaultSecretFromParams(test *testing.T) {
	vaultSecret := VaultSecretFromParams("secret", "kv2")
	if vaultSecret.Mount != "secret" || vaultSecret.Engine != "kv2" {
		test.Error("the VaultSecret did contain the expected member fields")
		test.Errorf("expected mount: secret, actual: %s", vaultSecret.Mount)
		test.Errorf("expected engine: kv2, actual: %s", vaultSecret.Engine)
	}
}

func TestSecretsToJsonFile(test *testing.T) {
	secretValues := concourse.SecretValues{}
	SecretsToJsonFile(".", secretValues)
	defer os.Remove("./vault.json")
}
