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
	files := get_file_list()

	for _, file := range files {
		html := get_html_by_file(file)
		logs := get_logs_from_html(html)

		for _, log := range logs {
			paifu := get_paifu(log)
			fmt.Println(paifu)
			reaches := get_reaches_from_paifu(paifu)
			fmt.Println(reaches)
			break
		}
		fmt.Println(logs)
		break
	}
}

func get_reaches_from_paifu(paifu string) int {
	r := regexp.MustCompile(`REACH[^>]*?step="2"/>`)
	matches := r.FindAllString(paifu, -1)
	return len(matches)
}

func get_paifu(id string) string {
	url := "https://tenhou.net/0/log/?" + id
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func get_logs_from_html(html *goquery.Document) []string {
	logs := []string{}

	html.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			r := regexp.MustCompile(`log=([\w\-]+)`)
			matches := r.FindStringSubmatch(href)
			if len(matches) > 0 {
				logs = append(logs, matches[1])
			}
		}
	})

	return logs
}

func get_html_by_file(file string) *goquery.Document {
	url := "https://tenhou.net/sc/raw/dat/" + file
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

	return doc
}

func get_file_list() []string {
	res, err := http.Get("https://tenhou.net/sc/raw/list.cgi")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`\w+\.html\.gz`)
	files := []string{}
	for _, matches := range r.FindAllStringSubmatch(string(body), -1) {
		files = append(files, matches[0])
	}

	return files
}
