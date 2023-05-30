package helper

import (
	"os"
	"strconv"
	"testing"

	vault "github.com/hashicorp/vault/api"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// minimum coverage testing for helper functions
func TestVaultClientFromSource(test *testing.T) {
	source := concourse.Source{Token: "abcdefghijklmnopqrstuvwxyz09"}
	VaultClientFromSource(source)
}

func TestVaultSecretFromParams(test *testing.T) {
	vaultSecret := VaultSecretFromParams("secret", "kv2", "bar/baz")
	if vaultSecret.Mount != "secret" || vaultSecret.Engine != "kv2" || vaultSecret.Path != "bar/baz" {
		test.Error("the VaultSecret did contain the expected member fields")
		test.Errorf("expected mount: secret, actual: %s", vaultSecret.Mount)
		test.Errorf("expected engine: kv2, actual: %s", vaultSecret.Engine)
		test.Errorf("expected path: bar/baz, actual: %s", vaultSecret.Path)
	}
}

func TestSecretsToJsonFile(test *testing.T) {
	secretValues := concourse.SecretValues{}
	SecretsToJsonFile(".", secretValues)
	defer os.Remove("./vault.json")
}

func TestRawSecretToMetadata(test *testing.T) {
	rawSecret := &vault.Secret{
		LeaseID:       "abcdefg12345",
		LeaseDuration: 65535,
		Renewable:     false,
	}

	metadata := RawSecretToMetadata("secret-foo/bar", rawSecret)
	if len(metadata) != 3 {
		test.Error("metadata did not contain the expected number (three) entries per raw secret")
	}
	if metadata[0].Name != "secret-foo/bar-LeaseID" || metadata[0].Value != rawSecret.LeaseID {
		test.Error("first metadata entry is inaccurate")
		test.Errorf("expected name: secret-foo/bar-LeaseID, actual: %s", metadata[0].Name)
		test.Errorf("expected value: %s, actual: %s", rawSecret.LeaseID, metadata[0].Value)
	}
	if metadata[1].Name != "secret-foo/bar-LeaseDuration" || metadata[1].Value != strconv.Itoa(rawSecret.LeaseDuration) {
		test.Error("first metadata entry is inaccurate")
		test.Errorf("expected name: secret-foo/bar-LeaseDuration, actual: %s", metadata[1].Name)
		test.Errorf("expected value: %d, actual: %s", rawSecret.LeaseDuration, metadata[1].Value)
	}
	if metadata[2].Name != "secret-foo/bar-Renewable" || metadata[2].Value != strconv.FormatBool(rawSecret.Renewable) {
		test.Error("first metadata entry is inaccurate")
		test.Errorf("expected name: secret-foo/bar-Renewable, actual: %s", metadata[2].Name)
		test.Errorf("expected value: %t, actual: %s", rawSecret.Renewable, metadata[2].Value)
	}
}
