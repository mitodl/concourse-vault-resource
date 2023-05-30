package main

import (
	"os"
	_ "testing"
)

func ExampleMain() {
	// defer stdin close and establish workdir from argsp[1]
	defer os.Stdin.Close()
	os.Args[1] = "/opt/resource"

	// params secrets and source secret
	for _, secretKey := range []string{"params", "source"} {
		// deliver test pipeline file content as stdin to "in" the same as actual pipeline execution
		os.Stdin, _ = os.OpenFile("fixtures/token_kv_"+secretKey+".json", os.O_RDONLY, 0o644)

		// invoke main and validate stdout TODO validate vault.json content
		main()
		// Output: {"metadata":[{}],"version":{}}
	}
}
