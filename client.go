package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Auth         *auth
	Repositories *Repositories
	Users        *users
	User         *user
	Teams        *teams
}

type auth struct {
	app_id, secret string
	user, password string
}

func NewOAuth(i, s string) *Client {
	a := &auth{app_id: i, secret: s}
	c := &Client{Auth: a}
	return c
}

func NewBasicAuth(u, p string) *Client {
	a := &auth{user: u, password: p}
	c := &Client{Auth: a}
	return c
}

func (c *Client) execute(method, url, text string) interface{} {

	body := strings.NewReader(text)
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")

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
	return API_BASE_URL + fmt.Sprintf(template, args...)
}
