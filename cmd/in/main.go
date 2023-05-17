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
		Token:        inRequest.Source.Token,
		Insecure:     inRequest.Source.Insecure,
	}
	vaultConfig.New()
	vaultClient := vaultConfig.AuthClient()

	// perform secrets operations
	for mount, secretParams := range inRequest.Params {
		// validate parameters with generic and/or custom types
		engineAny, ok := secretParams["engine"]
		if !ok {
			log.Fatalf("the secret engine was not specified for mount: %s", mount)
		}
		engineString, ok := engineAny.(string)
		if !ok {
			log.Fatalf("the secret engine must be a string: %v, but was instead a type: %T", engineAny, engineAny)
		}
		engine := vault.SecretEngine(engineString)
		if len(engine) == 0 {
			log.Fatalf("an invalid secrets engine was specified: %s", engineString)
		}
		pathsAny, ok := secretParams["paths"]
		if !ok {
			log.Fatalf("the paths were not specified for mount: %s", mount)
		}
		pathsSlice, ok := pathsAny.([]any)
		if !ok {
			log.Fatalf("the secret paths must be a list of strings: %v, but was instead a type: %T", pathsAny, pathsAny)
		}
		// initialize vault secret
		secret := &vault.VaultSecret{
			Mount:  mount,
			Engine: engine,
		}
		// iterate through secret paths
		for _, pathAny := range pathsSlice {
			// validate path parameter can be converted to a string and assign to secret member
			secret.Path, ok = pathAny.(string)
			if !ok {
				log.Fatalf("the secret path must be a string: %v, but was instead a type: %T", pathAny, pathAny)
			}
			// invoke secret constructor
			secret.New()
			// return the secret value and assign to the response struct as key "<mount>-<path>" and value as secret keys and values
			inResponse.Metadata.Values[mount+"-"+secret.Path] = secret.SecretValue(vaultClient)
		}
	}

	// format inResponse into json TODO: verify how this is behaving in concourse and how it can be captured for later use
	if err := json.NewEncoder(os.Stdout).Encode(inResponse); err != nil {
		log.Fatal("unable to format secret values into JSON")
	}
}
