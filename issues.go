package bitbucket

import "net/url"

type Issues struct {
	c *Client
}

func (p *Issues) Gets(io *IssuesOptions) (interface{}, error) {
	url, err := url.Parse(p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/")
	if err != nil {
		return nil, err
	}

	if io.States != nil && len(io.States) != 0 {
		query := url.Query()
		for _, state := range io.States {
			query.Set("state", state)
		}
		url.RawQuery = query.Encode()
	}

	if io.Query != "" {
		query := url.Query()
		query.Set("q", io.Query)
		url.RawQuery = query.Encode()
	}

	if io.Sort != "" {
		query := url.Query()
		query.Set("sort", io.Sort)
		url.RawQuery = query.Encode()
	}

	return p.c.execute("GET", url.String(), "")
}
