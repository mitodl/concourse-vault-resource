package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/cmd"
	"github.com/mitodl/concourse-vault-resource/concourse"
)

// GET and primary
func main() {
	// initialize request from concourse pipeline and response storing secret values
	inRequest := concourse.NewInRequest(os.Stdin)
	inResponse := concourse.NewInResponse(inRequest.Version)
	// initialize vault client from concourse source
	vaultClient := helper.VaultClientFromSource(inRequest.Source)

	// declare err specifically to track any SecretValue failure and trigger only after all secret operations
	var err error

	// perform secrets operations
	for mount, secretParams := range inRequest.Params {
		// initialize vault secret from concourse params
		secret := helper.VaultSecretFromParams(mount, secretParams.Engine)

		// iterate through secret params' paths and assign each to each vault secret path
		for _, secret.Path = range secretParams.Paths {
			// invoke secret constructor
			secret.New()
			// return and assign the secret values for the given path
			var values interface{}
			values, err = secret.SecretValue(vaultClient)
			secretValue := concourse.MetadataSecretValue{
				Name:  mount + "-" + secret.Path,
				Value: values,
			}
			// append to the response struct metadata values as key "<mount>-<path>" and value as secret keys and values
			inResponse.Metadata = append(inResponse.Metadata, secretValue)
		}
	}

	// fatally exit if any secret Read operation failed
	if err != nil {
		log.Fatal("one or more attempted secret Read operations failed")
	}

	// write marshalled metadata to file in at /opt/resource/vault.json
	helper.MetadataToJsonFile(inResponse.Metadata)

	// TODO: temp dummy metadata to bypass concourse type check
	dummyInResponse := concourse.NewInResponse(inRequest.Version)
	// marshal, encode, and pass inResponse json as output to concourse
	if err = json.NewEncoder(os.Stdout).Encode(dummyInResponse); err != nil {
		log.Print("unable to marshal in response struct to JSON")
		log.Fatal(err)
	}

	// TODO investigate if/how metadata populates concourse env vars ELSE write to <mount>.json for `load_var` later in pipeline
}
