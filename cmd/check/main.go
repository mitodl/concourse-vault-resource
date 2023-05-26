package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mitodl/concourse-vault-resource/concourse"
)

// no PUT/POST associated with this custom resource NOW output actual versions instead of dummied
func main() {
	// initialize checkResponse
	checkResponse := concourse.NewCheckResponse(concourse.Version{"mount-path": "0"})

	if err := json.NewEncoder(os.Stdout).Encode(checkResponse); err != nil {
		log.Print("unable to marshal check response struct to JSON")
		log.Fatal(err)
	}
}
