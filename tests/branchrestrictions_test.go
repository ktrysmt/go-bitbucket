package tests

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/raphaeldevs/go-bitbucket"
)

var (
	user  = os.Getenv("BITBUCKET_TEST_USERNAME")
	pass  = os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner = os.Getenv("BITBUCKET_TEST_OWNER")
	repo  = os.Getenv("BITBUCKET_TEST_REPOSLUG")
)

func setup(t *testing.T) *bitbucket.Client {

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}
	if repo == "" {
		t.Error("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)
	return c
}

func TestBranchRestrictionsKindPush(t *testing.T) {

	c := setup(t)
	var pushRestrictionID int

	t.Run("create", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    owner,
			Pattern:  "develop",
			RepoSlug: repo,
			Kind:     "push",
			Users:    []string{user},
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
			Owner:    owner,
			RepoSlug: repo,
			ID:       strconv.Itoa(pushRestrictionID),
		}
		_, err := c.Repositories.BranchRestrictions.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestBranchRestrictionsKindRequirePassingBuilds(t *testing.T) {

	c := setup(t)
	var pushRestrictionID int

	t.Run("create", func(t *testing.T) {
		opt := &bitbucket.BranchRestrictionsOptions{
			Owner:    owner,
			Pattern:  "master",
			RepoSlug: repo,
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
			Owner:    owner,
			RepoSlug: repo,
			ID:       strconv.Itoa(pushRestrictionID),
		}
		_, err := c.Repositories.BranchRestrictions.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestBranchRestrictionsGets(t *testing.T) {
	c := setup(t)

	t.Run("gets", func(t *testing.T) {
		const expectedNumberOfRestrictions = 20
		var restrictionIDs []int

		defer func() {
			for i := 0; i < len(restrictionIDs); i++ {
				opt := &bitbucket.BranchRestrictionsOptions{
					Owner:    owner,
					RepoSlug: repo,
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
				Owner:    owner,
				Pattern:  fmt.Sprintf("branch-restrictions-gets-%d", i),
				RepoSlug: repo,
				Kind:     "push",
				Users:    []string{user},
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
			Owner:    owner,
			RepoSlug: repo,
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
