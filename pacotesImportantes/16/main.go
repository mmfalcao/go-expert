package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

const url = "http://google.com"
const APPJSON = "application/json"

func main() {

	c := http.Client{Timeout: time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", APPJSON)
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))

}
