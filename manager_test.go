package main

import (
	"testing"

	"github.com/anthdm/hollywood/actor"
)

// Test integration between actor levels, only mocking manager
// Propagation / timeout ???
func TestActors(t *testing.T) {
	resultCh := make(chan *JobResults, 1)
	input := &JobRequest{
		EmailAddr: "",
		Keywords: []string{"myynti"},
		Location: "",
		maxPages: 1,
	}
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		t.Fatal(err)
	}	
	mPid := e.Spawn(NewTestManager(resultCh), "test")
	e.Send(mPid, input)

	select {
	case msg := <- resultCh:
		if msg.Error != nil {
			t.Fatal(err)
		}
	}
}

