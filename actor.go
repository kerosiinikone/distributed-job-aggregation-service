package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
)

type Actor struct {
	MPID 	*actor.PID
	Link 	string
	Meta 	*JobRequest
}

type JobPosting struct {
	Text      []string
	Link      []string
}

func NewActor(link string, mpid *actor.PID, meta *JobRequest) actor.Producer {
	return func () actor.Receiver  {
		return &Actor{
			MPID: mpid,
			Link: link,
			Meta: meta,
		}
	}
	
}

func (a *Actor) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		// Perform the job
		// Send the result back to Manager 
		postings, err := a.scrapeJobService()
		if err != nil {
			log.Println(err)
		} 
		log.Println(postings, "postings")
		ctx.Engine().Poison(ctx.PID())
	}
}

// As long as there are related job postings, spawn new actors to dig deeper ??
// Spawn new actors to check other pages ? -> later
func (a *Actor) scrapeJobService() ([]JobPosting, error) {
	var allPostings []JobPosting 
	for i := 0; i < 10; i++ {
		postings, err := a.extractJobListings(a.Link + "?page=" + fmt.Sprintf("%d", i))
		if err != nil {
			return []JobPosting{}, err
		}
		allPostings = append(allPostings, postings...)
	}
	
	return allPostings, nil
}

func jobListingLinkMatcher(val string) bool {
	return (strings.Contains(val, "/tyopaikka/") || strings.Contains(val, "/tyopaikat/tyo/") || strings.Contains(val, "/avoimet-tyopaikat/"))
}

// Needs access to Meta object -> method instead of function
func (a *Actor) jobListingKeywordMatcher(val string) bool {
	var valid bool

	if len(a.Meta.Keywords) == 0 {
		return true
	}
	for _, k := range a.Meta.Keywords {
		if strings.Contains(strings.ToLower(val), strings.ToLower(k)) {
			valid = true
		}
	}
	return valid
}

// Needs to be able to access the next page of jobs on a listing site
// Initial checks based on listing title and description 
// Refactor the logic later
func (a *Actor) extractJobListings(link string) ([]JobPosting, error) {
	var (
		f func(*html.Node, *JobPosting)
		postings []JobPosting
	)

	res, err := http.Get(link)
	if err != nil {
		return []JobPosting{}, err
	}
	
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	f = func(n *html.Node, jp *JobPosting) {
		if n.Type == html.ElementNode && n.Data == "article" {
			jp = &JobPosting{} 
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, jp)
			}
			if len(jp.Text) != 0 && len(jp.Link) != 0 {
				postings = append(postings, JobPosting{Text: jp.Text, Link: jp.Link})
			}
		} else if n.Type == html.TextNode && jp != nil {
			text := strings.TrimSpace(n.Data)
			if text != "" && a.jobListingKeywordMatcher(text) {
				jp.Text = append(jp.Text, text)
			}
		} else {
			if n.Data == "a" && jp != nil {
				for _, a := range n.Attr {
					if a.Key == "href" && jobListingLinkMatcher(a.Val) {
						jp.Link = append(jp.Link, a.Val)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, jp)
			}
		}
	}
	f(doc, nil)
	
	return postings, nil
}