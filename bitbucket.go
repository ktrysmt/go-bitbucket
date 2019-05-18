package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var apiBaseURL = "https://api.bitbucket.org/2.0"

func GetApiBaseURL() string {
	return apiBaseURL
}

func SetApiBaseURL(urlStr string) {
	apiBaseURL = urlStr
}

type users interface {
	Get(username string) (interface{}, error)
	Followers(username string) (interface{}, error)
	Following(username string) (interface{}, error)
	Repositories(username string) (interface{}, error)
}

type user interface {
	Profile() (interface{}, error)
	Emails() (interface{}, error)
}

type pullrequests interface {
	Create(opt CreatePullRequestOpts) (interface{}, error)
	Update(opt CreatePullRequestOpts) (interface{}, error)
	List(opt CreatePullRequestOpts) (interface{}, error)
	Get(opt CreatePullRequestOpts) (interface{}, error)
	Activities(opt CreatePullRequestOpts) (interface{}, error)
	Activity(opt CreatePullRequestOpts) (interface{}, error)
	Commits(opt CreatePullRequestOpts) (interface{}, error)
	Patch(opt CreatePullRequestOpts) (interface{}, error)
	Diff(opt CreatePullRequestOpts) (interface{}, error)
	Merge(opt CreatePullRequestOpts) (interface{}, error)
	Decline(opt CreatePullRequestOpts) (interface{}, error)
}

type repository interface {
	Get(opt RepositoryOptions) (*Repository, error)
	Create(opt RepositoryOptions) (*Repository, error)
	Delete(opt RepositoryOptions) (interface{}, error)
	ListWatchers(opt RepositoryOptions) (interface{}, error)
	ListForks(opt RepositoryOptions) (interface{}, error)
	UpdatePipelineConfig(opt RepositoryPipelineOptions) (*Pipeline, error)
	AddPipelineVariable(opt RepositoryPipelineVariableOptions) (*PipelineVariable, error)
	AddPipelineKeyPair(opt RepositoryPipelineKeyPairOptions) (*PipelineKeyPair, error)
}

type repositories interface {
	ListForAccount(opt RepositoriesOptions) (interface{}, error)
	ListForTeam(opt RepositoriesOptions) (interface{}, error)
	ListPublic() (interface{}, error)
}

type commits interface {
	GetCommits(opt CommitsOptions) (interface{}, error)
	GetCommit(opt CommitsOptions) (interface{}, error)
	GetCommitComments(opt CommitsOptions) (interface{}, error)
	GetCommitComment(opt CommitsOptions) (interface{}, error)
	GetCommitStatus(opt CommitsOptions) (interface{}, error)
	GiveApprove(opt CommitsOptions) (interface{}, error)
	RemoveApprove(opt CommitsOptions) (interface{}, error)
	CreateCommitStatus(cmo CommitsOptions, cso CommitStatusOptions) (interface{}, error)
}

type branchrestrictions interface {
	Gets(opt BranchRestrictionsOptions) (interface{}, error)
	Get(opt BranchRestrictionsOptions) (interface{}, error)
	Create(opt BranchRestrictionsOptions) (interface{}, error)
	Update(opt BranchRestrictionsOptions) (interface{}, error)
	Delete(opt BranchRestrictionsOptions) (interface{}, error)
}

type diff interface {
	GetDiff(opt DiffOptions) (interface{}, error)
	GetPatch(opt DiffOptions) (interface{}, error)
}

type webhooks interface {
	Gets(opt WebhooksOptions) (interface{}, error)
	Get(opt WebhooksOptions) (interface{}, error)
	Create(opt WebhooksOptions) (interface{}, error)
	Update(opt WebhooksOptions) (interface{}, error)
	Delete(opt WebhooksOptions) (interface{}, error)
}

type teams interface {
	List(role string) (interface{}, error) // [WIP?] role=[admin|contributor|member]
	Profile(teamname string) (interface{}, error)
	Members(teamname string) (interface{}, error)
	Followers(teamname string) (interface{}, error)
	Following(teamname string) (interface{}, error)
	Repositories(teamname string) (interface{}, error)
	Projects(teamname string) (interface{}, error)
}

type ListOptions struct {
	Page    uint64 `json:"page"`
	PageLen uint64 `json:"pagelen"`
}

type RepositoriesOptions struct {
	ListOptions *ListOptions `json:"list_options"`
	Owner       string       `json:"owner"`
	Role        string       `json:"role"` // role=[owner|admin|contributor|member]
}

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
	Project     string `json:"project"`
}

type CommitsOptions struct {
	Owner       string `json:"owner"`
	RepoSlug    string `json:"repo_slug"`
	Revision    string `json:"revision"`
	Branchortag string `json:"branchortag"`
	Include     string `json:"include"`
	Exclude     string `json:"exclude"`
	CommentID   string `json:"comment_id"`
}

type CommitStatusOptions struct {
	Key         string `json:"key"`
	Url         string `json:"url"`
	State       string `json:"state"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
	Value    interface{}       `json:"value"`
}

type DiffOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Spec     string `json:"spec"`
}

type WebhooksOptions struct {
	Owner       string   `json:"owner"`
	RepoSlug    string   `json:"repo_slug"`
	Uuid        string   `json:"uuid"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"` // EX) {'repo:push','issue:created',..} REF) https://goo.gl/VTj93b
}

type RepositoryPipelineOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Enabled  bool   `json:"has_pipelines"`
}

type RepositoryPipelineVariableOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Uuid     string `json:"uuid"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Secured  bool   `json:"secured"`
}

type RepositoryPipelineKeyPairOptions struct {
	Owner      string `json:"owner"`
	RepoSlug   string `json:"repo_slug"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type DownloadsOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	FilePath string `json:"filepath"`
	FileName string `json:"filename"`
}

type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = "failed to parse unknown error format"
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}
