package helper

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/mitodl/concourse-vault-resource/concourse"
	"github.com/mitodl/concourse-vault-resource/vault"
)

// instantiates vault client from concourse source
func VaultClientFromSource(source concourse.Source) *vaultapi.Client {
	// initialize vault config and client
	vaultConfig := &vault.VaultConfig{
		Engine:       vault.AuthEngine(source.AuthEngine),
		Address:      source.Address,
		AWSMountPath: source.AWSMountPath,
		AWSRole:      source.AWSVaultRole,
		Token:        source.Token,
		Insecure:     source.Insecure,
	}
	vaultConfig.New()
	return vaultConfig.AuthClient()
}

// instantiates vault secret from concourse params or source
func VaultSecretFromParams(mount string, engineString string, path string) *vault.VaultSecret {
	// validate engine parameter
	engine := vault.SecretEngine(engineString)
	if len(engine) == 0 {
		log.Fatalf("an invalid secrets engine was specified: %s", engineString)
	}
	// initialize vault secret and return
	return &vault.VaultSecret{
		Mount:  mount,
		Engine: engine,
		Path:   path,
	}
}

// writes inResponse.Metadata marshalled to json to file at /opt/resource/vault.json
func SecretsToJsonFile(filePath string, secretValues concourse.SecretValues) {
	// marshal secretValues into json data
	secretsData, err := json.Marshal(secretValues)
	if err != nil {
		log.Print("unable to marshal SecretValues struct to json data")
		log.Fatal(err)
	}
	// write secrets to file at /opt/resource/vault.json
	secretsFile := filePath + "/vault.json"
	if err = os.WriteFile(secretsFile, secretsData, 0o600); err != nil {
		log.Printf("error writing secrets to destination file at %s", secretsFile)
		log.Fatal(err)
	}
}

// converts Vault raw secret information to Concourse metadata
func RawSecretToMetadata(prefix string, rawSecret *vaultapi.Secret) []concourse.MetadataEntry {
	// initialize metadata entries for raw secret
	var metadataEntries []concourse.MetadataEntry

	// append lease id, lease duration, and renewable converted to string to the entries
	metadataEntries = append(metadataEntries, concourse.MetadataEntry{
		Name:  prefix + "-LeaseID",
		Value: rawSecret.LeaseID,
	})
	metadataEntries = append(metadataEntries, concourse.MetadataEntry{
		Name:  prefix + "-LeaseDuration",
		Value: strconv.Itoa(rawSecret.LeaseDuration),
	})
	metadataEntries = append(metadataEntries, concourse.MetadataEntry{
		Name:  prefix + "-Renewable",
		Value: strconv.FormatBool(rawSecret.Renewable),
	})

	return metadataEntries
}
