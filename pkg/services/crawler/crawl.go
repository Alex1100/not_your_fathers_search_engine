package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func isValidUrl(src string) bool {
	_, err := url.ParseRequestURI(src)
	if err != nil {
		return false
	}

	u, err := url.Parse(src)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	if !isValidUrl(url) {
		return nil, fmt.Errorf("url is not a valid protocol: %s", url)
	}
	fmt.Println("Calkung it: ", url)
	networkClient := http.Client{
		Timeout: 2000 * time.Millisecond,
	}

	resp, err := networkClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// example of a semaphore
// if the size of the buffered
// channel was going to be 1
// then the solo thread could
// also be secured with a mutex
// aka import "sync"
// declare var (mu sync.Mutex)
// and use mu.Lock()
// before a given operation
// read/write
// and mu.Unlock()

var tokens = make(chan struct{}, 20000)

func crawl(url string) []string {
	tokens <- struct{}{} // acquire a token
	list, err := Extract(url)
	<-tokens // release the tokens

	if err != nil {
		log.Print(err)
	}

	return list
}

func StartCrawlProcess(srcURL string) []byte {
	worklist := make(chan []string)  // lists of URL's, may have dups
	unseenLinks := make(chan string) // deduped URL's

	go func() {
		worklist <- []string{srcURL}
	}()

	// Create 20000 crawler goroutines to fetch each unseen link
	for i := 0; i < 20000; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() {
					worklist <- foundLinks
				}()
			}
		}()
	}

	// The main goroutine dedups worklist items
	// and sends the unseen ones to the crawlers
	seen := make(map[string]bool)
	linkCollection := make([]byte, 0)

	for list := range worklist {
		for _, link := range list {
			if len(seen) >= 1000 {
				return linkCollection
			}

			if !seen[string(link)] && !strings.Contains(string(link), "localhost") {
				seen[string(link)] = true
				linkCollection = append(linkCollection, []byte(string(link)+"\n\n")...)
				unseenLinks <- string(link)
			}
		}
	}

	return linkCollection
}
