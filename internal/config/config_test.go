package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
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

	cfg, err := LoadConfig(tempConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.RootCAPath != "test_rootCA.crt" {
		t.Errorf("Expected root_ca_path to be 'test_rootCA.crt', got %s", cfg.RootCAPath)
	}

	if _, exists := cfg.Certificates["server"]; !exists {
		t.Error("Expected server certificate configuration, but not found")
	}
}
