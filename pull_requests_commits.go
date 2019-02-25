package bitbucket

type PullRequestCommits struct {
	PageLen int                  `json:"pagelen,omitempty"`
	Page    int                  `json:"page,omitempty"`
	Values  *[]PullRequestCommit `json:"values,omitempty"`
}

type PullRequestCommit struct {
}

func (p *PullRequestsService) Commits(owner, repo, id string, po *CreatePullRequestOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + owner + "/" + repo + "/pullrequests/" + id + "/commits"
	return p.client.execute("GET", urlStr, "", "")
}
