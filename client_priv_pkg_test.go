package bitbucket

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	caCerts, err := fetchCACertsForTest("bitbucket.org", "443")
	if err != nil {
		t.Fatalf("Error fetching CA certs using `fetchCACertsForTest`: %v", err)
	}
	httpClient, err := appendCaCerts(caCerts)
	if err != nil {
		t.Fatalf("Error returned from `appendCaCerts` failed to create the http client: %v", err)
	}
	assert.NotNil(t, httpClient)
}
