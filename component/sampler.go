package component

import (
	"github.com/gotracker/gomixing/sampling"
	"github.com/gotracker/gomixing/volume"

	"github.com/gotracker/voice/loop"
	"github.com/gotracker/voice/pcm"
)

// Sampler is a sampler component
type Sampler struct {
	sample       pcm.Sample
	pos          sampling.Pos
	keyOn        bool
	loopsEnabled bool
	wholeLoop    loop.Loop
	sustainLoop  loop.Loop

	temp0, temp1 volume.Matrix
}

// Setup sets up the sampler
func (s *Sampler) Setup(sample pcm.Sample, wholeLoop loop.Loop, sustainLoop loop.Loop) {
	s.sample = sample
	s.wholeLoop = wholeLoop
	s.sustainLoop = sustainLoop
}

// SetPos sets the current position of the sampler in the pcm data (and loops)
func (s *Sampler) SetPos(pos sampling.Pos) {
	s.pos = pos
}

// GetPos returns the current position of the sampler in the pcm data (and loops)
func (s *Sampler) GetPos() sampling.Pos {
	return s.pos
}

// Attack sets the key-on value (for loop processing)
func (s *Sampler) Attack() {
	s.keyOn = true
	s.loopsEnabled = true
}

// Release releases the key-on value (for loop processing)
func (s *Sampler) Release() {
	s.keyOn = false
}

// Fadeout disables the loops (for loop processing)
func (s *Sampler) Fadeout() {
	s.loopsEnabled = false
}

// GetSample returns a multi-channel sample at the specified position
func (s *Sampler) GetSample(pos sampling.Pos) volume.Matrix {
	v0 := s.getConvertedSample(pos.Pos, &s.temp0)
	if len(v0) == 0 {
		if s.canLoop() {
			v01 := s.getConvertedSample(pos.Pos, &s.temp1)
			panic(v01)
		}
		return v0
	}
	if pos.Frac == 0 {
		return v0
	}
	v1 := s.getConvertedSample(pos.Pos+1, &s.temp1)
	for c, s := range v1 {
		v0[c] += volume.Volume(pos.Frac) * (s - v0[c])
	}
	return v0
}

func (s *Sampler) canLoop() bool {
	switch {
	case !s.loopsEnabled:
		return false
	case s.keyOn && s.sustainLoop.Enabled():
		return true
	case s.wholeLoop.Enabled():
		return true
	}
	return false
}

func (s *Sampler) getConvertedSample(pos int, temp0 *volume.Matrix) volume.Matrix {
	if s.sample == nil {
		return volume.Matrix{}
	}
	sl := s.sample.Length()
	if pos >= sl && !s.canLoop() {
		return volume.Matrix{}
	}
	pos, _ = loop.CalcLoopPos(s.wholeLoop, s.sustainLoop, pos, sl, s.keyOn)
	if pos < 0 || pos >= sl {
		return volume.Matrix{}
	}
	s.sample.Seek(pos)
	ch := s.sample.Channels()
	if len(*temp0) != ch {
		*temp0 = make(volume.Matrix, ch)
	}
	data, err := s.sample.Read(*temp0)
	if err != nil {
		return volume.Matrix{}
	}
	return data
}
