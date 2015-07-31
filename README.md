# go-bitbucket

It is Bitbucket-API library for golang.

Support Bitbucket-API Version 2.0. 

And the response type is json format as defined Bitbucket-API.

- ref) <https://confluence.atlassian.com/display/BITBUCKET/Version+2>

## Install

```
go get github.com/aqafiam/go-bitbucket
```

## How to use

```
import "github.com/aqafiam/go-bitbucket"
```


## Example

```
package main

import (
        "github.com/aqafiam/go-bitbucket" 
        "fmt"
)

func main() {

        c := bitbucket.NewBasicAuth("username", "password")

        opt := bitbucket.PullRequestsOptions{
                Id:         "4",
                Owner:      "username",
                Repo_slug:  "awesome-project",
                Source_branch: "develop",
                Destination_branch: "master",
                Title: "fix bug. #9999",
                Close_source_branch: true
        }
        res := c.Repositories.PullRequests.Create(opt)

        fmt.Println(res) // receive the data as json format
}
```

## License

MIT

## Author

[aqafiam](https://github.com/aqafiam)
