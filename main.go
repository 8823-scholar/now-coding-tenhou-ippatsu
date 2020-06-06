package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	names := get_log_list()
	fmt.Println(names)
}

func get_log_list() []string {
	res, err := http.Get("https://tenhou.net/sc/raw/list.cgi")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`\w+\.html\.gz`)
	names := []string{}
	for _, matches := range r.FindAllStringSubmatch(string(body), -1) {
		names = append(names, matches[0])
	}

	return names
}
