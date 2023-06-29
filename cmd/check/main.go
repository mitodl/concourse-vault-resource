package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/mitodl/concourse-vault-resource/cmd"
	"github.com/mitodl/concourse-vault-resource/concourse"
)

// GET for kv2 versions (kv1 not possible and TODO others currently unsupported)
func main() {
	// initialize checkRequest and secretSource
	checkRequest := concourse.NewCheckRequest(os.Stdin)
	secretSource := checkRequest.Source.Secret

	// return immediately if secret unspecified in source
	if secretSource == (concourse.SecretSource{}) {
		// dummy check response
		dummyResponse := concourse.NewCheckResponse([]concourse.Version{concourse.Version{Version: "0"}})
		// format checkResponse into json
		if err := json.NewEncoder(os.Stdout).Encode(&dummyResponse); err != nil {
			log.Print("unable to marshal check response struct to JSON")
			log.Fatal(err)
		}

		return
	}

	// initialize vault client from concourse source
	vaultClient := helper.VaultClientFromSource(checkRequest.Source)

	// initialize vault secret from concourse params and invoke constructor
	secret := helper.VaultSecretFromParams(secretSource.Mount, secretSource.Engine, secretSource.Path)
	secret.New()

	// retrieve version for secret
	_, getVersion, _, err := secret.SecretValue(vaultClient)
	if err != nil {
		log.Fatalf("version could not be retrieved for %s engine, %s mount, and path %s secret", secretSource.Engine, secretSource.Mount, secretSource.Path)
	}

	// assign input and get version and initialize versions slice TODO supporting other engines impacts this and next block greatly
	getVersionInt, _ := strconv.Atoi(getVersion)
	inputVersion, _ := strconv.Atoi(checkRequest.Version.Version)
	versions := []concourse.Version{}

	if inputVersion > getVersionInt {
		log.Printf("the input version %d is later than the retrieved version %s", inputVersion, getVersion)
		log.Print("only the retrieved version will be returned to Concourse")

		versions = []concourse.Version{concourse.Version{Version: getVersion}}
	} else {
		// populate versions slice with delta
		for versionDelta := inputVersion; versionDelta <= getVersionInt; versionDelta++ {
			versions = append(versions, concourse.Version{Version: strconv.Itoa(versionDelta)})
		}
	}

	// input secret version to constructed response
	checkResponse := concourse.NewCheckResponse(versions)

	// format checkResponse into json
	if err := json.NewEncoder(os.Stdout).Encode(&checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
