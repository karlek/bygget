package main

import (
	"flag"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/karlek/bygget/plane"
	"github.com/karlek/bygget/ray"
	"github.com/karlek/bygget/sphere"
	"github.com/karlek/bygget/vector"
	"github.com/karlek/bygget/world"
	"github.com/pkg/profile"
)

func main() {
	flag.Parse()
	defer profile.Start(profile.CPUProfile).Stop()
	if err := render(); err != nil {
		logrus.Println(err)
	}
}

type Camera struct {
	origin     vector.Vector
	horizontal vector.Vector
	vertical   vector.Vector
	lowerLeft  vector.Vector
	lensRadius float64
}

func NewCamera(lookFrom, lookAt vector.Vector, vFov, aspect, aperture float64) Camera {
	// Convert vertical fov to radians.
	theta := vFov * math.Pi / 180

	// Vertical field of view, where 180 degrees is rays with (close to)
	// infinite ascent and decent. And 0.1 degrees is really narrow.
	halfHeight := math.Tan(theta / 2)
	// Aspect ratio of height and width to produce the width.
	halfWidth := aspect * halfHeight

	// Our skewed reference system.
	//	w is our "forward" vector
	//  u our "right" vector.
	//  v our "upwards" vector.
	w := lookFrom.Sub(lookAt).Norm()
	u := vector.Vector{0, 1, 0}.Cross(w).Norm()
	v := w.Cross(u)

	focusDist := lookFrom.Sub(lookAt).Length()

	x := u.MultiplyScalar(halfWidth * focusDist)
	y := v.MultiplyScalar(halfHeight * focusDist)
	z := w.MultiplyScalar(focusDist)

	lowerLeft := lookFrom.Sub(x).Sub(y).Sub(z)
	horizontal := x.MultiplyScalar(2)
	vertical := y.MultiplyScalar(2)
	return Camera{
		origin:     lookFrom,
		horizontal: horizontal,
		vertical:   vertical,
		lowerLeft:  lowerLeft,
		lensRadius: aperture / 2,
	}
}

func (c Camera) RayAt(xscale, yscale float64, width, height int) ray.Ray {
	position := c.horizontal.MultiplyScalar(xscale).Add(c.vertical.MultiplyScalar(yscale))
	direction := c.lowerLeft.Add(position)
	return ray.Ray{
		Origin:    c.origin,
		Direction: direction}
}

func render() (err error) {

	var (
		width, height = 1280, 720
		samples       = 10
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	h := -2.0
	cam := NewCamera(vector.Origo, vector.Vector{0, h, -5.5}, 70, float64(width)/float64(height), 5.4)

	s1 := &sphere.Sphere{vector.Vector{0, h, -6}, 0.5, world.Lambertian{vector.Vector{0.8, 0.7, 0.2}}}
	// s2 := &sphere.Sphere{vector.Vector{-1, h, -6}, 0.5, world.Dielectric{0.5}}
	// s3 := &sphere.Sphere{vector.Vector{-.2, h, -4}, 0.5, world.Lambertian{vector.Vector{0.88, 0.0, 0.5}}}
	// s4 := &sphere.Sphere{vector.Vector{-0.4, h, -3}, 0.5, world.Dielectric{0.8}}
	// s5 := &sphere.Sphere{vector.Vector{1, h, -6}, 0.5, world.Metal{vector.Vector{.7, .7, .9}, 0.3}}
	floor := &plane.Plane{P: vector.Vector{0, h - 0.5, 0}, N: vector.Vector{0, -1, 0}, Material: world.Lambertian{vector.Vector{.1, .1, .1}}}

	// roof := &plane.Plane{P: vector.Vector{0, 1, 0}, N: vector.Vector{0, 1, 0}, Material: world.Lambertian{vector.Vector{.9, .9, .9}}}
	wall1 := &plane.Plane{P: vector.Vector{0, h, -9}, N: vector.Vector{0, 0, -1}, Material: world.Lambertian{vector.Vector{.9, .9, .9}}}
	wall2 := &plane.Plane{P: vector.Vector{-3, h, 0}, N: vector.Vector{-1, -1, 0}, Material: world.Lambertian{vector.Vector{.9, .9, .9}}}
	wall3 := &plane.Plane{P: vector.Vector{3, h, 0}, N: vector.Vector{1, 0, 0}, Material: world.Lambertian{vector.Vector{.9, .9, .9}}}
	// wall4 := &plane.Plane{P: vector.Vector{0, 0, 1}, N: vector.Vector{0, 0, 1}, Material: world.Lambertian{vector.Vector{1, 1, 1}}}

	disc := &plane.Disc{P: vector.Vector{0, 3, 0}, N: vector.Vector{0, 1, 0}, R: 3, Material: world.Light{vector.Vector{1.0, 1.0, 1.0}}}
	world.Bulb = disc
	// world.Bulb = &sphere.Sphere{Center: vector.Vector{0, 0, 0}, Radius: 3, Material: world.Light{vector.Vector{1.0, 1.0, 1.0}}}
	w := &world.World{Elements: []world.Hitable{}}
	w.Add(floor)
	// w.Add(roof)
	w.Add(wall1)
	w.Add(wall2)
	w.Add(wall3)
	// w.Add(wall4)
	w.Add(s1)
	// w.Add(s2)
	// w.Add(s3)
	// w.Add(s4)
	// w.Add(s5)

	wg := new(sync.WaitGroup)
	wg.Add(height)
	for y := 0; y < height; y++ {
		go func(y int, wg *sync.WaitGroup) {
			defer wg.Done()
			r := ray.Ray{}
			rnd := rand.New(rand.NewSource(int64(y)))
			for x := 0; x < width; x++ {
				xscale := (float64(x) + rand.Float64()) / float64(width)
				yscale := (float64(y) + rand.Float64()) / float64(height)
				rgb := vector.Vector{}
				for i := 0; i < samples; i++ {
					r = cam.RayAt(xscale, yscale, width, height)
					col := world.ColorVector(&r, w, 0, rnd)
					rgb = rgb.Add(col)
				}
				c := rgb.DivideScalar(float64(samples)).Color()
				img.Set(int(x), int(height-y), c)
			}
		}(y, wg)
	}
	wg.Wait()
	f, err := os.Create("a.png")
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	f.Close()
	// }
	return nil
}
