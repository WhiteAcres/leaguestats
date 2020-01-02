package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func request() ([]byte, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func mainOld() {
	body, err := request()
	if err != nil {
		panic(err)
	}
	s := string(body[:])
	fmt.Printf(s)
}
