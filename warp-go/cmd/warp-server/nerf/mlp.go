package nerf

import "math"

type RenderOutput struct {
	Color   Vec3
	Density float64
}

// Dummy MLP evaluation for demonstration purposes.
func EvaluateMLP(position Vec3, direction Vec3) RenderOutput {
	// A real implementation would involve matrix multiplications and activation functions.
	// This is a placeholder that returns a synthetic sphere.
	radius := 1.0
	dist := math.Sqrt(position.X*position.X + position.Y*position.Y + position.Z*position.Z)

	if dist < radius {
		// Inside the sphere
		return RenderOutput{
			Color:   Vec3{1.0, 0.0, 0.0}, // Red sphere
			Density: 10.0,
		}
	}
	// Empty space
	return RenderOutput{
		Color:   Vec3{0, 0, 0},
		Density: 0,
	}
}
