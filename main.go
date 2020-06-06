package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	names := get_log_list()

	for _, name := range names {
		html := get_html_by_log_name(name)
		fmt.Println(html)
		break
	}
}

func get_html_by_log_name(name string) string {
	url := "https://tenhou.net/sc/raw/dat/" + name
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	html, err := doc.Html()
	if err != nil {
		log.Fatal(err)
	}

	return html
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
