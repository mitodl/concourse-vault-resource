package helper

import (
  "log"

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

// instantiates vault secret from concourse params
func VaultSecretFromParams(mount string, engineString string) *vault.VaultSecret {
  // validate engine parameter
  engine := vault.SecretEngine(engineString)
  if len(engine) == 0 {
    log.Fatalf("an invalid secrets engine was specified: %s", engineString)
  }
  // initialize vault secret
  return &vault.VaultSecret{
    Mount:  mount,
    Engine: engine,
  }
}
