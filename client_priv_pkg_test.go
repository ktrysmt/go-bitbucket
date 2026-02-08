package bitbucket

import (
	"testing"

	"github.com/ktrysmt/go-bitbucket/tests"
	"github.com/stretchr/testify/assert"
)

func TestAppendCaCerts_util_test(t *testing.T) {
	caCerts, err := tests.FetchCACerts("bitbucket.org", "443")
	if err != nil {
		t.Fatal(err)
	}
	httpClient, err := appendCaCerts(caCerts)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, httpClient)
}
