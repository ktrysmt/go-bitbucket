package openApiTests

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile(t *testing.T) {
	c := bitbucket.NewBasicAuth("username", "password")

	res, err := c.User.Profile()

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
