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
	MPID *actor.PID
	Link string
	Meta *JobRequest
}

type JobPosting struct {
	Text      []string
	Link      string
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
func (a *Actor) scrapeJobService() ([]JobPosting, error) {
	postings, err := a.extractJobListings(a.Link)
	if err != nil {
		return []JobPosting{}, err
	}
	return postings, nil
}

func jobListingLinkMatcher(val string) bool {
	return (strings.Contains(val, "/tyopaikka/") || strings.Contains(val, "/tyopaikat/tyo/") || strings.Contains(val, "/avoimet-tyopaikat/"))
}

func jobListingKeywordMatcher(val string) bool {
	return len(val) > 0
}

// Needs to be able to access the next page of jobs on a listing site
// Initial checks based on listing title and description 
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
			postings = append(postings, JobPosting{Text: jp.Text})
		} else if n.Type == html.TextNode && jp != nil {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				jp.Text = append(jp.Text, text)
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, jp)
			}
		}
	}
	f(doc, nil) 
	
	return postings, nil
}