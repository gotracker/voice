package pcm

import (
	"io"
	"math"

	"github.com/gotracker/gomixing/volume"
)

const (
	//cSample32BitFloatVolumeCoeff = volume.Volume(1)
	cSample32BitFloatBytes = 4
)

// Sample32BitFloat is a 32-bit floating-point sample
type Sample32BitFloat float32

// Volume returns the volume value for the sample
func (s Sample32BitFloat) Volume() volume.Volume {
	return volume.Volume(s)
}

// Size returns the size of the sample in bytes
func (s Sample32BitFloat) Size() int {
	return cSample32BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s *Sample32BitFloat) ReadAt(d *SampleData, ofs int64) error {
	if len(d.data) <= int(ofs)+(cSample32BitFloatBytes-1) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	*s = Sample32BitFloat(math.Float32frombits(d.byteOrder.Uint32(d.data[ofs:])))
	return nil
}

// SampleReader32BitFloat is a 32-bit floating-point PCM sample reader
type SampleReader32BitFloat struct {
	SampleData
}

// Read returns the next multichannel sample
func (s *SampleReader32BitFloat) Read(arg ...volume.Matrix) (volume.Matrix, error) {
	var v Sample32BitFloat
	return s.readData(&v, arg...)
}
