package world

import (
	"math"
	"math/rand"

	"github.com/karlek/bygget/ray"
	"github.com/karlek/bygget/vector"
)

var (
	Bulb Hitable
)

type Light struct {
	C vector.Vector
}

func (l Light) ColorVector() vector.Vector {
	return l.C
}

func (l Light) Bounce(input *ray.Ray, hit Hit, rand *rand.Rand) (bool, *ray.Ray) {
	return false, &ray.Ray{}
}

type World struct {
	Elements []Hitable
	Light    Hitable
}

func (w *World) Add(h Hitable) {
	w.Elements = append(w.Elements, h)
}

type Hitable interface {
	Hit(r *ray.Ray, min, max float64) (bool, Hit)
}

type Hit struct {
	T      float64       // Parameter t when the ray intersects the sphere.
	P      vector.Vector // The point where the ray intersects with the sphere.
	Normal vector.Vector // The surface normal of the point on the sphere.
	Material
}

func (w *World) Hit(r *ray.Ray, min, max float64) (bool, Hit) {
	hitAnything := false
	closest := max
	hit := Hit{}
	for _, elem := range w.Elements {
		hasHit, tmpHit := elem.Hit(r, min, closest)
		hitAnything = hitAnything || hasHit
		if hasHit {
			closest = tmpHit.T
			hit = tmpHit
		}
	}
	return hitAnything, hit
}

func ColorVector(r *ray.Ray, h Hitable, depth int, rand *rand.Rand) vector.Vector {
	hits, hit := h.Hit(r, 0.001, math.MaxFloat64)
	if hits {
		if depth < 1e5 {
			bounced, bouncedRay := hit.Bounce(r, hit, rand)
			if bounced {
				newColor := ColorVector(bouncedRay, h, depth+1, rand)
				return hit.Material.ColorVector().Mul(newColor)
			}
		}
		// return vector.UnitVector
		return vector.Vector{}
		// Normalize to [0, 1].
		// return hit.Normal.AddScalar(1).DivideScalar(2).Norm()
	}

	hits, hit = Bulb.Hit(r, 0.001, math.MaxFloat64)
	if hits {
		return hit.Material.ColorVector()
	} else {
		return vector.Vector{0, 0, 0}
	}

	// Make unit vector so y is between -1 and 1.
	unitDirection := r.Direction.Norm()

	// Scale t to between 0 and 1.
	t := 0.5 * (unitDirection.Y + 1.0)

	// Linear blend.
	// blended_value = (1 - t)*white + t*blue
	white := vector.Vector{1, 1, 1}
	blue := vector.Vector{0.5, 0.7, 1.0}
	return white.MultiplyScalar(1 - t).Add(blue.MultiplyScalar(t))
}

type Material interface {
	Bounce(input *ray.Ray, hit Hit, rand *rand.Rand) (bool, *ray.Ray)
	ColorVector() vector.Vector
}

type Lambertian struct {
	C vector.Vector
}

func (l Lambertian) Bounce(input *ray.Ray, hit Hit, rand *rand.Rand) (bool, *ray.Ray) {
	direction := hit.Normal.Add(vector.VectorInUnitSphere(rand))
	return true, &ray.Ray{Origin: hit.P, Direction: direction}
}

func (l Lambertian) ColorVector() vector.Vector {
	return l.C
}

type Metal struct {
	C    vector.Vector
	Fuzz float64
}

func (m Metal) Bounce(input *ray.Ray, hit Hit, rand *rand.Rand) (bool, *ray.Ray) {
	direction := input.Direction.Reflect(hit.Normal)
	fuzzed := vector.VectorInUnitSphere(rand).MultiplyScalar(m.Fuzz)
	bouncedRay := &ray.Ray{hit.P, direction.Add(fuzzed)}
	bounced := direction.Dot(hit.Normal) > 0
	return bounced, bouncedRay
}

func (m Metal) ColorVector() vector.Vector {
	return m.C
}

type Dielectric struct {
	Index float64
}

func (d Dielectric) Bounce(input *ray.Ray, hit Hit, rand *rand.Rand) (bool, *ray.Ray) {
	var outwardNormal vector.Vector
	var niOverNt, cosine float64

	if input.Direction.Dot(hit.Normal) > 0 {
		outwardNormal = hit.Normal.MultiplyScalar(-1)
		niOverNt = d.Index

		a := input.Direction.Dot(hit.Normal) * d.Index
		b := input.Direction.Length()

		cosine = a / b
	} else {
		outwardNormal = hit.Normal
		niOverNt = 1.0 / d.Index

		a := input.Direction.Dot(hit.Normal) * d.Index
		b := input.Direction.Length()

		cosine = -a / b
	}

	var success bool
	var refracted vector.Vector
	var reflectProbability float64

	if success, refracted = input.Direction.Refract(outwardNormal, niOverNt); success {
		reflectProbability = d.schlick(cosine)
	} else {
		reflectProbability = 1.0
	}

	if rand.Float64() < reflectProbability {
		reflected := input.Direction.Reflect(hit.Normal)
		return true, &ray.Ray{hit.P, reflected}
	}
	return true, &ray.Ray{hit.P, refracted}
}

func (d Dielectric) ColorVector() vector.Vector {
	return vector.Vector{1, 1, 1}
}

func (d Dielectric) schlick(cosine float64) float64 {
	r0 := (1 - d.Index) / (1 + d.Index)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
