package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type thread struct {
	id    string `db:"id"`
	title string `db:"title"`
}

var baseURL = "http://boards.4chan.org/biz"

func getDocument(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return html.Parse(strings.NewReader(string(response)))
}

func getAttr(n *html.Node, attr string) string {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}

func getThreadsForPage(url string, threads []string) []string {
	var err error

	var f func(*html.Node)
	f = func(n *html.Node) {

		// grab  thread links
		if n.Type == html.ElementNode && n.Data == "a" {
			class := getAttr(n, "class")
			if class == "replylink" {
				href := baseURL + "/" + getAttr(n, "href")
				threads = append(threads, href)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	doc, err := getDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	f(doc)

	return threads
}

func getThreads() []string {
	threads := getThreadsForPage(baseURL, make([]string, 0))

	for i := 2; i <= 10; i++ {
		page := strconv.Itoa(i)
		threads = append(threads, getThreadsForPage(baseURL+"/"+page, make([]string, 0))...)
	}

	return threads
}

func getThreadTitle(threadURL string) string {
	doc, err := getDocument(threadURL)
	if err != nil {
		log.Fatal(err)
	}

	title := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" && getAttr(n, "class") == "subject" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return title
}

func getThreadID(threadURL string) string {
	splitThread := strings.Split(threadURL, "/")

	if splitThread[len(splitThread)-2] == "thread" {
		return splitThread[len(splitThread)-1]
	}

	return splitThread[len(splitThread)-2]
}

func main() {
	ts := getThreads()
	threads := make(map[string]bool)
	for _, t := range ts {
		tid := getThreadID(t)
		title := getThreadTitle(t)
		td := thread{tid, title}
		if threads[tid] == false {
			fmt.Println(td)
			threads[tid] = true
		}
	}
}
