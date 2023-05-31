package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// no PUT/POST associated with this custom resource NOW output actual versions instead of dummied
func main() {
	// initialize checkRequest and checkResponse
	checkRequest := concourse.NewCheckRequest(os.Stdin)
	checkResponse := concourse.NewCheckResponse()

	// dummy from request
	*checkResponse = []concourse.Version{checkRequest.Version}

	if err := json.NewEncoder(os.Stdout).Encode(checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
