package main

import (
	"encoding/json"
	"log"
	"os"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/mitodl/concourse-vault-resource/cmd"
	"github.com/mitodl/concourse-vault-resource/concourse"
)

// GET and primary
func main() {
	// initialize request from concourse pipeline and response storing secret values
	inRequest := concourse.NewInRequest(os.Stdin)
	inResponse := concourse.NewResponse(inRequest.Version)
	// initialize vault client from concourse source
	vaultClient := helper.VaultClientFromSource(inRequest.Source)

	// declare err specifically to track any SecretValue failure and trigger only after all secret operations
	var err error
	// initialize secretValues to store aggregated retrieved secrets and secretSource for efficiency
	secretValues := concourse.SecretValues{}
	secretSource := inRequest.Source.Secret

	// read secrets from params
	if secretSource == (concourse.SecretSource{}) {
		// perform secrets operations
		for mount, secretParams := range inRequest.Params {
			// initialize vault secret from concourse params
			secret := helper.VaultSecretFromParams(mount, secretParams.Engine, "")

			// iterate through secret params' paths and assign each to each vault secret path
			for _, secret.Path = range secretParams.Paths {
				// invoke secret constructor
				secret.New()
				// declare identifier and rawSecret
				identifier := mount + "-" + secret.Path
				var rawSecret *vaultapi.Secret
				// return and assign the secret values for the given path
				secretValues[identifier], inResponse.Version[identifier], rawSecret, err = secret.SecretValue(vaultClient)
				// convert rawSecret to concourse metadata and append to metadata
				inResponse.Metadata = append(inResponse.Metadata, helper.RawSecretToMetadata(identifier, rawSecret)...)
			}
		}
	} else { // read secret from source TODO cleanup and dry with above
		// initialize vault secret from concourse params and invoke constructor
		secret := helper.VaultSecretFromParams(secretSource.Mount, secretSource.Engine, secretSource.Path)
		secret.New()
		// declare identifier and rawSecret
		identifier := secretSource.Mount + "-" + secretSource.Path
		var rawSecret *vaultapi.Secret
		// return and assign the secret values for the given path
		secretValues[identifier], inResponse.Version[identifier], rawSecret, err = secret.SecretValue(vaultClient)
		// convert rawSecret to concourse metadata and append to metadata
		inResponse.Metadata = append(inResponse.Metadata, helper.RawSecretToMetadata(identifier, rawSecret)...)
	}

	// fatally exit if any secret Read operation failed
	if err != nil {
		log.Fatal("one or more attempted secret Read operations failed")
	}

	// write marshalled metadata to file in at /opt/resource/vault.json
	helper.SecretsToJsonFile(os.Args[1], secretValues)

	// marshal, encode, and pass inResponse json as output to concourse
	if err = json.NewEncoder(os.Stdout).Encode(inResponse); err != nil {
		log.Print("unable to marshal in response struct to JSON")
		log.Fatal(err)
	}
}
