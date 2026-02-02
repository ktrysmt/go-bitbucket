package tests

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net"
	"time"
)

// fetchCACerts connects to the given host:port and returns the CA certs in PEM format
func fetchCACerts(host string, port string) ([]byte, error) {
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

/*
// fetchCACertFromHost connects to the host and retrieves the root CA certificate in PEM format
func fetchCACertFromHost(host string) ([]byte, error) {
	// Ensure host has port
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = net.JoinHostPort(host, "443")
	}

	// TLS configuration without skipping verification
	conf := &tls.Config{
		InsecureSkipVerify: false, // verify certs
	}

	// Connect with timeout
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Get the verified certificate chain
	state := conn.ConnectionState()
	if len(state.VerifiedChains) == 0 {
		return nil, fmt.Errorf("no verified certificate chains found")
	}

	// The last certificate in the chain is usually the root CA
	rootCert := state.VerifiedChains[0][len(state.VerifiedChains[0])-1]

	// Encode to PEM
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCert.Raw,
	})

	return pemBytes, nil
}

// fetchCACertFromHost connects to the host and retrieves the root CA certificate in PEM format
func fetchCACertFromHost(host, port string) ([]byte, error) {
	// Ensure host has port
	host = net.JoinHostPort(host, port)

	// TLS configuration without skipping verification
	conf := &tls.Config{
		InsecureSkipVerify: false, // verify certs
	}

	// Connect with timeout
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Get the verified certificate chain
	state := conn.ConnectionState()
	if len(state.VerifiedChains) == 0 {
		return nil, fmt.Errorf("no verified certificate chains found")
	}

	// The last certificate in the chain is usually the root CA
	rootCert := state.VerifiedChains[0][len(state.VerifiedChains[0])-1]

	// Encode to PEM
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCert.Raw,
	})

	return pemBytes, nil
}
*/
