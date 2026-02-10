module github.com/ktrysmt/go-bitbucket

go 1.24.0

// You can uncomment this for local testing and development.
// Ref: https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/
//replace (
//	github.com/ktrysmt/go-bitbucket => ./
//	github.com/ktrysmt/go-bitbucket/tests => ./tests
//)

require (
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mitchellh/mapstructure v1.5.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
	golang.org/x/net v0.49.0
	golang.org/x/oauth2 v0.34.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
