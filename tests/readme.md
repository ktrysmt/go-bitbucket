## How to test go-bitbucket

### Set env for self testing

URL Syntax: `https://<your_username>:<your_password>@bitbucket.org/<your_repo_owner>/<your_repo_name>.git`

```
export BITBUCKET_TEST_USERNAME=<your_username> 
export BITBUCKET_TEST_PASSWORD=<your_password> 
export BITBUCKET_TEST_OWNER=<your_repo_owner>  
export BITBUCKET_TEST_REPOSLUG=<your_repo_name>
cd ./test
go test
```
