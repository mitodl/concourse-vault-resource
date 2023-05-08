package concourse

import (
  "log"
)

// TODO: may need to convert engines to enum struct. and then some of these would be slices; unhappy with how this is global, but the typedef is needed in vault package; honestly could I just do a nested struct or something? (then move defaults to individual constructors?); could secrets params be a slice of secrets?
type Params struct {
  // client
  VaultAddr    string
  AWSMountPath string
  Token        string
  Insecure     bool
  // secrets
  SecretEngine  SecretEngine
  SecretPath    string
  MountPath     string
}

// params constructor with defaults; TODO: secrets
func NewParams(inputParams Params) *Params {
  // initialize params return
  outputParams := new(Params)

  // vault address
  if len(inputParams.VaultAddr) == 0 {
    outputParams.VaultAddr = "http://127.0.0.1:8200"
  } else {
    outputParams.VaultAddr = inputParams.VaultAddr
  }

  // validate authentication inputs
  if len(inputParams.Token) == 0 && len(inputParams.AWSMountPath) == 0 {
    log.Fatal("Token and AWS authentication were simultaneously selected; these are mutually exclusive options")
  }
  if len(inputParams.Token) > 0 && len(inputParams.Token) != 26 {
    log.Fatal("the specified Vault token is invalid")
  }

  // aws mount path
  if len(inputParams.AWSMountPath) == 0 {
    outputParams.AWSMountPath = "aws"
  } else {
    outputParams.AWSMountPath = inputParams.AWSMountPath
  }

  return outputParams
}
