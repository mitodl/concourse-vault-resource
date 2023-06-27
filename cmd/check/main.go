package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/cmd"
	"github.com/mitodl/concourse-vault-resource/concourse"
)

// NOW doublecheck
func main() {
	// initialize checkRequest and secretSource
	checkRequest := concourse.NewCheckRequest(os.Stdin)
	secretSource := checkRequest.Source.Secret
	var err error

	// initialize vault client from concourse source
	vaultClient := helper.VaultClientFromSource(checkRequest.Source)

	// initialize vault secret from concourse params and invoke constructor
	secret := helper.VaultSecretFromParams(secretSource.Mount, secretSource.Engine, secretSource.Path)
	secret.New()

	// retrieve version for secret
	secretVersion := concourse.Version{}

	_, secretVersion[secretSource.Mount+"-"+secretSource.Path], _, err = secret.SecretValue(vaultClient)
	if err != nil {
		log.Fatalf("version could not be retrieved for %s engine, %s mount, and path %s secret", secretSource.Engine, secretSource.Mount, secretSource.Path)
	}

	// input secret version to constructed response NOW actually desire set of versions between requested version and retrieved version
	checkResponse := concourse.NewCheckResponse([]concourse.Version{secretVersion})

	// format checkResponse into json
	if err := json.NewEncoder(os.Stdout).Encode(&checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
