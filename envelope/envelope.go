package envelope

import (
	"github.com/gotracker/voice/loop"
)

// State is the state information about an envelope
type State[T any] struct {
	position int
	length   int
	stopped  bool
	env      *Envelope[T]
}

// Stopped returns true if the envelope state is stopped
func (e *State[T]) Stopped() bool {
	return e.stopped
}

// Stop stops the envelope state
func (e *State[T]) Stop() {
	e.stopped = true
}

// Envelope returns the envelope that the state is based on
func (e *State[T]) Envelope() *Envelope[T] {
	return e.env
}

// Reset resets the envelope
func (e *State[T]) Reset(env *Envelope[T]) {
	e.env = env
	if e.env == nil || !e.env.Enabled {
		e.stopped = true
		return
	}

	e.position = 0
	pos, _, _ := e.calcLoopedPos(true)
	if pos < len(e.env.Values) {
		e.length = e.env.Values[pos].Length()
	}
}

func (e *State[T]) calcLoopedPos(keyOn bool) (int, int, bool) {
	nPoints := len(e.env.Values)
	var looped bool
	cur, _ := loop.CalcLoopPos(e.env.Loop, e.env.Sustain, e.position, nPoints, keyOn)
	next, _ := loop.CalcLoopPos(e.env.Loop, e.env.Sustain, e.position+1, nPoints, keyOn)
	if (keyOn && e.env.Sustain.Enabled()) || e.env.Loop.Enabled() {
		looped = true
	}
	return cur, next, looped
}

// GetCurrentValue returns the current value
func (e *State[T]) GetCurrentValue(keyOn bool) (*EnvPoint[T], *EnvPoint[T], float32) {
	if e.stopped {
		return nil, nil, 0
	}

	pos, npos, looped := e.calcLoopedPos(keyOn)
	if pos >= len(e.env.Values) {
		return nil, nil, 0
	}

	if npos >= len(e.env.Values) {
		npos = pos
	}

	cur := e.env.Values[pos]
	next := e.env.Values[npos]
	t := float32(0)
	tl := cur.Length()
	if tl > 0 {
		l := float32(e.length)
		if looped {
			if e.env.Sustain.Enabled() && keyOn && e.env.Sustain.Length() == 0 {
				l = 0
			} else {
				l = float32(e.length)
			}
		}
		t = 1 - (l / float32(tl))
	}
	switch {
	case t < 0:
		t = 0
	case t > 1:
		t = 1
	}
	return &cur, &next, t
}

// Advance advances the state by 1 tick
func (e *State[T]) Advance(keyOn bool, prevKeyOn bool) bool {
	if e.stopped {
		return false
	}

	if e.env.Sustain.Enabled() && keyOn {
		if e.env.Sustain.Length() == 0 {
			return false
		}
	} else if e.env.Loop.Enabled() {
		if e.env.Loop.Length() == 0 {
			return false
		}
	}

loopAdvance:
	e.length--
	if e.length > 0 {
		return false
	}
	if keyOn != prevKeyOn && prevKeyOn {
		p, _, _ := e.calcLoopedPos(prevKeyOn)
		e.position = p
	}

	e.position++
	pos, _, _ := e.calcLoopedPos(keyOn)
	if pos >= len(e.env.Values) {
		e.stopped = true
		return true
	}

	e.length = e.env.Values[pos].Length()
	if e.length <= 0 {
		goto loopAdvance
	}
	return false
}
