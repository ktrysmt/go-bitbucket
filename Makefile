.PHONY: test help
.DEFAULT_GOAL := help

help: ## print this help 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

env: ## check env for e2e testing
	ifndef BITBUCKET_TEST_USERNAME
	  $(error `BITBUCKET_TEST_USERNAME` is not set)
	endif
	ifndef BITBUCKET_TEST_PASSWORD
	  $(error `BITBUCKET_TEST_PASSWORD` is not set)
	endif
	ifndef BITBUCKET_TEST_OWNER
	  $(error `BITBUCKET_TEST_OWNER` is not set)
	endif
	ifndef BITBUCKET_TEST_REPOSLUG
	  $(error `BITBUCKET_TEST_REPOSLUG` is not set)
	endif

test: env ## run go test all
	go test -v ./tests

