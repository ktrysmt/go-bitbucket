package tests

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net"
	"time"
)

// fetchCACerts connects to the given host:port and returns the CA certs in PEM format
func FetchCACerts(host string, port string) ([]byte, error) {
	// Prepare TLS configuration (skip verification so we can inspect all certs)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // We only want to fetch certs, not verify
	}

	// Connect with a timeout
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", net.JoinHostPort(host, port), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Get the full certificate chain
	peerCerts := conn.ConnectionState().PeerCertificates
	if len(peerCerts) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	// Extract the root CA (last in chain) or all intermediates
	var pemData []byte
	for _, cert := range peerCerts {
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}
		pemData = append(pemData, pem.EncodeToMemory(block)...)
	}

	return pemData, nil
}
