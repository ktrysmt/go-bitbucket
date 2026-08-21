module github.com/ktrysmt/go-bitbucket

go 1.25.0

// You can uncomment this for local testing and development.
// Ref: https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/
//replace (
//	github.com/ktrysmt/go-bitbucket => ./
//	github.com/ktrysmt/go-bitbucket/tests => ./tests
//)

require (
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mitchellh/mapstructure v1.5.0
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
	golang.org/x/net v0.57.0
	golang.org/x/oauth2 v0.36.0
)

require (
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
