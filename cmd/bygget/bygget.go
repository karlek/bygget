package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/Sirupsen/logrus"
)

func main() {
	flag.Parse()
	if err := render(); err != nil {
		logrus.Println(err)
	}
}

func render() (err error) {
	f, err := os.Create("a.png")
	if err != nil {
		return err
	}
	defer f.Close()

	width, height := 512, 512
	value := 256.0
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			// Red and green values range from
			// 0.0 to 1.0
			xscale := float64(x) / float64(width)
			yscale := float64(y) / float64(height)
			zscale := 0.2

			ir := uint8(value * xscale)
			ig := uint8(value * yscale)
			ib := uint8(value * zscale)

			// fmt.Println(x, y, ir, ig, ib)
			img.SetRGBA(int(x), int(y), color.RGBA{ir, ig, ib, 255})
		}
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
