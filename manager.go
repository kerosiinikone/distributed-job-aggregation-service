package main

import (
	"log"
	"strings"

	"github.com/anthdm/hollywood/actor"
)

var (
	jobSites = []string{
		"https://www.jobly.fi/tyopaikat",
	}
)

type JobResult struct {
	Link string
}

type JobResults struct {
	Results []JobResult
	Error error
}

type finderMap map[*actor.PID]bool

type Manager struct {
	CurrentRequest *JobRequest
	FinderMap finderMap
}

func NewManager() actor.Producer {
	return func() actor.Receiver {
		return &Manager{
			FinderMap: make(finderMap),
		}
	}
}

func (m *Manager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		// Started
	case *JobRequest:
		m.CurrentRequest = msg
		m.findJobService(ctx, msg)
	case *JobResults:
		if len(msg.Results) == 0 {
			log.Println("No jobs found") // For now
		} else {
			mailer := &EmailService{}
			if err := mailer.SendEmail(msg, m.CurrentRequest); err != nil {
				log.Fatalln(err)
			}
		}
		m.CurrentRequest = nil
	}
} 

// Search different jobSites
func (m *Manager) findJobService(ctx *actor.Context, meta *JobRequest) error {
	// Spawn a worker node / actor on each site
	for _, l := range jobSites {
		if strings.Contains(l, "jobly") {
			pid := ctx.SpawnChild(NewJobActor(NewJoblyFinder(l), ctx.PID(), meta), "finder-"+l)
			m.FinderMap[pid] = true
		}
	}
	// The actors will perform the business logic / scraping and send a list of links
	// After receiving the list, the manager kills the actor
	return nil
}


