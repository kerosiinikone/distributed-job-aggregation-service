package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type visitorMap map[*actor.PID]bool

type Finder interface {
	extractJobListings(string, chan<- JobPosting, context.CancelFunc) error
	scrapeJobService(chan JobPosting, context.CancelFunc)
}

type jobActor struct {
    finder 			Finder 
	Meta 			*JobRequest
	VisitorMap 		visitorMap
	timer 			*time.Timer
	MPID 			*actor.PID
	Results			JobResults
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

func NewJobActor(f Finder, mpid *actor.PID, meta *JobRequest) actor.Producer {
	return func() actor.Receiver {
		return &jobActor{
			finder: f,
			MPID: mpid,
			Meta: meta,
			VisitorMap: make(visitorMap),
		}
	}
}

func (a *jobActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *FilteredJobPost:
		// Check whether any Visitors are still alive
		// After that send the results to the Manager
		// Kill the process
		a.Results = append(a.Results, JobResult{
			Link: msg.Link,
		})
	case *KillVisitor:
		// Job done -> all spawned visitors have died
		if a.timer != nil {
            a.timer.Stop()
        }
		// Make the debounce time a global variable or part of config
        a.timer = time.AfterFunc(10*time.Second, func() {
            results := a.Results
			ctx.Send(a.MPID, &results)
			ctx.Engine().Poison(ctx.PID())
        })
	case actor.Started:
		// Perform the job
		// Send the result back to Manager

		// Pull from config (time)
		pCtx, cancel := context.WithTimeout(context.Background(), time.Second * 60)
		jobChan := make(chan JobPosting) 

		go a.finder.scrapeJobService(jobChan, cancel)
		go a.handleJobSites(ctx, jobChan, pCtx)
	}
}

// Spawn Visitors (with link array on each)
// Is concurrent and receives a single posting once it has been scraped through a channel
func (a *jobActor) handleJobSites(ctx *actor.Context, in <-chan JobPosting, pCtx context.Context) {
	var idx int
	for {
		select {
		case job := <-in:
			pid := ctx.SpawnChild(NewVisitor(job.Link, ctx.PID(), a.Meta), fmt.Sprintf("visitor-%d", idx))
			a.VisitorMap[pid] = true
			idx++
		case <-pCtx.Done():
            return 
        }
	}
}

func jobListingLinkMatcher(val string) bool {
	return (strings.Contains(val, "/tyopaikka/") || strings.Contains(val, "/tyopaikat/tyo/") || strings.Contains(val, "/avoimet-tyopaikat/"))
}