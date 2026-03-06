package bitbucket

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fetchCACertsForTest connects to the given host:port and returns the CA certs in PEM format.
// This is a test-local copy to avoid an import cycle with the tests package.
func fetchCACertsForTest(host string, port string) ([]byte, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", net.JoinHostPort(host, port), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer func() { _ = conn.Close() }()

	peerCerts := conn.ConnectionState().PeerCertificates
	if len(peerCerts) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

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

func TestAppendCaCerts_util_test(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring network access")
	}
	caCerts, err := fetchCACertsForTest("bitbucket.org", "443")
	if err != nil {
		t.Fatalf("Error fetching CA certs using `fetchCACertsForTest`: %v", err)
	}
	httpClient, err := appendCaCerts(caCerts)
	require.NoError(t, err)
	require.NotNil(t, httpClient)

	transport, ok := httpClient.Transport.(*http.Transport)
	require.True(t, ok, "Transport should be *http.Transport")
	require.NotNil(t, transport.TLSClientConfig, "TLSClientConfig should be set")
	assert.NotNil(t, transport.TLSClientConfig.RootCAs, "RootCAs cert pool should be set")
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion, "MinVersion should be TLS 1.2")
}

func TestAppendCaCerts_InvalidCert(t *testing.T) {
	t.Parallel()
	invalidPEM := []byte("this is not a valid PEM certificate")

	httpClient, err := appendCaCerts(invalidPEM)

	assert.Error(t, err)
	assert.Nil(t, httpClient)
	assert.Contains(t, err.Error(), "unable to append CA Certs")
}
