package main

import (
	"math"
	"math/rand"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const width = 2700
const height = 1000
const padding = 50

var particles []*particle = make([]*particle, 0)
var cv *canvas.Canvas

func randomX() float64 {
	return (rand.Float64() * (float64(cv.Width()) - padding*2)) + padding
}
func randomY() float64 {
	return (rand.Float64() * (float64(cv.Height()) - padding*2)) + padding
}

func draw(x float64, y float64, c string, s float64) {
	cv.SetFillStyle(c)
	cv.FillRect(x, y, s, s)
}

func create(number int, color string) []particle {
	group := make([]particle, number)
	for i := 0; i < number; i++ {
		group[i] = particle{x: randomX(), y: randomY(), color: color}
		particles = append(particles, &group[i])
	}
	return group
}

func rule(particles1 []particle, particles2 []particle, g float64) {
	for i := 0; i < len(particles1); i++ {
		a := &particles1[i]
		fx, fy := 0.0, 0.0
		for j := 0; j < len(particles2); j++ {
			b := &particles2[j]
			dx, dy := a.x-b.x, a.y-b.y
			if dx == 0 && dy == 0 {
				continue
			}
			d := math.Sqrt(dx*dx + dy*dy)
			if d > 80 { // only points nearby affect each other
				continue
			}
			F := g / (d)
			fx += (F * dx)
			fy += (F * dy)
		}
		a.vx = (a.vx + fx) * 0.5
		a.vy = (a.vy + fy) * 0.5
		a.x += a.vx
		a.y += a.vy
		if a.x <= 0 || a.x >= float64(cv.Width()) {
			a.vx *= -1
		}
		if a.y <= 0 || a.y >= float64(cv.Height()) {
			a.vy *= -1
		}
	}
}

func main() {
	var wnd *sdlcanvas.Window
	var err error
	wnd, cv, err = sdlcanvas.CreateWindow(width, height, "Artificial Life")
	if err != nil {
		panic(err)
	}

	yellow := create(1500, "#FFFF00")
	red := create(1500, "#FF0000")
	green := create(1500, "#00FF00")

	wnd.MainLoop(func() {
		//startTime := time.Now()
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)
		rule(green, green, -0.32)
		rule(green, red, -0.17)
		rule(green, yellow, 0.34)
		rule(red, red, -0.1)
		rule(red, green, -0.34)
		rule(yellow, yellow, 0.15)
		rule(yellow, green, -0.20)
		for _, p := range particles {
			draw(p.x, p.y, p.color, 3)
		}
		//time.Sleep(18 * time.Millisecond)
		//elapsedTime := time.Since(startTime)
		//fmt.Println(1/elapsedTime.Seconds(), "fps")
	})
}