package bitbucket

type Users struct {
	c *Client
}

func (u *Users) Get(t string) (interface{}, error) {

	urlStr := GetApiBaseURL() + "/users/" + t + "/"
	return u.c.execute("GET", urlStr, "")
}

func (c *Client) Get(t string) (interface{}, error) {

	urlStr := GetApiBaseURL() + "/users/" + t + "/"
	return c.execute("GET", urlStr, "")
}

func (u *Users) Followers(t string) (interface{}, error) {

	urlStr := GetApiBaseURL() + "/users/" + t + "/followers"
	return u.c.execute("GET", urlStr, "")
}

func (u *Users) Following(t string) (interface{}, error) {

	urlStr := GetApiBaseURL() + "/users/" + t + "/following"
	return u.c.execute("GET", urlStr, "")
}
func (u *Users) Repositories(t string) (interface{}, error) {

	urlStr := GetApiBaseURL() + "/users/" + t + "/repositories"
	return u.c.execute("GET", urlStr, "")
}
