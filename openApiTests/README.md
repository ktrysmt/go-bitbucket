# Open API Tests

The existing tests, under [/tests](../tests), have some shortcomings, the primary one being they are not run as part of a Continuous Integration (CI) process. Hence, they are not always run, and currently don't all pass (see issue [#111](https://github.com/ktrysmt/go-bitbucket/issues/111)).

The intention is to replace those tests with test that reference the swagger documentation from Bitbucket ([https://bitbucket.org/api/swagger.json](https://bitbucket.org/api/swagger.json)), using [Stoplight's Prism](https://stoplight.io/open-source/prism). And for those run as part of a CI process using [Github Actions](https://github.com/features/actions).

This will take time implement, so creating the new tests under this new folder and eventually, when ready will delete the existing `/tests` directory and will rename this one to `/tests`.

## Running tests locally

Run in a shell terminal:
```
docker run --rm -it -p 4010:4010 stoplight/prism:3 mock -h 0.0.0.0 https://bitbucket.org/api/swagger.json
```

Then in another shell terminal session run:
```
env BITBUCKET_API_BASE_URL=http://0.0.0.0:4010 go test -v ./openApiTests/...
```
