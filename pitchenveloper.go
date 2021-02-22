package voice

import (
	"github.com/gotracker/voice/period"
)

// PitchEnveloper is a pitch envelope interface
type PitchEnveloper interface {
	EnablePitchEnvelope(enabled bool)
	IsPitchEnvelopeEnabled() bool
	GetCurrentPitchEnvelope() period.Delta
	SetPitchEnvelopePosition(pos int)
}
