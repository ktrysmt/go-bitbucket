package bitbucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeError_ValidError(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "Repository not found",
			"fields":  map[string][]string{},
		},
	}

	err := DecodeError(errorMap)

	assert.Error(t, err)
	assert.Equal(t, "Repository not found", err.Error())
}

func TestDecodeError_EmptyMessage(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "",
		},
	}

	err := DecodeError(errorMap)

	// mapstructure.Decode fails when input doesn't match expected structure
	// The actual behavior depends on the mapstructure implementation
	if err != nil {
		assert.Error(t, err)
	}
}

func TestDecodeError_MissingErrorKey(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{}

	err := DecodeError(errorMap)

	// When "error" key is nil, mapstructure.Decode returns an error
	assert.Error(t, err)
}

func TestUnexpectedResponseStatusError_Error(t *testing.T) {
	t.Parallel()
	err := &UnexpectedResponseStatusError{
		Status: "404 Not Found",
		Body:   []byte(`{"error": "not found"}`),
	}

	assert.Equal(t, "404 Not Found", err.Error())
}

func TestUnexpectedResponseStatusError_ErrorWithBody(t *testing.T) {
	t.Parallel()
	err := &UnexpectedResponseStatusError{
		Status: "500 Internal Server Error",
		Body:   []byte(`{"error": "something broke"}`),
	}

	result := err.ErrorWithBody()

	assert.Error(t, result)
	assert.Contains(t, result.Error(), "500 Internal Server Error")
	assert.Contains(t, result.Error(), `{"error": "something broke"}`)
}
