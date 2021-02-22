package component

import (
	"github.com/gotracker/voice"
	"github.com/gotracker/voice/envelope"
	"github.com/gotracker/voice/period"
)

// PitchEnvelope is an frequency modulation envelope
type PitchEnvelope struct {
	enabled   bool
	state     envelope.State
	delta     period.Delta
	keyOn     bool
	prevKeyOn bool
}

// Reset resets the state to defaults based on the envelope provided
func (e *PitchEnvelope) Reset(env *envelope.Envelope) {
	e.state.Reset(env)
	e.keyOn = false
	e.prevKeyOn = false
	e.update()
}

// SetEnabled sets the enabled flag for the envelope
func (e *PitchEnvelope) SetEnabled(enabled bool) {
	e.enabled = enabled
}

// IsEnabled returns the enabled flag for the envelope
func (e *PitchEnvelope) IsEnabled() bool {
	return e.enabled
}

// GetCurrentValue returns the current cached envelope value
func (e *PitchEnvelope) GetCurrentValue() period.Delta {
	return e.delta
}

// SetEnvelopePosition sets the current position in the envelope
func (e *PitchEnvelope) SetEnvelopePosition(pos int) voice.Callback {
	keyOn := e.keyOn
	prevKeyOn := e.prevKeyOn
	env := e.state.Envelope()
	e.state.Reset(env)
	// TODO: this is gross, but currently the most optimal way to find the correct position
	for i := 0; i < pos; i++ {
		if doneCB := e.Advance(keyOn, prevKeyOn); doneCB != nil {
			return doneCB
		}
	}
	return nil
}

// Advance advances the envelope state 1 tick and calculates the current envelope value
func (e *PitchEnvelope) Advance(keyOn bool, prevKeyOn bool) voice.Callback {
	e.keyOn = keyOn
	e.prevKeyOn = prevKeyOn
	var doneCB voice.Callback
	if done := e.state.Advance(e.keyOn, e.prevKeyOn); done {
		doneCB = e.state.Envelope().OnFinished
	}
	e.update()
	return doneCB
}

func (e *PitchEnvelope) update() {
	cur, next, t := e.state.GetCurrentValue(e.keyOn)

	y0 := float32(0)
	if cur != nil {
		cur.Value(&y0)
	}

	y1 := float32(0)
	if next != nil {
		next.Value(&y1)
	}

	val := y0 + t*(y1-y0)
	e.delta = period.Delta(-val)
}