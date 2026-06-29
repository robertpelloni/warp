package nerf

import "math"

func VolumeRender(ray Ray, near, far float64, numSamples int) Vec3 {
	samples := StratifiedSampling(ray, near, far, numSamples)

	accumulatedColor := Vec3{0, 0, 0}
	transmittance := 1.0

	for i := 0; i < len(samples)-1; i++ {
		delta := samples[i+1].T - samples[i].T
		output := EvaluateMLP(samples[i].Position, ray.Direction)

		alpha := 1.0 - math.Exp(-output.Density*delta)
		weight := alpha * transmittance

		accumulatedColor = accumulatedColor.Add(output.Color.Mul(weight))
		transmittance *= (1.0 - alpha)

		// Early ray termination
		if transmittance < 0.01 {
			break
		}
	}
	return accumulatedColor
}

func RenderScene(width, height int) {
	cameraPos := Vec3{0, 0, -3}
	lookAt := Vec3{0, 0, 0}
	fov := math.Pi / 3.0 // 60 degrees

	rays := GenerateRays(cameraPos, lookAt, fov, width, height)

	// Simple loop to demonstrate rendering
	// In reality this would write to an image buffer and parallelize
	for _, ray := range rays {
		_ = VolumeRender(ray, 1.0, 5.0, 64)
	}
}
