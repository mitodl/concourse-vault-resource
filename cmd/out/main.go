package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
	"github.com/mitodl/concourse-vault-resource/vault"
)

// PUT/POST
func main() {
	// initialize request from concourse pipeline
	outRequest := concourse.NewOutRequest(os.Stdin)

	// initialize response to satisfy concourse requirement
	outResponse := concourse.NewOutResponse(concourse.Version{Version: "1"})

	// initialize vault config and client
	vaultConfig := &vault.VaultConfig{
		Engine:       vault.AuthEngine(outRequest.Source.AuthEngine),
		Address:      outRequest.Source.Address,
		AWSMountPath: outRequest.Source.AWSMountPath,
		AWSRole:      outRequest.Source.AWSVaultRole,
		Token:        outRequest.Source.Token,
		Insecure:     outRequest.Source.Insecure,
	}
	vaultConfig.New()
	vaultClient := vaultConfig.AuthClient()
	// declare err specifically to track any SecretValue failure and trigger only after all secret operations
	var err error

	// perform secrets operations
	for mount, secretParams := range outRequest.Params {
		// validate engine parameter
		engineString := secretParams.Engine
		engine := vault.SecretEngine(engineString)
		if len(engine) == 0 {
			log.Fatalf("an invalid secrets engine was specified: %s", engineString)
		}
		// initialize vault secret
		secret := &vault.VaultSecret{
			Mount:  mount,
			Engine: engine,
		}
		// iterate through secrets and assign each path to each vault secret path, and write each secret value to the path
		var secretValue concourse.SecretValue
		for secret.Path, secretValue = range secretParams.Secrets {
			// invoke secret constructor
			secret.New()
			// write the secret value to the path for the specified mount and engine
			err = secret.PopulateKVSecret(vaultClient, secretValue)
		}
	}

	// fatally exit if any secret Read operation failed
	if err != nil {
		log.Fatal("one or more attempted secret Create/Update operations failed")
	}

	// format outResponse into json TODO: verify how this is behaving in concourse and how it can be captured for later use
	if err = json.NewEncoder(os.Stdout).Encode(outResponse); err != nil {
		log.Fatal("unable to unmarshal out response struct to JSON")
	}
}
