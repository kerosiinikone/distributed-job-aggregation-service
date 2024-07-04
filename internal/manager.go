package actors

import (
	"github.com/anthdm/hollywood/actor"
)

// Go Set
type nodeMap map[*actor.PID]bool

type JobRequest struct{}

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
	switch ctx.Message().(type) {
	case actor.Started:
		// New Client
	case JobRequest:
		// A new Job
	}
} 

func HandleNewJobRequest() {
	// Send Message to Manager
}
