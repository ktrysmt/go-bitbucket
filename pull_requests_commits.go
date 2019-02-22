package bitbucket

type PullRequestCommits struct {
	PageLen int                  `json:"pagelen,omitempty"`
	Page    int                  `json:"page,omitempty"`
	Values  *[]PullRequestCommit `json:"values,omitempty"`
}

type PullRequestCommit struct {
}

func (p *PullRequestsService) Commits(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/commits"
	return p.client.execute("GET", urlStr, "", "")
}
