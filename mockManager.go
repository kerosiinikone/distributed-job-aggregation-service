package main

import (
	"strings"

	"github.com/anthdm/hollywood/actor"
)

type TestManager struct {
	CurrentRequest *JobRequest
	FinderMap       finderMap
	resultCh 		chan *JobResults
}

func NewTestManager(ch chan *JobResults) actor.Producer {
	return func() actor.Receiver {
		return &TestManager{
			FinderMap: make(finderMap),
			resultCh: ch,
		}
	}
}

func (m *TestManager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *JobRequest:
		m.CurrentRequest = msg
		m.findJobService(ctx, msg)
	case *JobResults:
		m.resultCh <- msg
	}
}

func (m *TestManager) findJobService(ctx *actor.Context, meta *JobRequest) error {
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