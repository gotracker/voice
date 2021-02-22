package optional

import (
	"github.com/gotracker/gomixing/panning"
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/voice/period"
)

// Value is an optional value
type Value struct {
	set   bool
	value interface{}
}

// Reset clears the memory on the value
func (o *Value) Reset() {
	o.value = nil
	o.set = false
}

// Set updates the value and sets the set flag
func (o *Value) Set(value interface{}) {
	o.value = value
	o.set = true
}

func (o *Value) IsSet() bool {
	return o.set
}

// Get returns the value and its set flag
func (o *Value) Get() (interface{}, bool) {
	return o.value, o.set
}

// GetBool returns the stored value as a boolean and if it has been set
func (o *Value) GetBool() (bool, bool) {
	if v, ok := o.value.(bool); ok {
		return v, o.set
	}
	return false, false
}

// GetInt returns the stored value as an integer and if it has been set
func (o *Value) GetInt() (int, bool) {
	if v, ok := o.value.(int); ok {
		return v, o.set
	}
	return 0, false
}

// GetVolume returns the stored value as a volume and if it has been set
func (o *Value) GetVolume() (volume.Volume, bool) {
	if v, ok := o.value.(volume.Volume); ok {
		return v, o.set
	}
	return volume.Volume(1), false
}

// GetPeriod returns the stored value as a period and if it has been set
func (o *Value) GetPeriod() (period.Period, bool) {
	if v, ok := o.value.(period.Period); ok {
		return v, o.set
	}
	return nil, false
}

// GetPeriodDelta returns the stored value as a period and if it has been set
func (o *Value) GetPeriodDelta() (period.Delta, bool) {
	if v, ok := o.value.(period.Delta); ok {
		return v, o.set
	}
	return period.Delta(0), false
}

// GetPanning returns the stored value as a panning position and if it has been set
func (o *Value) GetPanning() (panning.Position, bool) {
	if v, ok := o.value.(panning.Position); ok {
		return v, o.set
	}
	return panning.CenterAhead, false
}

// GetPosition returns the stored value as a sample position and if it has been set
func (o *Value) GetPosition() (sampling.Pos, bool) {
	if v, ok := o.value.(sampling.Pos); ok {
		return v, o.set
	}
	return sampling.Pos{}, false
}
