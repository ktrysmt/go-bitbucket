package bitbucket

import (
//"github.com/k0kubun/pp"
)

type User struct {
	c *Client
}

func (u *User) Profile() interface{} {

	url := API_BASE_URL + "/user/"
	return u.c.execute("GET", url, "")
}

func (u *User) Emails() interface{} {

	url := API_BASE_URL + "/user/emails"
	return u.c.execute("GET", url, "")
}
