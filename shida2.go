package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	N  int     = 20
	xm float64 = 0.0
	ym float64 = 0.5
	h  float64 = 0.6
)

var (
	width    int    = 500
	height   int    = 500
	filename string = "shida.png"
)

var (
	bgcolor   color.Color = color.RGBA{255, 255, 255, 255}
	linecolor color.Color = color.RGBA{0, 128, 0, 255}
)

func W1x(x, y float64) float64 {
	return 0.836*x + 0.044*y
}

func W1y(x, y float64) float64 {
	return -0.044*x + 0.836*y + 0.169
}

func W2x(x, y float64) float64 {
	return -0.141*x + 0.302*y
}

func W2y(x, y float64) float64 {
	return 0.302*x + 0.141*y + 0.127
}

func W3x(x, y float64) float64 {
	return 0.141*x - 0.302*y
}

func W3y(x, y float64) float64 {
	return 0.302*x + 0.141*y + 0.169
}

func W4x(x, y float64) float64 {
	return 0
}

func W4y(x, y float64) float64 {
	return 0.175337 * y
}

func f(m *image.RGBA, k int, x, y float64, r *rand.Rand) {
	if 0 < k {
		// ch1 := make(chan int)
		ch2 := make(chan int)

		// 必ず通るものは毎回待つ。
		// go func(ch chan<- int) {
		f(m, k-1, W1x(x, y), W1y(x, y), r)
		// 	ch <- 1
		// }(ch1)

		go func(ch chan<- int) {
			if r.Float64() < 0.3 {
				f(m, k-1, W2x(x, y), W2y(x, y), r)
			}
			if r.Float64() < 0.3 {
				f(m, k-1, W2x(x, y), W2y(x, y), r)
			}
			if r.Float64() < 0.3 {
				f(m, k-1, W3x(x, y), W3y(x, y), r)
			}
			if r.Float64() < 0.3 {
				f(m, k-1, W4x(x, y), W4y(x, y), r)
			}
			ch <- 2
		}(ch2)

		// <-ch1
		<-ch2

	} else {
		var s float64 = 490.0
		m.Set(int(x*s+float64(width)*0.5), int(float64(height)-y*s), linecolor)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	src := rand.NewSource(time.Now().Unix())

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	f(m, N, 0, 0, rand.New(src))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
	}
}
