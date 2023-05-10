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
	// read, decode, unmarshal, and store the json stdin in the inRequest pointer
	var inRequest concourse.InRequest

	if err := json.NewDecoder(os.Stdin).Decode(&inRequest); err != nil {
		log.Print("error decoding stdin from JSON")
		log.Fatal(err)
	}

	// initialize response storing secret values
	inResponse := concourse.InResponse{Version: concourse.Version{Version: 1}}

	// initialize vault config and client
	vaultConfig := &vault.VaultConfig{
		Engine:       vault.AuthEngine(inRequest.Source.Engine),
		Address:      inRequest.Source.Address,
		AWSMountPath: inRequest.Source.AWSMountPath,
		Token:        inRequest.Source.Token,
		Insecure:     inRequest.Source.Insecure,
	}
	vaultConfig.New()
	vaultClient := vaultConfig.AuthClient()

	// perform secrets operations
	for mount, secretParams := range inRequest.Source.Secrets {
		// validate parameters
		engineAny, ok := secretParams["engine"]
		if !ok {
			log.Fatalf("The secret engine was not specified for mount: %s", mount)
		}
		engineString, ok := engineAny.(string)
		if !ok {
			log.Fatalf("The secret engine must be a string: %v", engineAny)
		}
		engine := vault.SecretEngine(engineString)
		if len(engine) == 0 {
			log.Fatalf("An invalid secrets engine was specified: %s", engineString)
		}
		pathsAny, ok := secretParams["paths"]
		if !ok {
			log.Fatalf("The paths were not specified for mount: %s", mount)
		}
		paths, ok := pathsAny.([]string)
		if !ok {
			log.Fatalf("The secret paths must be a list of strings: %v", pathsAny)
		}
		// initialize vault secret
		secret := &vault.VaultSecret{
			Mount:  mount,
			Engine: engine,
		}
		// iterate through secret paths
		for _, path := range paths {
			// assign path to secret member and then invoke constructor
			secret.Path = path
			secret.New()
			// populate the secret value
			secret.PopulateSecret(vaultClient)
			// assign to the response struct as key "<mount>-<path>" and value as secret keys and values
			inResponse.Metadata.Values[mount+"-"+path] = secret.Value
		}
	}

	// format inResponse into json
	if err := json.NewEncoder(os.Stdout).Encode(inResponse); err != nil {
		log.Fatal("unable to format secret values into JSON")
	}
}
