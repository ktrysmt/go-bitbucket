package bitbucket

// TODO: This file is WIP. The JSON response from the /activity endpoint is fluid so need a strategy on how to best manage this

type PullRequestActivities struct{}

func (p *PullRequestsService) Activities(owner, repo string, po *CreatePullRequestOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + owner + "/" + repo + "/pullrequests/activity"
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) Activity(owner, repo, id string, po *CreatePullRequestOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + owner + "/" + repo + "/pullrequests/" + id + "/activity"
	return p.client.execute("GET", urlStr, "", "")
}
