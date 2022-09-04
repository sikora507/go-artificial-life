package main

import (
	"fmt"
	"math"
	"math/rand"
	"path"
	"sync"
	"time"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const width = 2700
const height = 1000
const padding = 50
const particleSize = 5

var particles []*particle = make([]*particle, 0)
var cv *canvas.Canvas
var wg sync.WaitGroup

func randomX() float64 {
	return (rand.Float64() * (float64(cv.Width()) - padding*2)) + padding
}
func randomY() float64 {
	return (rand.Float64() * (float64(cv.Height()) - padding*2)) + padding
}

func draw(p *particle) {
	cv.SetFillStyle(p.color)
	cv.FillRect(p.x, p.y, particleSize, particleSize)
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
	wg.Add(len(particles1))
	for i := 0; i < len(particles1); i++ {
		go func(i int) {
			defer wg.Done()
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
			if a.x <= 0 {
				a.vx *= -1
				a.x = 0
			} else if a.x >= float64(cv.Width()-particleSize) {
				a.vx *= -1
				a.x = float64(cv.Width()) - particleSize
			}
			if a.y <= 0 {
				a.vy *= -1
				a.y = 0
			} else if a.y >= float64(cv.Height()-particleSize) {
				a.vy *= -1
				a.y = float64(cv.Height()) - particleSize
			}
		}(i)
	}
	wg.Wait()
}

func printFps(elapsedTime time.Duration) {
	cv.SetFillStyle("#FFFFFF")
	cv.SetStrokeStyle("#000000")
	cv.SetLineWidth(2)
	fpsValue := int(1 / elapsedTime.Seconds())
	fpsText := fmt.Sprintf("FPS: %d", fpsValue)
	cv.FillText(fpsText, 5, 35)
	cv.StrokeText(fpsText, 5, 35)
}

func main() {
	var wnd *sdlcanvas.Window
	var err error
	wnd, cv, err = sdlcanvas.CreateWindow(width, height, "Artificial Life")
	if err != nil {
		panic(err)
	}

	font := path.Join("assets", "fonts", "montserrat.ttf")
	cv.SetFont(font, 32)

	yellow := create(2000, "#FFFF00")
	red := create(1000, "#FF0000")
	green := create(2000, "#00FF00")

	wnd.MainLoop(func() {
		startTime := time.Now()
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
			draw(p)
		}
		elapsedTime := time.Since(startTime)
		printFps(elapsedTime)
	})
}
