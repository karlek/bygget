package plane

import (
	"github.com/karlek/bygget/ray"
	"github.com/karlek/bygget/vector"
	"github.com/karlek/bygget/world"
)

type Plane struct {
	N, P vector.Vector
	world.Material
}

// http://www.scratchapixel.com/images/upload/ray-simple-shapes/raysphereisect1.png
func (p *Plane) Hit(r *ray.Ray, min, max float64) (bool, world.Hit) {
	denom := p.N.Dot(r.Direction)
	// If they are perpendicualar.
	if denom < 1e-6 {
		return false, world.Hit{}
	}

	p0l0 := p.P.Sub(r.Origin)
	numerator := p0l0.Dot(p.N)
	t := numerator / denom
	hit := world.Hit{Material: p.Material}
	if t < max && t > min {
		hit.T = t
		hit.P = r.Point(t)
		hit.Normal = hit.P.Sub(p.P)
		return true, hit
	}
	return false, world.Hit{}
}

type Disc struct {
	N, P vector.Vector
	R    float64
	world.Material
}

// http://www.scratchapixel.com/images/upload/ray-simple-shapes/raysphereisect1.png
func (d *Disc) Hit(r *ray.Ray, min, max float64) (bool, world.Hit) {
	denom := d.N.Dot(r.Direction)
	// If they are perpendicualar.
	if denom < 1e-6 {
		return false, world.Hit{}
	}

	p0l0 := d.P.Sub(r.Origin)
	numerator := p0l0.Dot(d.N)
	t := numerator / denom
	hit := world.Hit{Material: d.Material}
	if r.Point(t).Sub(d.P).Length() > d.R {
		return false, world.Hit{}
	}
	if t < max && t > min {
		hit.T = t
		hit.P = r.Point(t)
		hit.Normal = hit.P.Sub(d.P)
		return true, hit
	}
	return false, world.Hit{}
}
