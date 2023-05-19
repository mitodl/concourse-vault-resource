package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
	"github.com/mitodl/concourse-vault-resource/vault"
)

// GET and primary
func main() {
	// initialize request from concourse pipeline
	inRequest := concourse.NewInRequest(os.Stdin)

	// initialize response storing secret values
	inResponse := concourse.NewInResponse(inRequest.Version)

	// initialize vault config and client
	vaultConfig := &vault.VaultConfig{
		Engine:       vault.AuthEngine(inRequest.Source.AuthEngine),
		Address:      inRequest.Source.Address,
		AWSMountPath: inRequest.Source.AWSMountPath,
		AWSRole:      inRequest.Source.AWSVaultRole,
		Token:        inRequest.Source.Token,
		Insecure:     inRequest.Source.Insecure,
	}
	vaultConfig.New()
	vaultClient := vaultConfig.AuthClient()
	// declare err specifically to track any SecretValue failure and trigger only after all secret operations
	var err error

	// perform secrets operations
	for mount, secretParams := range inRequest.Params {
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
		// iterate through secret params' paths and assign each to each vault secret path
		for _, secret.Path = range secretParams.Paths {
			// invoke secret constructor
			secret.New()
			// return and assign the secret values for the given path
			secretValue := concourse.SecretValue{}
			secretValue[mount+"-"+secret.Path], err = secret.SecretValue(vaultClient)
			// append to the response struct metadata values as key "<mount>-<path>" and value as secret keys and values
			inResponse.Metadata = append(inResponse.Metadata, secretValue)
		}
	}

	// fatally exit if any secret Read operation failed
	if err != nil {
		log.Fatal("one or more attempted secret Read operations failed")
	}

	// format inResponse into json TODO: verify how this is behaving in concourse and how it can be captured for later use
	if err = json.NewEncoder(os.Stdout).Encode(inResponse); err != nil {
		log.Fatal("unable to unmarshal in response struct to JSON")
	}
}
