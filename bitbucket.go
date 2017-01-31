package bitbucket

var apiBaseURL = "https://bitbucket.org/api/2.0"

// GetAPIBaseURL returns the given URL.
func GetAPIBaseURL() string {
	return apiBaseURL
}

// SetAPIBaseURL recieves the URL.
func SetAPIBaseURL(url string) {
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

// RepositoriesOptions is configuration for resource
type RepositoriesOptions struct {
	Owner string `json:"owner"`
	Team  string `json:"team"`
	Role  string `json:"role"` // role=[owner|admin|contributor|member]
}

// RepositoryOptions is configuration for resource
type RepositoryOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Scm      string `json:"scm"`
	//	Name        string `json:"name"`
	IsPrivate   string `json:"is_private"`
	Description string `json:"description"`
	ForkPolicy  string `json:"fork_policy"`
	Language    string `json:"language"`
	HasIssues   string `json:"has_issues"`
	HasWiki     string `json:"has_wiki"`
}

// PullRequestsOptions is configuration for resource
type PullRequestsOptions struct {
	ID                string   `json:"id"`
	CommentID         string   `json:"comment_id"`
	Owner             string   `json:"owner"`
	RepoSlug          string   `json:"repo_slug"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	CloseSourceBranch bool     `json:"close_source_branch"`
	SourceBranch      string   `json:"source_branch"`
	SourceRepository  string   `json:"source_repository"`
	DestinationBranch string   `json:"destination_branch"`
	DestinationCommit string   `json:"destination_repository"`
	Message           string   `json:"message"`
	Reviewers         []string `json:"reviewers"`
}

// CommitsOptions is configuration for resource
type CommitsOptions struct {
	Owner       string `json:"owner"`
	RepoSlug    string `json:"repo_slug"`
	Revision    string `json:"revision"`
	Branchortag string `json:"branchortag"`
	Include     string `json:"include"`
	Exclude     string `json:"exclude"`
	CommentID   string `json:"comment_id"`
}

// BranchRestrictionsOptions is configuration for resource
type BranchRestrictionsOptions struct {
	Owner    string            `json:"owner"`
	RepoSlug string            `json:"repo_slug"`
	ID       string            `json:"id"`
	Groups   map[string]string `json:"groups"`
	Pattern  string            `json:"pattern"`
	Users    []string          `json:"users"`
	Kind     string            `json:"kind"`
	FullSlug string            `json:"full_slug"`
	Name     string            `json:"name"`
}

// DiffOptions is configuration for resource
type DiffOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Spec     string `json:"spec"`
}

// WebhooksOptions is configuration for resource
type WebhooksOptions struct {
	Owner       string   `json:"owner"`
	RepoSlug    string   `json:"repo_slug"`
	UUID        string   `json:"uuid"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"` // EX) {'repo:push','issue:created',..} REF) https://goo.gl/VTj93b
}
