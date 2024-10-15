package tests

import (
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

func TestUserSSHKey(t *testing.T) {
	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	var sshKeyResourceUuid string

	label := "go-user-test"
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCjuh+EUAXrLtlQ5LfiSf4nWVOAuWUwMy+Cb+AkqFolyw/tuZh0tx9cEzSHddWgeFSJa5Zj0OYUEVTMkhUUfvKb8tfqbTGr0EWZrW6Odc6bqQXBDa48IfSqPYHmmdJh07MpRRQRqEMHB4WfnNuEUhOuNHr2lOX7BtCyp4r38gkuNBFmT6nheSoxSjJ6t3VbViyO+p2RY1RaGL77kUMgt4ti4MR4lNuUBT+BOxiILHqwWfY0z0i7Cc1zW4PvDbFtgHzSzQBdBel3vjk5AALZV31tiu0R21Gxm35n5L2N12ZgTXVXOVC1qfGzh6OR+7ZG0/iWyCmOoi+cOznXlnQEC/k5"
	t.Run("create", func(t *testing.T) {
		keyOptions := &bitbucket.SSHKeyOptions{
			Label: label,
			Key:   key,
			Owner: user,
		}
		sshUserKey, err := c.Users.SSHKeys.Create(keyOptions)
		if err != nil {
			t.Error(err)
		}
		if sshUserKey == nil {
			t.Error("The User SSH Key could not be created.")
		}
		sshKeyResourceUuid = sshUserKey.Uuid
	})
	t.Run("get", func(t *testing.T) {
		keyOptions := &bitbucket.SSHKeyOptions{
			Owner: user,
			Uuid:  sshKeyResourceUuid,
		}
		sshKey, err := c.Users.SSHKeys.Get(keyOptions)
		if err != nil {
			t.Error(err)
		}
		if sshKey == nil {
			t.Error("The Deploy Key could not be retrieved.")
		}
		if sshKey.Uuid != sshKeyResourceUuid {
			t.Error("The SSH Key `id` attribute does not match the expected value.")
		}
		if sshKey.Label != label {
			t.Error("The SSH Key `label` attribute does not match the expected value.")
		}
		if sshKey.Key != key {
			t.Error("The SSH Key `key` attribute does not match the expected value.")
		}
	})

	t.Run("delete", func(t *testing.T) {
		keyOptions := &bitbucket.SSHKeyOptions{
			Owner: user,
			Uuid:  sshKeyResourceUuid,
		}
		_, err := c.Users.SSHKeys.Delete(keyOptions)
		if err != nil {
			t.Error(err)
		}
	})
}
