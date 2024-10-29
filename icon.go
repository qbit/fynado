package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
)

const (
	size = 300
)

type CountIcon struct {
	data    *image.RGBA
	Enabled bool
}

func (c CountIcon) Name() string {
	return "icon.png"
}

func (m CountIcon) Content() []byte {
	var buf bytes.Buffer
	err := png.Encode(&buf, m.data)
	if err != nil {
		panic(fmt.Errorf("SAD: %v", err))
	}
	return buf.Bytes()
}

func (m CountIcon) Draw(percentage float64) CountIcon {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	m.data = img

	center := size / 2
	radius := float64(center - 10)

	angle := 2 * math.Pi * percentage

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			a := float64(x - center)
			b := float64(y - center)

			if a*a+b*b <= radius*radius {
				pointAngle := math.Atan2(a, b)
				if pointAngle < 0 {
					pointAngle += 2 * math.Pi
				}

				if m.Enabled {
					if pointAngle <= angle {
						m.data.Set(x, y, blue)
					} else {
						m.data.Set(x, y, sixes)
					}
				} else {
					m.data.Set(x, y, color.Black)
				}
			}
		}
	}
	return m
}
