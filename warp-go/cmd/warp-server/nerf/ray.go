package nerf

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vec3) Mul(scalar float64) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) Normalize() Vec3 {
	len := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if len == 0 {
		return Vec3{}
	}
	return Vec3{v.X / len, v.Y / len, v.Z / len}
}

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

func GenerateRays(cameraPos Vec3, lookAt Vec3, fov float64, width, height int) []Ray {
	var rays []Ray
	// Simplified ray generation for demonstration purposes.
	forward := lookAt.Sub(cameraPos).Normalize()
	up := Vec3{0, 1, 0}
	right := crossProduct(forward, up).Normalize()
	actualUp := crossProduct(right, forward).Normalize()

	aspectRatio := float64(width) / float64(height)
	tanFov := math.Tan(fov / 2)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ndcX := (float64(x) + 0.5) / float64(width)
			ndcY := (float64(y) + 0.5) / float64(height)

			px := (2*ndcX - 1) * aspectRatio * tanFov
			py := (1 - 2*ndcY) * tanFov

			dir := forward.Add(right.Mul(px)).Add(actualUp.Mul(py)).Normalize()
			rays = append(rays, Ray{Origin: cameraPos, Direction: dir})
		}
	}
	return rays
}

func crossProduct(v1, v2 Vec3) Vec3 {
	return Vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}
