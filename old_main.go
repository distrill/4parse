package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"golang.org/x/net/html"
// )

// type thread struct {
// 	id    string `db:"id"`
// 	title string `db:"title"`
// }

// type comment struct {
// 	id       string `db:"id"`
// 	imageURL string `db:"image_url"`
// 	posterID string `db:"poster_id"`
// 	isOP     bool   `db:"is_op"`
// 	threadID string `db:"thread_id"`
// }

// type commentData struct {
// 	id        string `db:"id"`
// 	dataType  string `db:"data_type"`
// 	contents  string `db:"contents"`
// 	commentID string `db:"comment_id"`
// }

// type commentReference struct {
// 	id              string `db:"id"`
// 	parentCommentID string `db:"parent_comment_id"`
// 	childCommentID  string `db:"child_comment_id"`
// }

// func getAttr(n *html.Node, attr string) string {
// 	for _, a := range n.Attr {
// 		if a.Key == attr {
// 			return a.Val
// 		}
// 	}
// 	return ""
// }

// func getDocument(url string) (*html.Node, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	response, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return html.Parse(strings.NewReader(string(response)))
// }

// func getThreadsForPage(url string, threads []string) []string {
// 	var err error

// 	var f func(*html.Node)
// 	f = func(n *html.Node) {

// 		// grab  thread links
// 		if n.Type == html.ElementNode && n.Data == "a" {
// 			class := getAttr(n, "class")
// 			if class == "replylink" {
// 				href := "http://boards.4chan.org/biz/" + getAttr(n, "href")
// 				threads = append(threads, href)
// 			}
// 		}

// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			f(c)
// 		}
// 	}

// 	doc, err := getDocument(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	f(doc)

// 	return threads
// }

// func getThreads() []string {
// 	threads := getThreadsForPage("http://boards.4chan.org/biz", make([]string, 0))

// 	for i := 2; i <= 10; i++ {
// 		page := strconv.Itoa(i)
// 		threads = append(threads, getThreadsForPage("http://boards.4chan.org/biz/"+page, make([]string, 0))...)
// 	}

// 	return threads
// }

// func getID(n *html.Node) string {
// 	return strings.Replace(getAttr(n, "id"), "m", "", 1)
// }

// func getGreentexts(n *html.Node) []string {
// 	greentexts := make([]string, 0)

// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Data == "span" {
// 			greentexts = append(greentexts, c.FirstChild.Data)
// 		}
// 	}
// 	return greentexts
// }

// func getReferences(n *html.Node) []string {
// 	references := make([]string, 0)

// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Data == "a" {
// 			references = append(references, strings.Replace(c.FirstChild.Data, ">>", "", -1))
// 		}
// 	}

// 	return references
// }

// func getContents(n *html.Node) []string {
// 	contents := make([]string, 0)

// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Data != "a" && c.Data != "span" && c.Data != "br" {
// 			contents = append(contents, c.Data)
// 		}
// 	}

// 	return contents

// }

// func doWhatever(threadURL string) {
// 	doc, err := getDocument(threadURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var f func(*html.Node)
// 	f = func(n *html.Node) {

// 		// grab  thread links
// 		if n.Type == html.ElementNode && n.Data == "blockquote" {
// 			// com := comment{}
// 			// fmt.Println(getID(n))
// 			// fmt.Println(getReferences(n))
// 			// fmt.Println(getGreentexts(n))
// 			// fmt.Println(getContents(n))
// 			// fmt.Printf("\n\n")
// 		}

// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			f(c)
// 		}
// 	}

// 	f(doc)
// }

// func printNode(s string, n *html.Node) {
// 	fmt.Printf("%v\n", s)
// 	fmt.Printf("type: %v\n", n.Type)
// 	fmt.Printf("data: %v\n", n.Data)
// 	fmt.Printf("namespace: %v\n", n.Namespace)
// 	fmt.Printf("attr: %v\n", n.Attr)
// 	fmt.Print("\n")
// }

// func getThreadTitle(threadURL string) string {
// 	doc, err := getDocument(threadURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	title := ""

// 	var f func(*html.Node)
// 	f = func(n *html.Node) {
// 		if n.Type == html.ElementNode && n.Data == "span" && getAttr(n, "class") == "subject" && n.FirstChild != nil {
// 			title = n.FirstChild.Data
// 		}
// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			f(c)
// 		}
// 	}
// 	f(doc)

// 	return title
// }

// func getThreadID(threadURL string) string {
// 	splitThread := strings.Split(threadURL, "/")

// 	if splitThread[len(splitThread)-2] == "thread" {
// 		return splitThread[len(splitThread)-1]
// 	}

// 	return splitThread[len(splitThread)-2]
// }

// func oldMain() {
// 	ts := getThreads()
// 	threads := make(map[string]bool)
// 	for _, t := range ts {
// 		tid := getThreadID(t)
// 		title := getThreadTitle(t)
// 		td := thread{tid, title}
// 		if threads[tid] == false {
// 			fmt.Println(td)
// 			threads[tid] = true
// 		}
// 	}

// 	// thread := "http://boards.4chan.org/biz/thread/9493753/where-are-those-fools-who-once-proudly-posted"
// 	// tid := getThreadID(thread)
// 	// title := getThreadTitle(thread)
// 	// fmt.Printf("id: %v\ntitle: %v\n", tid, title)
// }
