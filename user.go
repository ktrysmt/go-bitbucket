package bitbucket

import (
//"github.com/k0kubun/pp"
)

type User struct {
	c *Client
}

func (u *User) Profile() interface{} {

	url := GetApiBaseUrl() + "/user/"
	return u.c.execute("GET", url, "")
}

func (u *User) Emails() interface{} {

	url := GetApiBaseUrl() + "/user/emails"
	return u.c.execute("GET", url, "")
}
