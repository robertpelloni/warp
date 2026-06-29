package nerf

import (
	"math/rand"
)

type SamplePoint struct {
	Position Vec3
	T        float64
}

func StratifiedSampling(ray Ray, near, far float64, numSamples int) []SamplePoint {
	samples := make([]SamplePoint, numSamples)
	step := (far - near) / float64(numSamples)

	for i := 0; i < numSamples; i++ {
		t := near + float64(i)*step + rand.Float64()*step
		pos := ray.Origin.Add(ray.Direction.Mul(t))
		samples[i] = SamplePoint{Position: pos, T: t}
	}
	return samples
}
