package vector

import (
	"math"
	"math/rand"

	"github.com/Sirupsen/logrus"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var (
	Up      = Vector{0, 1, 0}
	Right   = Vector{1, 0, 0}
	Forward = Vector{0, 0, 1}
	Origo   = Vector{0, 0, 0}
)

type Vector struct {
	X, Y, Z float64
}

// DivideScalar divides each component of the vector v with a scalar s.
func (v Vector) DivideScalar(s float64) Vector {
	return Vector{v.X / s, v.Y / s, v.Z / s}
}

// MultiplyScalar multiples the vector v with a scalar s.
func (v Vector) MultiplyScalar(s float64) Vector {
	return Vector{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector) Dot(u Vector) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vector) Cross(u Vector) Vector {
	x := v.Y*u.Z - u.Y*v.Z
	y := u.X*v.Z - v.X*u.Z
	z := v.X*u.Y - u.X*v.Y
	return Vector{x, y, z}
}

// Sub subtracts two vectors.
func (v Vector) Sub(u Vector) Vector {
	return Vector{v.X - u.X, v.Y - u.Y, v.Z - u.Z}
}

func (v Vector) AddScalar(s float64) Vector {
	return Vector{v.X + s, v.Y + s, v.Z + s}
}

func (v Vector) Color() colorful.Color {
	return colorful.LinearRgb(v.X, v.Y, v.Z)
}

func (v Vector) Mul(u Vector) Vector {
	return Vector{v.X * u.X, v.Y * u.Y, v.Z * u.Z}
}

// Add adds two vectors together.
func (v Vector) Add(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y, v.Z + u.Z}
}

func (v Vector) Norm() Vector {
	len := v.Length()
	if len == 0 {
		logrus.Println("norm: bad length")
		return v
	}
	return v.DivideScalar(len)
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

var UnitVector = Vector{1, 1, 1}

func (v Vector) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func VectorInUnitSphere(rand *rand.Rand) Vector {
	for {
		v := Vector{rand.Float64(), rand.Float64(), rand.Float64()}
		v.MultiplyScalar(2).Sub(UnitVector)
		if v.SquaredLength() >= 1 {
			return v
		}
	}
}

func (v Vector) Reflect(n Vector) Vector {
	b := n.MultiplyScalar(2 * n.Dot(v))
	return v.Sub(b)
}

func (v Vector) Refract(ov Vector, n float64) (bool, Vector) {
	uv := v.Norm()
	uo := ov.Norm()
	dt := uv.Dot(uo)

	discriminant := 1.0 - (n * n * (1 - dt*dt))
	if discriminant > 0 {
		a := uv.Sub(ov.MultiplyScalar(dt)).MultiplyScalar(n)
		b := ov.MultiplyScalar(math.Sqrt(discriminant))
		return true, a.Sub(b)
	}
	return false, Vector{}
}
