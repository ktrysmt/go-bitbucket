package bitbucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.Error(t, err)
	assert.Equal(t, "", err.Error())
}

func TestDecodeError_MissingErrorKey(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{}

	err := DecodeError(errorMap)

	require.Error(t, err)
	assert.Equal(t, "", err.Error())
}

func TestDecodeError_InvalidStructure(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{
		"error": "not a map",
	}

	err := DecodeError(errorMap)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected a map")
}

func TestDecodeError_WithFields(t *testing.T) {
	t.Parallel()
	errorMap := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "Bad request",
			"fields": map[string][]string{
				"username": {"Username is required"},
			},
		},
	}

	err := DecodeError(errorMap)

	require.Error(t, err)
	assert.Equal(t, "Bad request", err.Error())
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
