package bitbucket

type Users struct {
	c *Client
}

func (u *Users) Get(t string) interface{} {

	url := GetApiBaseUrl() + "/users/" + t + "/"
	return u.c.execute("GET", url, "")
}

func (c *Client) Get(t string) interface{} {

	url := GetApiBaseUrl() + "/users/" + t + "/"
	return c.execute("GET", url, "")
}

func (u *Users) Followers(t string) interface{} {

	url := GetApiBaseUrl() + "/users/" + t + "/followers"
	return u.c.execute("GET", url, "")
}

func (u *Users) Following(t string) interface{} {

	url := GetApiBaseUrl() + "/users/" + t + "/following"
	return u.c.execute("GET", url, "")
}
func (u *Users) Repositories(t string) interface{} {

	url := GetApiBaseUrl() + "/users/" + t + "/repositories"
	return u.c.execute("GET", url, "")
}
