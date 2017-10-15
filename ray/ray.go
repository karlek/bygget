package ray

import "github.com/karlek/bygget/vector"

// Ray starts from an origin and continues in a direction forever.
type Ray struct {
	Origin, Direction vector.Vector
}

func (r Ray) Point(t float64) vector.Vector {
	b := r.Direction.MultiplyScalar(t)
	a := r.Origin
	return a.Add(b)
}
