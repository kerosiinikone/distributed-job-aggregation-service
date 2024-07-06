package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
)

// Why not stream every found job listing to the handleJobSites() and spawn new Visitors as new JobPostings arrive ???

type KillVisitor struct {
	PID *actor.PID
}

type FilteredJobPost struct {
	Link string
}

type visitorMap map[*actor.PID]bool

type Finder struct {
	MPID 	*actor.PID
	Link 	string
	Results []JobResult // Can be a struct also
	Meta 	*JobRequest
	VisitorMap visitorMap
}

type JobPosting struct {
	Text      []string
	Link      string
}

func NewFinder(link string, mpid *actor.PID, meta *JobRequest) actor.Producer {
	return func () actor.Receiver  {
		return &Finder{
			MPID: mpid,
			Link: link,
			Meta: meta,
			VisitorMap: make(visitorMap),
		}
	}
}

func (fi *Finder) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *FilteredJobPost:
		// Check whether any Visitors are still alive
		// After that send the results to the Manager
		// Kill the process
		fi.handleResults(ctx, msg)
	case *KillVisitor:
		// After the Visitor has been poisoned, 
		// remove from the VisitorMap
		delete(fi.VisitorMap, msg.PID)
	case actor.Started:
		// Perform the job
		// Send the result back to Manager 
		postings, err := fi.scrapeJobService()
		if err != nil {
			log.Fatalln(err)
		} 
		err = fi.handleJobSites(ctx, postings)
		if err != nil {
			log.Fatalln(err)
		} 
	}
}

func (fi *Finder) handleResults(ctx *actor.Context, job *FilteredJobPost) {
	var results JobResults
	
	fi.Results = append(fi.Results, JobResult{
		Link: job.Link,
	})
	if len(fi.VisitorMap) > 0 {
		return
	}
	// Any cleanup necessary ???
	results = fi.Results
	ctx.Send(fi.MPID, &results)
	ctx.Engine().Poison(ctx.PID())
}

// As long as there are related job postings, spawn new actors to dig deeper ??
// Spawn new actors to check other pages ? -> later
func (fi *Finder) scrapeJobService() ([]JobPosting, error) {
	var allPostings []JobPosting 
	for i := 0; i < 10; i++ {
		postings, err := fi.extractJobListings(fi.Link + "?page=" + fmt.Sprintf("%d", i))
		if err != nil {
			return []JobPosting{}, err
		}
		allPostings = append(allPostings, postings...)
	}
	
	return allPostings, nil
}

// Spawn Visitors (with link array on each)
func (fi *Finder) handleJobSites(ctx *actor.Context, postings []JobPosting) error {
	for i, p := range postings {
		pid := ctx.SpawnChild(NewVisitor(p.Link, ctx.PID(), fi.Meta), fmt.Sprintf("visitor-%d", i))
		fi.VisitorMap[pid] = true
	}
	return nil
}

func jobListingLinkMatcher(val string) bool {
	return (strings.Contains(val, "/tyopaikka/") || strings.Contains(val, "/tyopaikat/tyo/") || strings.Contains(val, "/avoimet-tyopaikat/"))
}

// Needs to be able to access the next page of jobs on a listing site
// Initial checks based on listing title and description 
// Refactor the logic later
func (fi *Finder) extractJobListings(link string) ([]JobPosting, error) {
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
			if text != "" {
				jp.Text = append(jp.Text, text)
			}
		} else {
			if n.Data == "a" && jp != nil {
				for _, a := range n.Attr { 
					// Don't add duplicate links !!!
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
	
	return postings, nil
}