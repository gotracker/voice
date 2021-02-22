package envelope

import (
	"github.com/gotracker/voice"
	"github.com/gotracker/voice/loop"
)

// EnvPoint is a point for the envelope
type EnvPoint interface {
	Length() int
	Value(out interface{})
	Init(ticks int, value interface{})
}

// Envelope is an envelope for instruments
type Envelope struct {
	Enabled    bool
	Loop       loop.Loop
	Sustain    loop.Loop
	Values     []EnvPoint
	OnFinished voice.Callback
}
