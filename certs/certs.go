package certs

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetClientWithCerts() *http.Client {
	certsDir := filepath.Base(os.Getenv("CERTS_LOCATION"))

	// Create a CA certificate pool and add certs from certsDir
	caCertPool := x509.NewCertPool()
	caCertFiles, err := os.ReadDir(certsDir)
	if err != nil {
		log.Fatalf("Failed to read certs directory: %v", err)
	}

	for _, certFile := range caCertFiles {
		certPath := filepath.Join(certsDir, certFile.Name())
		caCert, err := os.ReadFile(certPath)
		if err != nil {
			log.Fatalf("Failed to read cert file %s: %v", certPath, err)
		}
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatalf("Failed to append cert file %s to pool", certPath)
		}
	}

	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Disable certificate validation
	}

	// Create an HTTP client with the custom TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client
}
