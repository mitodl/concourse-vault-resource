package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
	"github.com/mitodl/concourse-vault-resource/cmd"
)

// PUT/POST
func main() {
	// initialize request from concourse pipeline and response to satisfy concourse requirement
	outRequest := concourse.NewOutRequest(os.Stdin)
	outResponse := concourse.NewOutResponse(concourse.Version{Version: "1"})
	// initialize vault client from concourse source
	vaultClient := helper.VaultClientFromSource(outRequest.Source)

	// declare err specifically to track any SecretValue failure and trigger only after all secret operations
	var err error

	// perform secrets operations
	for mount, secretParams := range outRequest.Params {
		// initialize vault secret from concourse params
		secret := helper.VaultSecretFromParams(mount, secretParams.Engine)

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

	// format outResponse into json
	if err = json.NewEncoder(os.Stdout).Encode(outResponse); err != nil {
		log.Fatal("unable to marshal out response struct to JSON")
	}
}
