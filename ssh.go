package bitbucket

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type SSHKeys struct {
	c *Client
}
type SSHKey struct {
	Id      int    `json:"id"`
	Label   string `json:"label"`
	Key     string `json:"key"`
	Comment string `json:"comment"`
}

func decodeSSHKey(response interface{}) (*SSHKey, error) {
	respMap := response.(map[string]interface{})

	if respMap["type"] == "error" {
		return nil, DecodeError(respMap)
	}

	var SSHKey = new(SSHKey)
	err := mapstructure.Decode(respMap, SSHKey)
	if err != nil {
		return nil, err
	}

	return SSHKey, nil
}

func buildSSHKeysBody(opt *SSHKeyCreationOptions) (string, error) {
	body := map[string]interface{}{}
	body["label"] = opt.Label
	body["key"] = opt.Key

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *SSHKeys) Create(opt *SSHKeyCreationOptions) (interface{}, error) {
	urlStr := s.c.requestUrl("/users/%s/ssh-keys/", opt.User)

	data, err := buildSSHKeysBody(opt)
	if err != nil {
		return nil, err
	}
	response, err := s.c.execute("POST", urlStr, data)
	if err != nil {
		return nil, err
	}

	return decodeSSHKey(response)
}
