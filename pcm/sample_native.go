package pcm

import (
	"errors"

	"github.com/gotracker/gomixing/volume"
)

var (
	// ErrIndexOutOfRange is for when a slice is iterated with an index that's out of the range
	ErrIndexOutOfRange = errors.New("index out of range")
)

type NativeSampleData struct {
	baseSampleData
	data []volume.Volume
}

// SampleReaderNative is a native (pre-converted) PCM sample reader
type SampleReaderNative struct {
	NativeSampleData
}

// Channels returns the channel count from the sample data
func (s *NativeSampleData) Channels() int {
	return s.channels
}

// Length returns the sample length from the sample data
func (s *NativeSampleData) Length() int {
	return s.length
}

// Seek sets the current position in the sample data
func (s *NativeSampleData) Seek(pos int) {
	s.pos = pos
}

// Tell returns the current position in the sample data
func (s *NativeSampleData) Tell() int {
	return s.pos
}

// Read returns the next multichannel sample
func (s *SampleReaderNative) Read(arg ...volume.Matrix) (volume.Matrix, error) {
	return s.readData(arg...)
}

func (s *NativeSampleData) readData(arg ...volume.Matrix) (volume.Matrix, error) {
	actualPos := int64(s.pos * s.channels)
	if actualPos < 0 {
		actualPos = 0
	}

	dl := len(s.data)

	var out volume.Matrix
	if len(arg) > 0 {
		out = arg[0]
	} else {
		out = make(volume.Matrix, s.channels)
	}
	start := actualPos
	end := actualPos + int64(s.channels)
	if int(end) > dl {
		return nil, ErrIndexOutOfRange
	}
	copy(out, s.data[start:end])
	s.pos++
	return out, nil
}

func NewSampleNative(data []volume.Volume, length int, channels int) Sample {
	return &SampleReaderNative{
		NativeSampleData: NativeSampleData{
			baseSampleData: baseSampleData{
				length:   length,
				channels: channels,
			},
			data: data,
		},
	}
}
