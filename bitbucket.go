package bitbucket

var apiBaseURL = "https://bitbucket.org/api/2.0"

func GetApiBaseURL() string {
	return apiBaseURL
}

func SetApiBaseURL(url string) {
	apiBaseURL = url
}

type users interface {
	Get(username string) interface{}
	Followers(username string) interface{}
	Following(username string) interface{}
	Repositories(username string) interface{}
}

type user interface {
	Profile() interface{}
	Emails() interface{}
}

type pullrequests interface {
	Create(opt PullRequestsOptions) interface{}
	Update(opt PullRequestsOptions) interface{}
	List(opt PullRequestsOptions) interface{}
	Get(opt PullRequestsOptions) interface{}
	Activities(opt PullRequestsOptions) interface{}
	Activity(opt PullRequestsOptions) interface{}
	Commits(opt PullRequestsOptions) interface{}
	Patch(opt PullRequestsOptions) interface{}
	Diff(opt PullRequestsOptions) interface{}
	Merge(opt PullRequestsOptions) interface{}
	Decline(opt PullRequestsOptions) interface{}
}

type repository interface {
	Get(opt RepositoryOptions) interface{}
	Create(opt RepositoryOptions) interface{}
	Delete(opt RepositoryOptions) interface{}
	ListWatchers(opt RepositoryOptions) interface{}
	ListForks(opt RepositoryOptions) interface{}
}

type repositories interface {
	ListForAccount(opt RepositoriesOptions) interface{}
	ListForTeam(opt RepositoriesOptions) interface{}
	ListPublic() interface{}
}

type commits interface {
	GetCommits(opt CommitsOptions) interface{}
	GetCommit(opt CommitsOptions) interface{}
	GetCommitComments(opt CommitsOptions) interface{}
	GetCommitComment(opt CommitsOptions) interface{}
	GetCommitStatus(opt CommitsOptions) interface{}
	GiveApprove(opt CommitsOptions) interface{}
	RemoveApprove(opt CommitsOptions) interface{}
}

type branchrestrictions interface {
	Gets(opt BranchRestrictionsOptions) interface{}
	Get(opt BranchRestrictionsOptions) interface{}
	Create(opt BranchRestrictionsOptions) interface{}
	Update(opt BranchRestrictionsOptions) interface{}
	Delete(opt BranchRestrictionsOptions) interface{}
}

type diff interface {
	GetDiff(opt DiffOptions) interface{}
	GetPatch(opt DiffOptions) interface{}
}

type webhooks interface {
	Gets(opt WebhooksOptions) interface{}
	Get(opt WebhooksOptions) interface{}
	Create(opt WebhooksOptions) interface{}
	Update(opt WebhooksOptions) interface{}
	Delete(opt WebhooksOptions) interface{}
}

type teams interface {
	List(role string) interface{} // [WIP?] role=[admin|contributor|member]
	Profile(teamname string) interface{}
	Members(teamname string) interface{}
	Followers(teamname string) interface{}
	Following(teamname string) interface{}
	Repositories(teamname string) interface{}
}

type RepositoriesOptions struct {
	Owner string `json:"owner"`
	Team  string `json:"team"`
	Role  string `json:"role"` // role=[owner|admin|contributor|member]
}

type RepositoryOptions struct {
	Owner     string `json:"owner"`
	Repo_slug string `json:"repo_slug"`
	Scm       string `json:"scm"`
	//	Name        string `json:"name"`
	Is_private  string `json:"is_private"`
	Description string `json:"description"`
	Fork_policy string `json:"fork_policy"`
	Language    string `json:"language"`
	Has_issues  string `json:"has_issues"`
	Has_wiki    string `json:"has_wiki"`
}

type PullRequestsOptions struct {
	Id                  string   `json:"id"`
	Comment_id          string   `json:"comment_id"`
	Owner               string   `json:"owner"`
	Repo_slug           string   `json:"repo_slug"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Close_source_branch bool     `json:"close_source_branch"`
	Source_branch       string   `json:"source_branch"`
	Source_repository   string   `json:"source_repository"`
	Destination_branch  string   `json:"destination_branch"`
	Destination_commit  string   `json:"destination_repository"`
	Message             string   `json:"message"`
	Reviewers           []string `json:"reviewers"`
}

type CommitsOptions struct {
	Owner       string `json:"owner"`
	Repo_slug   string `json:"repo_slug"`
	Revision    string `json:"revision"`
	Branchortag string `json:"branchortag"`
	Include     string `json:"include"`
	Exclude     string `json:"exclude"`
	Comment_id  string `json:"comment_id"`
}

type BranchRestrictionsOptions struct {
	Owner     string            `json:"owner"`
	Repo_slug string            `json:"repo_slug"`
	Id        string            `json:"id"`
	Groups    map[string]string `json:"groups"`
	Pattern   string            `json:"pattern"`
	Users     []string          `json:"users"`
	Kind      string            `json:"kind"`
	Full_slug string            `json:"full_slug"`
	Name      string            `json:"name"`
}

type DiffOptions struct {
	Owner     string `json:"owner"`
	Repo_slug string `json:"repo_slug"`
	Spec      string `json:"spec"`
}

type WebhooksOptions struct {
	Owner       string   `json:"owner"`
	Repo_slug   string   `json:"repo_slug"`
	Uuid        string   `json:"uuid"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"` // EX) {'repo:push','issue:created',..} REF) https://goo.gl/VTj93b
}
