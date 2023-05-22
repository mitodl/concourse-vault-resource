package main

import (
	"os"
	_ "testing"
)

func ExampleMain() {
	// deliver test pipeline file content as stdin to "in" the same as actual pipeline execution
	os.Stdin, _ = os.OpenFile("fixtures/token_kv.json", os.O_RDONLY, 0o644)
	defer os.Stdin.Close()
	os.Args[1] = "/opt/resource"

	// invoke main and validate stdout
	main()
	// Output: {"metadata":[],"version":{"version":""}}
}
