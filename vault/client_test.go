package vault

import (
	"testing"
)

// global test helpers
const (
	testVaultAddress = "http://127.0.0.1:8200"
	testVaultToken   = "hvs.IZrMVkhZTIYgyArfMEhmLXsP"
)

var basicVaultConfig = &VaultConfig{
	VaultAddr: testVaultAddress,
	Token:     testVaultToken,
	Insecure:  true,
}

// test config constructor
func TestNewVaultConfig(test *testing.T) {
	basicVaultConfig.New()

	if basicVaultConfig.Engine != token || basicVaultConfig.VaultAddr != testVaultAddress || basicVaultConfig.AWSMountPath != "aws"  || basicVaultConfig.Token != testVaultToken || !basicVaultConfig.Insecure {
		test.Error("the Vault config constructor returned unexpected values.")
		test.Errorf("expected Auth Engine: %s, actual: %s", token, basicVaultConfig.Engine)
		test.Errorf("expected Vault Address: %s, actual: %s", testVaultAddress, basicVaultConfig.VaultAddr)
		test.Errorf("expected AWS Mount Path: aws, actual: %s", basicVaultConfig.AWSMountPath)
		test.Errorf("expected Vault Token: %s, actual: %s", testVaultToken, basicVaultConfig.Token)
		test.Errorf("expected Vault Insecure: %t, actual: %t", basicVaultConfig.Insecure, basicVaultConfig.Insecure)
	}
}

// test client error messages

// test client token authentication
func TestAuthClient(test *testing.T) {
  basicVaultClient := basicVaultConfig.authClient()

	if basicVaultClient.Token() != testVaultToken {
		test.Error("the authenticated Vault client return failed basic validation")
		test.Errorf("expected Vault token: %s, actual: %s", testVaultToken, basicVaultClient.Token())
	}
}
