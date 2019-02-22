package bitbucket

// IssuesService handles communication with the issue related methods
// of the Bitbucket API.
//
// Bitbucket API docs: https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/issues
type IssuesService struct {
	client *Client
}

type Issues struct {
	Page    int      `json:"page,omitempty"`
	Size    int      `json:"size,omitempty"`
	PageLen int      `json:"pagelen,omitempty"`
	Values  *[]Issue `json:"values,omitempty"`
}

type Issue struct {
	Priority string `json:"priority,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Links    struct {
		HTML struct {
			Href string `json:"href,omitempty"`
		} `json:"html,omitempty"`
	} `json:"links,omitempty"`
	Title     string       `json:"title,omitempty"`
	Votes     int          `json:"votes,omitempty"`
	Watches   int          `json:"watches,omitempty"`
	Content   IssueContent `json:"content,omitempty"`
	State     string       `json:"state,omitempty"`
	IssueType string       `json:"type,omitempty"`
	ID        int64        `json:"id,omitempty"`
}

type IssueContent struct {
	Raw    string `json:"raw,omitempty"`
	Markup string `json:"markup,omitempty"`
	Html   string `json:"html,omitempty"`
	Type   string `json:"type,omitempty"`
}

type CreateIssueOpts struct {
	Title    string                   `json:"title,omitempty"`
	Kind     string                   `json:"kind,omitempty"`
	Priority string                   `json:"priority,omitempty"`
	Content  *CreateIssueContentOpts  `json:"content,omitempty"`
	Assignee *CreateIssueAssigneeOpts `json:"assignee,omitempty"`
}

type CreateIssueContentOpts struct {
	Raw *string `json:"raw,omitempty"`
}

type CreateIssueAssigneeOpts struct {
	Username *string `json:"username,omitempty"`
}

func (i *IssuesService) List(owner, repoSlug string) (*Issues, *Response, error) {
	result := new(Issues)
	urlStr := i.client.requestUrl("/repositories/%s/%s/issues", owner, repoSlug)

	response, err := i.client.executeNew("GET", urlStr, result, nil, "")

	return result, response, err
}

func (i *IssuesService) Get(owner, repoSlug, issueId string) (*Issue, *Response, error) {
	result := new(Issue)
	urlStr := i.client.requestUrl("/repositories/%s/%s/issues/%s", owner, repoSlug, issueId)
	response, err := i.client.executeNew("GET", urlStr, result, nil, "")

	return result, response, err
}

func (i *IssuesService) Create(owner, repoSlug string, io *CreateIssueOpts) (*Issue, *Response, error) {
	result := new(Issue)
	urlStr := i.client.requestUrl("/repositories/%s/%s/issues", owner, repoSlug)
	response, err := i.client.executeNew("POST", urlStr, result, io, "")

	return result, response, err
}
