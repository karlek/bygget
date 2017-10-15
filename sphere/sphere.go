package sphere

import (
	"math"

	"github.com/karlek/bygget/ray"
	"github.com/karlek/bygget/vector"
	"github.com/karlek/bygget/world"
)

type Sphere struct {
	Center vector.Vector
	Radius float64
	world.Material
}

// http://www.scratchapixel.com/images/upload/ray-simple-shapes/raysphereisect1.png
func (s *Sphere) Hit(r *ray.Ray, min, max float64) (bool, world.Hit) {
	// Solve the quadratic equation for P^2 - R^2 = 0 to find the discriminant
	// which tells us how many (0, 1 or 2) real-valued solutions exists to our equation.
	// Hence, if we intersect the sphere or not.
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c

	hit := world.Hit{Material: s.Material}
	if discriminant <= 0 {
		return false, world.Hit{}
	}
	// We try to find the solution to the quadratic equation.
	t := (-b - math.Sqrt(b*b-a*c)) / a
	if t < max && t > min {
		hit.T = t
		hit.P = r.Point(t)
		hit.Normal = (hit.P.Sub(s.Center)).DivideScalar(s.Radius)
		return true, hit
	}
	// Try the other solution.
	t = (-b + math.Sqrt(b*b-a*c)) / a
	if t < max && t > min {
		hit.T = t
		hit.P = r.Point(t)
		hit.Normal = (hit.P.Sub(s.Center)).DivideScalar(s.Radius)
		return true, hit
	}

	return false, world.Hit{}
}

func (s *Sphere) Point() vector.Vector {
	return s.Center
}
