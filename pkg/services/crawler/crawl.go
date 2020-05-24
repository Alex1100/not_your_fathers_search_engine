package crawler

import (
	"fmt"
	"golang.org/x/net/html"

	"log"
	"net/http"
	// "os"
	"strings"
)

// WORK IN PROGRESS


// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
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
		// add go-routine and check to upsert/insert into
		// memory and cockroachdb here
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
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := Extract(url)
	<-tokens // release the tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

func StartCrawlProcess(srcURL string) []string {
	worklist := make(chan []string)  // lists of URL's, may have dups
	unseenLinks := make(chan string) // deduped URL's
	// Add command-line args to worklist
	go func() {
		// change from getting from terminal to an api/client side request perhaps
		// worklist <- os.Args[1:]
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
	linkCollection := make([]string, 0)
	for list := range worklist {
		for _, link := range list {
			if len(seen) >= 5 {
				return linkCollection
			}

			if !seen[link] && !strings.Contains(link, "localhost") {
				seen[link] = true
				linkCollection = append(linkCollection, link)
				unseenLinks <- link
			}
		}
	}

	return linkCollection
}