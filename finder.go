package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
)

type visitorMap map[*actor.PID]bool

type Finder struct {
	timer 			*time.Timer
	MPID 			*actor.PID
	Link 			string
	Results			JobResults
	Meta 			*JobRequest
	VisitorMap 		visitorMap
}

type KillVisitor struct {
	PID *actor.PID
}

type FilteredJobPost struct {
	Link string
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
		fi.handleResults(msg)
	case *KillVisitor:
		// Job done -> all spawned visitors have died
		if fi.timer != nil {
            fi.timer.Stop()
        }
		// Make the debounce time a global variable or part of config
        fi.timer = time.AfterFunc(10*time.Second, func() {
            results := fi.Results
			ctx.Send(fi.MPID, &results)
			ctx.Engine().Poison(ctx.PID())
        })
	case actor.Started:
		// Perform the job
		// Send the result back to Manager
		jobChan := make(chan JobPosting)
		doneCh := make(chan struct{})  

		go fi.scrapeJobService(jobChan, doneCh)
		go fi.handleJobSites(ctx, jobChan, doneCh)
	}
}

func (fi *Finder) handleResults(job *FilteredJobPost) {
	fi.Results = append(fi.Results, JobResult{
		Link: job.Link,
	})
}

// As long as there are related job postings, spawn new actors to dig deeper ??
func (fi *Finder) scrapeJobService(jobCh chan JobPosting, doneCh chan struct{}) {
	for i := 0; i < 10; i++ {
		go func() {
			if err := fi.extractJobListings(fi.Link + "?page=" + fmt.Sprintf("%d", i), jobCh, doneCh); err != nil {
				log.Fatalln(err)
			}
		}()
	}
}

// Spawn Visitors (with link array on each)
// Is concurrent and receives a single posting once it has been scraped through a channel
func (fi *Finder) handleJobSites(ctx *actor.Context, in <-chan JobPosting, doneCh <-chan struct{}) {
	var idx int
	for {
		select {
		case job := <-in:
			pid := ctx.SpawnChild(NewVisitor(job.Link, ctx.PID(), fi.Meta), fmt.Sprintf("visitor-%d", idx))
			fi.VisitorMap[pid] = true
			idx++
		case <-doneCh:
            return 
        }
	}
}

func jobListingLinkMatcher(val string) bool {
	return (strings.Contains(val, "/tyopaikka/") || strings.Contains(val, "/tyopaikat/tyo/") || strings.Contains(val, "/avoimet-tyopaikat/"))
}

// Needs to be able to access the next page of jobs on a listing site
// Initial checks based on listing title and description 
// Refactor the logic later
func (fi *Finder) extractJobListings(link string, out chan<- JobPosting, doneCh chan<- struct{}) error {
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

	close(doneCh)
	
	return nil
}