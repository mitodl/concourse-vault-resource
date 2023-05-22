package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// no PUT/POST associated with this custom resource
func main() {
	// initialize checkResponse
	checkResponse := concourse.CheckResponse([]concourse.Version{{Version: "1"}})

	if err := json.NewEncoder(os.Stdout).Encode(checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
