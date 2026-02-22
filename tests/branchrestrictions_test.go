package tests

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
)

func TestBranchRestrictionsKindPush(t *testing.T) {

	c, err := setupBasicAuthTest(t)
	if err != nil {
		assert.Nilf(t, err, "failed to setup basic auth test: %w", err)
	}
	var pushRestrictionID int

	t.Run("create", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    ownerEnv,
			Pattern:  "develop",
			RepoSlug: repoEnv,
			Kind:     "push",
		}
		res, err := c.Repositories.BranchRestrictions.Create(opt)
		if err != nil {
			t.Error(err)
		}
		if res.Kind != "push" {
			t.Error("did not match branchrestriction kind")
		}
		pushRestrictionID = res.ID
	})

	t.Run("delete", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    ownerEnv,
			RepoSlug: repoEnv,
			ID:       strconv.Itoa(pushRestrictionID),
		}
		_, err := c.Repositories.BranchRestrictions.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestBranchRestrictionsKindRequirePassingBuilds(t *testing.T) {

	c, err := setupBasicAuthTest(t)
	if err != nil {
		assert.Nilf(t, err, "failed to setup basic auth test: %w", err)
	}
	var pushRestrictionID int

	t.Run("create", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    ownerEnv,
			Pattern:  "master",
			RepoSlug: repoEnv,
			Kind:     "require_passing_builds_to_merge",
			Value:    2,
		}
		res, err := c.Repositories.BranchRestrictions.Create(opt)
		if err != nil {
			t.Error(err)
		}
		if res.Kind != "require_passing_builds_to_merge" {
			t.Error("did not match branchrestriction kind")
		}
		pushRestrictionID = res.ID
	})

	t.Run("delete", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    ownerEnv,
			RepoSlug: repoEnv,
			ID:       strconv.Itoa(pushRestrictionID),
		}
		_, err := c.Repositories.BranchRestrictions.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestBranchRestrictionsGets(t *testing.T) {
	c, err := setupBasicAuthTest(t)
	if err != nil {
		assert.Nilf(t, err, "failed to setup basic auth test: %w", err)
	}

	t.Run("gets", func(t *testing.T) {
		const expectedNumberOfRestrictions = 20
		var restrictionIDs []int

		defer func() {
			for i := 0; i < len(restrictionIDs); i++ {
				opt := &bitbucket.BranchRestrictionsOptions{
					Owner:    ownerEnv,
					RepoSlug: repoEnv,
					ID:       strconv.Itoa(restrictionIDs[i]),
				}
				_, err := c.Repositories.BranchRestrictions.Delete(opt)
				if err != nil {
					t.Error(err)
				}
			}
		}()

		for i := 0; i < expectedNumberOfRestrictions; i++ {
			opt := &bitbucket.BranchRestrictionsOptions{
				Owner:    ownerEnv,
				Pattern:  fmt.Sprintf("branch-restrictions-gets-%d", i),
				RepoSlug: repoEnv,
				Kind:     "push",
			}
			res, err := c.Repositories.BranchRestrictions.Create(opt)
			if err != nil {
				t.Error(err)
				return
			}

			restrictionIDs = append(restrictionIDs, res.ID)
		}

		c.Pagelen = 5

		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    ownerEnv,
			RepoSlug: repoEnv,
		}

		res, err := c.Repositories.BranchRestrictions.Gets(opt)
		if err != nil {
			t.Error(err)
			return
		}

		responseMap, ok := res.(map[string]interface{})
		if !ok {
			t.Error(errors.New("response could not be decoded"))
			return
		}

		values := responseMap["values"].([]interface{})
		if len(values) != expectedNumberOfRestrictions {
			t.Error(fmt.Errorf("Expected %d branch restrictions but got %d. Response: %v", expectedNumberOfRestrictions, len(values), res))
			return
		}
	})
}
