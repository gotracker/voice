package pcm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/gotracker/gomixing/volume"
)

// Sample is the interface to a sample
type Sample interface {
	SampleReader
	Channels() int
	Length() int
	Seek(pos int)
	Tell() int
}

// SampleData is the presentation of the core data of the sample
type SampleData struct {
	pos       int // in multichannel samples
	length    int // in multichannel samples
	byteOrder binary.ByteOrder
	channels  int
	data      []byte
}

// Channels returns the channel count from the sample data
func (s *SampleData) Channels() int {
	return s.channels
}

// Length returns the sample length from the sample data
func (s *SampleData) Length() int {
	return s.length
}

// Seek sets the current position in the sample data
func (s *SampleData) Seek(pos int) {
	s.pos = pos
}

// Tell returns the current position in the sample data
func (s *SampleData) Tell() int {
	return s.pos
}

// NewSample constructs a sampler that can handle the requested sampler format
func NewSample(data []byte, length int, channels int, format SampleDataFormat) Sample {
	switch format {
	case SampleDataFormat8BitSigned:
		return &SampleReader8BitSigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat8BitUnsigned:
		return &SampleReader8BitUnsigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat16BitLESigned:
		return &SampleReader16BitSigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat16BitLEUnsigned:
		return &SampleReader16BitUnsigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat16BitBESigned:
		return &SampleReader16BitSigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.BigEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat16BitBEUnsigned:
		return &SampleReader16BitUnsigned{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.BigEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat32BitLEFloat:
		return &SampleReader32BitFloat{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat32BitBEFloat:
		return &SampleReader32BitFloat{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.BigEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat64BitLEFloat:
		return &SampleReader64BitFloat{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.LittleEndian,
				channels:  channels,
				data:      data,
			},
		}
	case SampleDataFormat64BitBEFloat:
		return &SampleReader64BitFloat{
			SampleData: SampleData{
				length:    length,
				byteOrder: binary.BigEndian,
				channels:  channels,
				data:      data,
			},
		}
	default:
		panic("unhandled sampler type")
	}
}

func ConvertTo(from Sample, format SampleDataFormat) (Sample, error) {
	cvt := &bytes.Buffer{}
	length := from.Length()
	channels := from.Channels()
	for i := 0; i < length; i++ {
		v, _ := from.Read() // ignore error
		for c := 0; c < channels; c++ {
			var vol volume.Volume
			if len(v) > c {
				vol = v[c]
			}
			switch format {
			case SampleDataFormat8BitUnsigned:
				cv := (vol * 0x80) + 0x80
				if err := binary.Write(cvt, binary.LittleEndian, uint8(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat8BitSigned:
				cv := (vol * 0x80)
				if err := binary.Write(cvt, binary.LittleEndian, int8(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitLEUnsigned:
				cv := (vol * 0x8000) + 0x8000
				if err := binary.Write(cvt, binary.LittleEndian, uint16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitLESigned:
				cv := (vol * 0x8000)
				if err := binary.Write(cvt, binary.LittleEndian, int16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitBEUnsigned:
				cv := (vol * 0x8000) + 0x8000
				if err := binary.Write(cvt, binary.BigEndian, uint16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat16BitBESigned:
				cv := (vol * 0x8000)
				if err := binary.Write(cvt, binary.BigEndian, int16(cv)); err != nil {
					return nil, err
				}
			case SampleDataFormat32BitLEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.LittleEndian, math.Float32bits(float32(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat32BitBEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.BigEndian, math.Float32bits(float32(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat64BitLEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.LittleEndian, math.Float64bits(float64(cv))); err != nil {
					return nil, err
				}
			case SampleDataFormat64BitBEFloat:
				cv := vol
				if err := binary.Write(cvt, binary.BigEndian, math.Float64bits(float64(cv))); err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unhandled format type")
			}
		}
	}
	to := NewSample(cvt.Bytes(), length, channels, format)
	return to, nil
}
