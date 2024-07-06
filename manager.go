package main

import (
	"github.com/anthdm/hollywood/actor"
)

var jobSites = []string{
	"https://www.jobly.fi/tyopaikat",
}

type JobResult struct {
	Link string
}

type JobResults []JobResult

type finderMap map[*actor.PID]bool

type Manager struct {
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
		m.findJobService(ctx, msg)
	case *JobResults:
		// -> Continue
		// Email service, etc
	}
} 

func (m *Manager) findJobService(ctx *actor.Context, meta *JobRequest) error {
	// Spawn a worker node / actor on each site
	for _, l := range jobSites {
		pid := ctx.SpawnChild(NewFinder(l, ctx.PID(), meta), "finder-"+l)
		m.FinderMap[pid] = true
	}
	// The actors will perform the business logic / scraping and send a list of links
	// After receiving the list, the manager kills the actor
	return nil
}

