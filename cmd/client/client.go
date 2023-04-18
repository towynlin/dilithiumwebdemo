package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/towynlin/dilithiumwebdemo/eddilithium2jwt"
)

func main() {
	url := os.Getenv("SERVER_URL")
	if url == "" || url[len(url)-1] == '/' {
		panic("Please set the env var SERVER_URL with no trailing slash")
	}
	url += "/jobs"

	if len(os.Args) < 2 {
		panic("Pass at least one arg: get, post, delete")
	}

	cmd := os.Args[1]
	if cmd != "get" && cmd != "post" && cmd != "delete" {
		panic("Invalid command: " + cmd)
	}

	if cmd == "delete" {
		if len(os.Args) < 3 {
			panic("Pass the job you want to delete: delete 7854460d-bab8-412e-8f35-a0659ae65da6")
		}
		url += "/" + os.Args[2]
	}

	signedString, err := eddilithium2jwt.NewSignedString(
		os.Getenv("SIGNING_ID"), os.Getenv("SIGNING_KEY"))
	if err != nil {
		panic(err)
	}

	method := strings.ToUpper(cmd)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+signedString)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resBody))
}
