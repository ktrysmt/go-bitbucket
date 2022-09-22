package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

func TestWebhook(t *testing.T) {
	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

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

	var webhookResourceUuid string

	t.Run("create", func(t *testing.T) {
		opt := &bitbucket.WebhooksOptions{
			Owner:       owner,
			RepoSlug:    repo,
			Description: "go-bb-test",
			Url:         "https://example.com",
			Active:      false,
			Events:      []string{bitbucket.RepoPushEvent, bitbucket.IssueCreatedEvent},
		}

		webhook, err := c.Repositories.Webhooks.Create(opt)
		if err != nil {
			t.Error(err)
		}

		if webhook == nil {
			t.Error("The webhook could not be created.")
		}

		if webhook.Description != "go-bb-test" {
			t.Error("The webhook `description` attribute does not match the expected value.")
		}
		if webhook.Url != "https://example.com" {
			t.Error("The webhook `url` attribute does not match the expected value.")
		}
		if webhook.Active != false {
			t.Error("The webhook `active` attribute does not match the expected value.")
		}
		if len(webhook.Events) != 2 {
			t.Error("The webhook `events` attribute does not match the expected value.")
		}

		webhookResourceUuid = webhook.Uuid
	})

	t.Run("get", func(t *testing.T) {
		opt := &bitbucket.WebhooksOptions{
			Owner:    owner,
			RepoSlug: repo,
			Uuid:     webhookResourceUuid,
		}
		webhook, err := c.Repositories.Webhooks.Get(opt)
		if err != nil {
			t.Error(err)
		}

		if webhook == nil {
			t.Error("The webhook could not be retrieved.")
		}

		if webhook.Description != "go-bb-test" {
			t.Error("The webhook `description` attribute does not match the expected value.")
		}
		if webhook.Url != "https://example.com" {
			t.Error("The webhook `url` attribute does not match the expected value.")
		}
		if webhook.Active != false {
			t.Error("The webhook `active` attribute does not match the expected value.")
		}
		if len(webhook.Events) != 2 {
			t.Error("The webhook `events` attribute does not match the expected value.")
		}
	})

	t.Run("update", func(t *testing.T) {
		opt := &bitbucket.WebhooksOptions{
			Owner:       owner,
			RepoSlug:    repo,
			Uuid:        webhookResourceUuid,
			Description: "go-bb-test-new",
			Url:         "https://new-example.com",
			Events:      []string{bitbucket.RepoPushEvent, bitbucket.IssueCreatedEvent, bitbucket.RepoForkEvent},
		}
		webhook, err := c.Repositories.Webhooks.Update(opt)
		if err != nil {
			t.Error(err)
		}

		if webhook == nil {
			t.Error("The webhook could not be retrieved.")
		}

		if webhook.Description != "go-bb-test-new" {
			t.Error("The webhook `description` attribute does not match the expected value.")
		}
		if webhook.Url != "https://new-example.com" {
			t.Error("The webhook `url` attribute does not match the expected value.")
		}
		if webhook.Active != false {
			t.Error("The webhook `active` attribute does not match the expected value.")
		}
		if len(webhook.Events) != 3 {
			t.Error("The webhook `events` attribute does not match the expected value.")
		}
	})

	t.Run("delete", func(t *testing.T) {
		opt := &bitbucket.WebhooksOptions{
			Owner:    owner,
			RepoSlug: repo,
			Uuid:     webhookResourceUuid,
		}
		_, err := c.Repositories.Webhooks.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("list/gets", func(t *testing.T) {
		const expectedNumberOfWebhooks = 20
		var webhookUUIDs []string
		defer func() {
			for _, uuid := range webhookUUIDs {
				_, err := c.Repositories.Webhooks.Delete(&bitbucket.WebhooksOptions{
					Owner:    owner,
					RepoSlug: repo,
					Uuid:     uuid,
				})

				if err != nil {
					t.Errorf("Failed to delete webhook UUID %s: %v", uuid, err)
				}
			}
		}()

		for i := 0; i < expectedNumberOfWebhooks; i++ {
			opt := &bitbucket.WebhooksOptions{
				Owner:       owner,
				RepoSlug:    repo,
				Description: fmt.Sprintf("go-bb-test-%d", i),
				Url:         fmt.Sprintf("https://example.com/%d", i),
				Active:      false,
				Events:      []string{bitbucket.RepoPushEvent, bitbucket.IssueCreatedEvent},
			}

			webhook, err := c.Repositories.Webhooks.Create(opt)
			if err != nil {
				t.Errorf("Failed to create webhook %d: %v", i, err)
			}

			webhookUUIDs = append(webhookUUIDs, webhook.Uuid)
		}

		// Use a page length of 5 to ensure the auto paging is working
		c.Pagelen = 5

		getsResponse, err := c.Repositories.Webhooks.Gets(&bitbucket.WebhooksOptions{
			Owner:    owner,
			RepoSlug: repo,
		})
		if err != nil {
			t.Errorf("Failed to list webhooks: %v", err)
			return
		}

		responseMap, ok := getsResponse.(map[string]interface{})
		if !ok {
			t.Error(errors.New("response could not be decoded"))
			return
		}

		values := responseMap["values"].([]interface{})
		if len(values) != expectedNumberOfWebhooks {
			t.Error(fmt.Errorf("Expected %d webhooks but got %d. Response: %v", expectedNumberOfWebhooks, len(values), getsResponse))
			return
		}

		listResponse, err := c.Repositories.Webhooks.List(&bitbucket.WebhooksOptions{
			Owner:    owner,
			RepoSlug: repo,
		})

		if err != nil {
			t.Errorf("Failed to list webhooks: %v", err)
			return
		}

		if len(listResponse) != expectedNumberOfWebhooks {
			t.Error(fmt.Errorf("Expected %d webhooks but got %d. Response: %v", expectedNumberOfWebhooks, len(values), listResponse))
			return
		}
	})
}
