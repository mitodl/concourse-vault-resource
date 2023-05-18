package main

import "fmt"

// no PUT/POST associated with this custom resource
func main() {
	fmt.Print("{\"metadata\":[{}],\"version\":{\"version\":\"\"}}")
}
