package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type JoblyFinder struct {
	Link string
}

func NewJoblyFinder(link string) Finder {
	return &JoblyFinder{
		Link: link,
	}
}

func (fi *JoblyFinder) getLink() string {
	return fi.Link
}

// Needs to be able to access the next page of jobs on a listing site
// Initial checks based on listing title and description 
// Refactor the logic later
func (fi *JoblyFinder) extractJobListings(link string, out chan<- JobPosting, cancel context.CancelFunc) error {
	var (
		f func(*html.Node, *JobPosting)
	)

	res, err := http.Get(link)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return err
	}

	f = func(n *html.Node, jp *JobPosting) {
		if n.Type == html.ElementNode && n.Data == "article" {
			jp = &JobPosting{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, jp)
			}
			if len(jp.Text) != 0 && len(jp.Link) != 0 {
				// Pipe to a channel that takes care of posting
				out <- JobPosting{Text: jp.Text, Link: jp.Link}
			}
		} else if n.Type == html.TextNode && jp != nil {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				jp.Text = append(jp.Text, text)
			}
		} else {
			if n.Data == "a" && jp != nil {
				for _, a := range n.Attr {
					// Don't add duplicate links !!!
					// Assumes that a post "article" has only a single valid job link
					if a.Key == "href" && jobListingLinkMatcher(a.Val) {
						jp.Link = a.Val
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, jp)
			}
		}
	}
	f(doc, nil)

	defer cancel()

	return nil
}