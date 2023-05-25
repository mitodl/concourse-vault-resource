package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// no PUT/POST associated with this custom resource TODO maybe https://pkg.go.dev/github.com/hashicorp/vault/api#LifetimeWatcher
func main() {
	// initialize checkResponse
	checkResponse := concourse.CheckResponse([]concourse.Version{})

	if err := json.NewEncoder(os.Stdout).Encode(checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
