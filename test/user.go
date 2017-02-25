package main

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"os"
)

func main() {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")

	c := bitbucket.NewBasicAuth(user, pass)

	res := c.User.Profile()

	fmt.Println(res) // receive the data as json format
}
