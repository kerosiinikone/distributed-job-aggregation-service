package main

import (
	"log"

	"github.com/anthdm/hollywood/actor"
)

var jobSites = []string{
	"https://tyopaikat.oikotie.fi/tyopaikat",
	"https://www.jobly.fi/tyopaikat",
	"https://duunitori.fi/tyopaikat",
}

type JobResult struct {
	Links string
	// ...
}

type nodeMap map[*actor.PID]bool

type Manager struct {
	NodeMap nodeMap
}

func NewManager() actor.Producer {
	return func() actor.Receiver {
		return &Manager{
			NodeMap: make(nodeMap),
		}
	}
}

func (m *Manager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		// Started
	case *JobRequest:
		log.Println("New message to Manager")
		m.findJobService(ctx, msg)
	case *JobResult:
		log.Println("New message from Actor")
	}
} 

// Link fetching can be done once the sercver is up, not on every request !!!

func (m *Manager) findJobService(ctx *actor.Context, meta *JobRequest) error {
	// Spawn a worker node / actor on each site
	for _, l := range jobSites {
		pid := ctx.SpawnChild(NewActor(l, ctx.PID(), meta), "actor-"+l)
		m.NodeMap[pid] = true
	}

	// The actors will perform the business logic / scraping and send a list of links
	// After receiving the list, the manager kills the actor
	return nil
}

