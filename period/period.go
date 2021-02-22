package period

// Period is an interface that defines a sampler period
type Period interface {
	AddDelta(delta Delta) Period
	GetFrequency() Frequency
	GetSamplerAdd(float64) float64
}
