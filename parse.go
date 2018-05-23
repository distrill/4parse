package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

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

func classIs(n *html.Node, value string) bool {
	return getAttr(n, "class") == value
}

func titleIs(n *html.Node, value string) bool {
	return getAttr(n, "title") == value
}

func nodeIs(n *html.Node, value string) bool {
	return n.Type == html.ElementNode && n.Data == value
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

/*
 * thread stuffs
 */
func getThreadTitle(threadURL string) string {
	doc, err := getDocument(threadURL)
	if err != nil {
		log.Fatal(err)
	}

	title := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if nodeIs(n, "span") && classIs(n, "subject") && n.FirstChild != nil {
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

func getThreadInfo(threadURL string) thread {
	id := getThreadID(threadURL)
	title := getThreadTitle(threadURL)
	return thread{id, title}
}

/*
 * comment stuffs
 */
func getCommentID(n *html.Node) string {
	return strings.Replace(getAttr(n, "id"), "m", "", 1)
}

func getCommentImageURL(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nodeIs(c, "div") && classIs(c, "file") {
			return strings.Replace(getAttr(c.FirstChild.FirstChild.NextSibling, "href"), "//", "", 1)
		}
	}
	return ""
}

func getPosterID(n *html.Node) string {
	posterID := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if nodeIs(n, "span") && classIs(n, "hand") && titleIs(n, "Highlight posts by this ID") {
			posterID = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)

	return posterID
}

func getCommentsInfo(threadURL string, threadID string) []comment {
	isOP := true
	comments := make([]comment, 0)

	doc, err := getDocument(threadURL)
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if nodeIs(n, "blockquote") {
			cmt := comment{
				getCommentID(n),
				getCommentImageURL(n.Parent),
				getPosterID(n.Parent),
				isOP,
				threadID,
			}
			comments = append(comments, cmt)
			isOP = false
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return comments
}

/*
 * comment data stuffs (greentexts and text content)
 */
func getGreentexts(n *html.Node, commentID string) []commentData {
	greentexts := make([]commentData, 0)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "span" {
			greentexts = append(greentexts, commentData{
				DataType:  "GREENTEXT",
				Contents:  c.FirstChild.Data,
				CommentID: commentID,
			})
		}
	}
	return greentexts
}

func getTextContents(n *html.Node, commentID string) []commentData {
	textContents := make([]commentData, 0)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data != "a" && c.Data != "span" && c.Data != "br" {
			textContents = append(textContents, commentData{
				DataType:  "TEXTCONTENT",
				Contents:  c.Data,
				CommentID: commentID,
			})
		}
	}

	return textContents

}

// returns [greentexts..., textcontents...]
func getCommentData(threadURL string) []commentData {
	greenTexts := make([]commentData, 0)
	textContents := make([]commentData, 0)

	doc, err := getDocument(threadURL)
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "blockquote" {
			commentID := getCommentID(n)
			greenTexts = append(greenTexts, getGreentexts(n, commentID)...)
			textContents = append(textContents, getTextContents(n, commentID)...)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return append(greenTexts, textContents...)
}

func isValidCommentReference(href string) bool {
	return !strings.Contains(href, "/thread/") && strings.Contains(href, "#p")
}

/*
 * comment reference stuffs
 */
func getReferences(n *html.Node, parentCommentID string) []commentReference {
	references := make([]commentReference, 0)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "a" && isValidCommentReference(getAttr(c, "href")) {
			references = append(references, commentReference{
				ParentCommentID: parentCommentID,
				ChildCommentID:  strings.Replace(getAttr(c, "href"), "#p", "", -1),
			})
		}
	}

	return references
}

func getCommentReferences(threadURL string) []commentReference {
	references := make([]commentReference, 0)

	doc, err := getDocument(threadURL)
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "blockquote" {
			commentID := getCommentID(n)
			references = append(references, getReferences(n, commentID)...)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return references
}
