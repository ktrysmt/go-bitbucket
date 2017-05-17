.PHONY: test

env:
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

test: env
	go test -v ./tests

