package main

import (
	"math"
	"time"

	"github.com/vecno-io/go-magi"

	maps "github.com/vecno-io/arc-gfx/particle"
)

const cfgCardWidth = 2100
const cfgCardHeight = 3450

func main() {
	dc := magi.NewContext(cfgCardWidth, cfgCardHeight)

	dc.Push()
	dc.DrawRoundedRectangle(0, 0, cfgCardWidth, cfgCardHeight, 64)
	dc.SetRGB(0.035, 0, 0.075)
	dc.Fill()
	dc.Pop()

	BuildFieldDefault(dc)
	DrawHex(dc)

	dc.SavePNG("starfield.png")
}

////////////////////////////////////////////////////////////////

const PI = float64(3.14159265)

func DrawHex(dc *magi.Context) {
	c := maps.Vec2{
		X: cfgCardWidth / 2,
		Y: cfgCardHeight / 2,
	}

	list := [6]maps.Vec2{}
	for i := uint(1); i < 7; i++ {
		pnt := hexCorner(c, 256.0, i)
		list[i-1] = pnt
	}

	dc.SetRGBA(0.8, 0.8, 1.0, 1.0)
	dc.SetLineWidth(4.0)

	old := list[len(list)-1]
	for _, pnt := range list {
		dc.DrawLine(
			old.X, old.Y,
			pnt.X, pnt.Y,
		)
		dc.Stroke()
		old = pnt
	}
}

func hexCorner(c maps.Vec2, s float64, i uint) maps.Vec2 {
	deg := float64(60.0*i - 30.0)
	rad := PI / 180.0 * deg
	return maps.Vec2{
		X: c.X + (s * math.Cos(rad)),
		Y: c.Y + (s * math.Sin(rad)),
	}
}

////////////////////////////////////////////////////////////////

func BuildFieldDefault(dc *magi.Context) {
	seed := uint64(time.Now().UnixNano())
	ref := maps.NewNeuronMap(maps.Cfg{
		Cnt:  8,
		Seed: seed,

		// MinRadius: 128, // 33, // +1
		// AvrRadius: 192, // 48,
		// MaxRadius: 256, // 64,

		MinRadius: 289, // +1
		AvrRadius: 336,
		MaxRadius: 384,

		MaxChecks:  24,
		Iterations: 24,
	}, maps.Size{
		Width:  cfgCardWidth - 384,
		Height: cfgCardHeight - 384,
	})
	ref.Generate()
	ref.LogOut()

	// This is not what we want, we want to map the old
	// points in to a new map before generating it again.
	ref.Set(maps.Cfg{
		Cnt:  32,
		Seed: uint64(time.Now().UnixNano()),

		MinRadius: 49, // +1
		AvrRadius: 96,
		MaxRadius: 96,

		MaxChecks:  192,
		Iterations: 192,
	})
	ref.Generate()
	ref.LogOut()

	ref.Render(maps.Point{
		X: 384 / 2,
		Y: 384 / 2,
	}, dc)
}

////////////////////////////////////////////////////////////////

func MakeDefault(seed uint64) *maps.NeuronMap {
	return maps.NewNeuronMap(maps.Cfg{
		Cnt:  32,
		Seed: seed,

		MinRadius: 33, // +1
		AvrRadius: 48,
		MaxRadius: 64,

		MaxChecks:  256,
		Iterations: 256,
	}, maps.Size{
		Width:  2048,
		Height: 2048,
	})
}
