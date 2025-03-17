package tls

import (
	"testing"
)

func TestGenerateCA(t *testing.T) {
	rootCert, rootKey, err := GenerateCA("test_rootCA.crt", "test_rootCA.key", nil, nil, true)
	if err != nil {
		t.Fatalf("Failed to generate Root CA: %v", err)
	}
	if rootCert == nil || rootKey == nil {
		t.Fatal("Root CA certificate or key is nil")
	}
}

func TestGenerateIntermediateCA(t *testing.T) {
	rootCert, rootKey, _ := GenerateCA("test_rootCA.crt", "test_rootCA.key", nil, nil, true)
	intermediateCert, intermediateKey, err := GenerateCA("test_intermediateCA.crt", "test_intermediateCA.key", rootCert, rootKey, false)
	if err != nil {
		t.Fatalf("Failed to generate Intermediate CA: %v", err)
	}
	if intermediateCert == nil || intermediateKey == nil {
		t.Fatal("Intermediate CA certificate or key is nil")
	}
}

func TestGenerateCert(t *testing.T) {
	rootCert, rootKey, _ := GenerateCA("test_rootCA.crt", "test_rootCA.key", nil, nil, true)
	intermediateCert, intermediateKey, _ := GenerateCA("test_intermediateCA.crt", "test_intermediateCA.key", rootCert, rootKey, false)

	testCertPath := "test_server.crt"
	testKeyPath := "test_server.key"
	err := GenerateCert(testCertPath, testKeyPath, intermediateCert, intermediateKey, "server.local")
	if err != nil {
		t.Fatalf("Failed to generate server certificate: %v", err)
	}
}
