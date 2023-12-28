// In this exercise you'll use Go's concurrency features to parallelize a web crawler.
// Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.

// Hint: you can keep a cache of the URLs that have been fetched on a map,
// but maps alone are not safe for concurrent use!

package main

import (
	"fmt"
	"sync"
)

type (
	Fetcher interface {
		// Fetch returns the body of URL and
		// a slice of URLs found on that page.
		Fetch(url string) (body string, urls []string, err error)
	}

	// fakeFetcher is Fetcher that returns canned results.
	fakeFetcher map[string]*fakeResult

	fakeResult struct {
		body string
		urls []string
	}
)

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

type cachedURLs struct {
	mu   sync.Mutex
	urls map[string]bool
}

func (s *cachedURLs) add(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.urls[url]; ok {
		return false
	}
	s.urls[url] = true
	return true
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, visited *cachedURLs, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	if !visited.add(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	var childWg sync.WaitGroup

	for _, u := range urls {
		childWg.Add(1)
		go Crawl(u, depth-1, fetcher, visited, &childWg)
	}

	childWg.Wait()
	return
}

func main() {
	visited := &cachedURLs{urls: make(map[string]bool)}
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("https://golang.org/", 4, fetcher, visited, &wg)
	wg.Wait()
}
