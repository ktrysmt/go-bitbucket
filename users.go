package bitbucket

type Users struct {
	c *Client
}

func (u *Users) Get(t string) interface{} {

	url := API_BASE_URL + "/users/" + t + "/"
	return u.c.execute("GET", url, "")
}

func (c *Client) Get(t string) interface{} {

	url := API_BASE_URL + "/users/" + t + "/"
	return c.execute("GET", url, "")
}

func (u *Users) Followers(t string) interface{} {

	url := API_BASE_URL + "/users/" + t + "/followers"
	return u.c.execute("GET", url, "")
}

func (u *Users) Following(t string) interface{} {

	url := API_BASE_URL + "/users/" + t + "/following"
	return u.c.execute("GET", url, "")
}
func (u *Users) Repositories(t string) interface{} {

	url := API_BASE_URL + "/users/" + t + "/repositories"
	return u.c.execute("GET", url, "")
}
