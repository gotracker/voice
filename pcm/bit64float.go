package pcm

import (
	"io"
	"math"

	"github.com/gotracker/gomixing/volume"
)

const (
	//cSample64BitFloatVolumeCoeff = volume.Volume(1)
	cSample64BitFloatBytes = 8
)

// Sample64BitFloat is a 64-bit floating-point sample
type Sample64BitFloat float64

// Volume returns the volume value for the sample
func (s Sample64BitFloat) Volume() volume.Volume {
	return volume.Volume(s)
}

// Size returns the size of the sample in bytes
func (s Sample64BitFloat) Size() int {
	return cSample64BitFloatBytes
}

// ReadAt reads a value from the reader provided in the byte order provided
func (s *Sample64BitFloat) ReadAt(d *SampleData, ofs int64) error {
	if len(d.data) <= int(ofs)+(cSample64BitFloatBytes-1) {
		return io.EOF
	}
	if ofs < 0 {
		ofs = 0
	}

	*s = Sample64BitFloat(math.Float64frombits(d.byteOrder.Uint64(d.data[ofs:])))
	return nil
}

// SampleReader64BitFloat is a 64-bit floating-point PCM sample reader
type SampleReader64BitFloat struct {
	SampleData
}

// Read returns the next multichannel sample
func (s *SampleReader64BitFloat) Read(arg ...volume.Matrix) (volume.Matrix, error) {
	var v Sample64BitFloat
	return s.readData(&v, arg...)
}
