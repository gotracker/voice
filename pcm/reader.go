package pcm

import (
	"github.com/gotracker/gomixing/volume"
)

// SampleReader is a reader interface that can return a whole multichannel sample at the current position
type SampleReader interface {
	Read(arg ...volume.Matrix) (volume.Matrix, error)
}

func (s *SampleData) readData(converter SampleConverter, arg ...volume.Matrix) (volume.Matrix, error) {
	bps := converter.Size()
	actualPos := int64(s.pos * s.channels * bps)
	if actualPos < 0 {
		actualPos = 0
	}

	var out volume.Matrix
	if len(arg) > 0 {
		out = arg[0]
	} else {
		out = make(volume.Matrix, s.channels)
	}
	for c := range out {
		if err := converter.ReadAt(s, actualPos); err != nil {
			return nil, err
		}

		out[c] = converter.Volume()
		actualPos += int64(bps)
	}

	s.pos++
	return out, nil
}
