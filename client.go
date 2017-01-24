package bitbucket

import (
	"encoding/json"
	"fmt"
	//	"github.com/k0kubun/pp"
	//	"os"

	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Auth         *auth
	Users        users
	User         user
	Teams        teams
	Repositories *Repositories
}

type auth struct {
	app_id, secret string
	user, password string
}

func NewOAuth(i, s string) *Client {
	a := &auth{app_id: i, secret: s}
	return injectClient(a)
}

func NewBasicAuth(u, p string) *Client {
	a := &auth{user: u, password: p}
	return injectClient(a)
}

func injectClient(a *auth) *Client {
	c := &Client{Auth: a}
	c.Repositories = &Repositories{
		c:                  c,
		PullRequests:       &PullRequests{c: c},
		Repository:         &Repository{c: c},
		Commits:            &Commits{c: c},
		Diff:               &Diff{c: c},
		BranchRestrictions: &BranchRestrictions{c: c},
		Webhooks:           &Webhooks{c: c},
	}
	c.Users = &Users{c: c}
	c.User = &User{c: c}
	c.Teams = &Teams{c: c}
	return c
}

func (c *Client) execute(method, url, text string) interface{} {

	body := strings.NewReader(text)
	req, err := http.NewRequest(method, url, body)
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return err
	}

	if c.Auth.user != "" && c.Auth.password != "" {
		req.SetBasicAuth(c.Auth.user, c.Auth.password)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result interface{}
	json.Unmarshal(buf, &result)

	return result
}

func (c *Client) requestUrl(template string, args ...interface{}) string {

	if len(args) == 1 && args[0] == "" {
		return GetApiBaseUrl() + template
	} else {
		return GetApiBaseUrl() + fmt.Sprintf(template, args...)
	}
}
