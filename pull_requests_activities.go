package bitbucket

// TODO: This file is WIP. The JSON response from the /activity endpoint is fluid so need a strategy on how to best manage this

type PullRequestActivities struct{}

func (p *PullRequestsService) Activities(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/activity"
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) Activity(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/activity"
	return p.client.execute("GET", urlStr, "", "")
}
