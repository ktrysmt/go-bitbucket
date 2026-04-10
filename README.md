# go-bitbucket

> **This is a temporary fork of [ktrysmt/go-bitbucket](https://github.com/ktrysmt/go-bitbucket).**
> It exists to adopt Bitbucket Cloud API changes (deprecated endpoint migration, safer error handling)
> ahead of upstream. We intend to return to the upstream module once it incorporates these changes.

Bitbucket API v2.0 library for Go.

- Bitbucket API v2.0 <https://developer.atlassian.com/bitbucket/api/2/reference/resource/>
- Swagger for API v2.0 <https://api.bitbucket.org/swagger.json>

## Changes from upstream

This fork introduces the following breaking and non-breaking changes relative to `ktrysmt/go-bitbucket v0.9.x`:

### Breaking

- **Constructor signatures changed**: All client constructors (`NewBasicAuth`, `NewOAuthbearerToken`, `NewOAuthClientCredentials`, etc.) now return `(*Client, error)` instead of `*Client`. Callers must handle the error.
- `NewOAuthWithCode` and `NewOAuthWithRefreshToken` return three values: `(*Client, string, error)`.
- **`log.Fatal` removed**: All `log.Fatal` calls in constructors have been replaced with returned errors, so the library no longer terminates the calling process on failure.

### Non-breaking

- **Workspace listing**: Uses `GET /user/workspaces` instead of the deprecated `GET /workspaces`. Response decoding handles the `workspace_access` envelope from the new endpoint.
- **Custom CA certificate support**: New constructor variants (`NewBasicAuthWithCaCert`, `NewOAuthbearerTokenWithCaCert`, etc.) accept custom CA certificates for Bitbucket Server / Data Center deployments with self-signed certs.
- **Base URL constructors**: New variants (`NewBasicAuthWithBaseUrlStr`, `NewOAuthbearerTokenWithBaseUrlStr`, etc.) accept the API base URL as a parameter instead of relying on the `BITBUCKET_API_BASE_URL` environment variable.
- **`PullRequests.List()`**: Added as the properly named method; `Gets()` is kept as a backward-compatible alias.
- **`PullRequestsMergeStrategy` type**: New string enum for merge strategy options.
- **Pipeline variable methods**: `GetPipelineVariable` and `UpdatePipelineVariable` added to the repository interface.

## Install

```sh
go get github.com/trufflesecurity/go-bitbucket
```

## Usage

### Create a pull request

```go
package main

import (
        "fmt"
        "log"

        "github.com/trufflesecurity/go-bitbucket"
)

func main() {
        c, err := bitbucket.NewBasicAuth("username", "password")
        if err != nil {
                log.Fatal(err)
        }

        opt := &bitbucket.PullRequestsOptions{
                Owner:             "your-team",
                RepoSlug:          "awesome-project",
                SourceBranch:      "develop",
                DestinationBranch: "master",
                Title:             "fix bug. #9999",
                CloseSourceBranch: true,
        }

        res, err := c.Repositories.PullRequests.Create(opt)
        if err != nil {
                log.Fatal(err)
        }

        fmt.Println(res)
}
```

### Create a repository

```go
package main

import (
        "fmt"
        "log"

        "github.com/trufflesecurity/go-bitbucket"
)

func main() {
        c, err := bitbucket.NewBasicAuth("username", "password")
        if err != nil {
                log.Fatal(err)
        }

        opt := &bitbucket.RepositoryOptions{
                Owner:    "project_name",
                RepoSlug: "repo_name",
                Scm:      "git",
        }

        res, err := c.Repositories.Repository.Create(opt)
        if err != nil {
                log.Fatal(err)
        }

        fmt.Println(res)
}
```

## FAQ

### Support Bitbucket API v1.0 ?

It does not correspond yet. Because there are many differences between v2.0 and v1.0.

- Bitbucket API v1.0 <https://confluence.atlassian.com/bitbucket/version-1-423626337.html>

It is officially recommended to use v2.0.
But unfortunately Bitbucket Server (formerly: Stash) API is still v1.0.
And The API v1.0 covers resources that the v2.0 API and API v2.0 is yet to cover.

## Development

### Get dependencies

It's using `go mod`.

### How to testing

Set your available user account to Global Env.

```sh
export BITBUCKET_TEST_USERNAME=<your_username>
export BITBUCKET_TEST_PASSWORD=<your_password>
export BITBUCKET_TEST_OWNER=<your_repo_owner>
export BITBUCKET_TEST_REPOSLUG=<your_repo_name>
export BITBUCKET_TEST_ACCESS_TOKEN=<your_repo_access_token>
```

And just run;

```sh
make test
```

If you want to test individually;

```sh
go test -v ./tests/diff_test.go
```
E2E Integration tests;
```sh
make test/e2e
```

Unit tests;
```sh
make test/unit
```

Mock tests;

```sh
make test/mock
```
Individually;
```sh
go test ./mock_tests/repository_mock_test.go
```

For documented workflow of the go:qmock test structure in ```/mock_tests/repository_mock_test.go``` refer to;
- TestMockRepositoryPipelineVariables_List_Success
- TestMockRepositoryPipelineVariables_List_Error

## License

[Apache License 2.0](./LICENSE)

## Author

Originally created by [ktrysmt](https://github.com/ktrysmt). Forked and maintained by [TruffleSecurity](https://github.com/trufflesecurity).
