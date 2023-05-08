package vault

import (
  "testing"

  "github.com/mitodl/concourse-vault-resource/concourse"
)

// global test helpers
const (
  testVaultAddress = "http://127.0.0.1:8200"
  testVaultToken   = "hvs.IZrMVkhZTIYgyArfMEhmLXsP"
)

var basicConcourseParams = &concourse.Params{
  VaultAddr: testVaultAddress,
  Token:     testVaultToken,
  Insecure:  true,
}

var basicVaultConfig = NewVaultConfig(basicConcourseParams)
var BasicVaultClient = basicVaultConfig.authClient()

// test config constructor
func TestNewVaultConfig(test *testing.T) {
  newVaultConfig := NewVaultConfig(basicConcourseParams)

  if newVaultConfig.vaultAddr != testVaultAddress || newVaultConfig.token != testVaultToken || newVaultConfig.insecure != basicConcourseParams.Insecure {
    test.Error("the Vault config constructor returned unexpected values.")
    test.Errorf("expected Vault Address: %s, actual: %s", testVaultAddress, newVaultConfig.vaultAddr)
    test.Errorf("expected Vault Token: %s, actual: %s", testVaultToken, newVaultConfig.token)
    test.Errorf("expected Vault Insecure: %t, actual: %t", basicConcourseParams.Insecure, newVaultConfig.insecure)
  }
}

// test client error messages

// test client token authentication
func TestAuthClient(test *testing.T) {
  vaultClient := basicVaultConfig.authClient()

  if vaultClient.Token() != testVaultToken {
    test.Error("the authenticated Vault client return failed basic validation")
    test.Errorf("expected Vault token: %s, actual: %s", testVaultToken, vaultClient.Token())
  }
}
