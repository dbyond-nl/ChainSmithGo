package main

import (
	"os"
	"testing"
)

func TestMainExecution(t *testing.T) {
	tempConfig := "test_config.yml"
	configContent := `
root_ca_path: "test_rootCA.crt"
intermediate_ca_path: "test_intermediateCA.crt"
certificates:
  server:
    cert_path: "test_server.crt"
    key_path: "test_server.key"
    common_name: "server.local"
  client:
    cert_path: "test_client.crt"
    key_path: "test_client.key"
    common_name: "client.local"
`

	err := os.WriteFile(tempConfig, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove(tempConfig)

	if err := run(); err != nil {
		t.Fatalf("Application execution failed: %v", err)
	}
}
