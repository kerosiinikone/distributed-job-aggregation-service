package main

import (
	"github.com/anthdm/hollywood/actor"
)

type Actor struct {
	MPID *actor.PID
}

func NewActor() actor.Receiver {
	return &Actor{}
}

func (m *Actor) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		// Perform the job
		// Send the result back to Manager 
		ctx.Engine().Poison(ctx.PID())
	}
}
