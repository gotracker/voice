package envelope

import (
	"github.com/gotracker/voice"
	"github.com/gotracker/voice/loop"
)

// Envelope is an envelope for instruments
type Envelope[T any] struct {
	Enabled    bool
	Loop       loop.Loop
	Sustain    loop.Loop
	Values     []EnvPoint[T]
	OnFinished voice.Callback
}

// EnvPoint is a point for the envelope
type EnvPoint[T any] struct {
	Ticks int
	Y     T
}

func (p EnvPoint[T]) Length() int {
	return p.Ticks
}

func (p EnvPoint[T]) Value() T {
	return p.Y
}

func (p *EnvPoint[T]) Init(ticks int, value T) {
	p.Ticks = ticks
	p.Y = value
}
