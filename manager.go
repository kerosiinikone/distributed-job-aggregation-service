package main

import (
	"log"

	"github.com/anthdm/hollywood/actor"
)

type JobRequest struct {
	Profession  string `json:"profession,omitempty"`
	JobTitle    string `json:"jobTitle,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	EmailAddr   string `json:"email,omitempty"`
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
		m.findJobService(msg)
	case *JobResult:
		log.Println("New message from Actor")
	}
} 

func (m *Manager) findJobService(job *JobRequest) {
	// Fetch the list of job listing sites
	// Spawn 1-2 worker nodes / actors on each site
	// The actors will perform the business logic / scraping and send a list of links
	// After receiving the list, the manager kills the actor
}

