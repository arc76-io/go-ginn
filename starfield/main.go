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

	BuildFieldPath(dc)
	//BuildFieldDefault(dc)
	DrawHex(dc, maps.Vec2{
		cfgCardWidth / 2,
		cfgCardHeight / 2,
	}, false)

	// Right
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 56,
		(cfgCardHeight / 2) + ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 111,
		(cfgCardHeight / 2),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 168,
		(cfgCardHeight / 2) + ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 56,
		(cfgCardHeight / 2) - ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 222,
		(cfgCardHeight / 2),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) + 168,
		(cfgCardHeight / 2) - ((SIZE * 2) * 0.75),
	}, false)

	// Left
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 56,
		(cfgCardHeight / 2) + ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 111,
		(cfgCardHeight / 2),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 168,
		(cfgCardHeight / 2) + ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 56,
		(cfgCardHeight / 2) - ((SIZE * 2) * 0.75),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 222,
		(cfgCardHeight / 2),
	}, false)
	DrawSpaceHex(dc, maps.Vec2{
		(cfgCardWidth / 2) - 168,
		(cfgCardHeight / 2) - ((SIZE * 2) * 0.75),
	}, false)

	dc.SavePNG("starfield.png")
}

////////////////////////////////////////////////////////////////

const PI = float64(3.14159265)
const SIZE = float64(64.0)

func DrawHex(dc *magi.Context, c maps.Vec2, flat bool) {

	list := [6]maps.Vec2{}
	if flat {
		for i := uint(1); i < 7; i++ {
			pnt := flatHexCorner(c, SIZE, i)
			list[i-1] = pnt
		}
	} else {
		for i := uint(1); i < 7; i++ {
			pnt := pointHexCorner(c, SIZE, i)
			list[i-1] = pnt
		}
	}

	dc.SetRGBA(0.8, 0.8, 1.0, 1.0)
	dc.SetLineWidth(2.0)

	old := list[len(list)-1]
	for _, pnt := range list {
		dc.DrawLine(
			old[0], old[1],
			pnt[0], pnt[1],
		)
		dc.Stroke()
		old = pnt
	}
}

func DrawSpaceHex(dc *magi.Context, c maps.Vec2, flat bool) {
	list := [6]maps.Vec2{}
	if flat {
		for i := uint(1); i < 7; i++ {
			pnt := flatHexCorner(c, 64.0, i)
			list[i-1] = pnt
		}
	} else {
		for i := uint(1); i < 7; i++ {
			pnt := pointHexCorner(c, 64.0, i)
			list[i-1] = pnt
		}
	}

	dc.SetRGBA(0.8, 0.8, 1.0, 1.0)
	dc.SetLineWidth(2.0)

	old := list[len(list)-1]
	for _, pnt := range list {

		// a + (b - a).norm() * d
		a := old
		b := pnt
		maps.Vec2Sub(&b, &a)
		maps.Vec2Normalize(&b)
		maps.Vec2Mulf(&b, 8.0)

		o := old
		p := pnt
		maps.Vec2Add(&o, &b)
		maps.Vec2Sub(&p, &b)

		dc.DrawLine(
			o[0], o[1],
			p[0], p[1],
		)
		dc.Stroke()
		old = pnt
	}
}

func flatHexCorner(c maps.Vec2, s float64, i uint) maps.Vec2 {
	deg := float64(60.0 * i)
	rad := PI / 180.0 * deg
	return maps.Vec2{
		c[0] + (s * math.Cos(rad)),
		c[1] + (s * math.Sin(rad)),
	}
}

func pointHexCorner(c maps.Vec2, s float64, i uint) maps.Vec2 {
	deg := float64(60.0*i - 30.0)
	rad := PI / 180.0 * deg
	return maps.Vec2{
		c[0] + s*math.Cos(rad),
		c[1] + s*math.Sin(rad),
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

func BuildFieldPath(dc *magi.Context) {
	seed := uint64(time.Now().UnixNano())
	ref := maps.NewNeuronMap(maps.Cfg{
		Cnt:  32,
		Seed: seed,

		MinRadius: 33, // +1
		AvrRadius: 48,
		MaxRadius: 64,

		MaxChecks:  384,
		Iterations: 384,
	}, maps.Size{
		Width:  cfgCardWidth - 384,
		Height: cfgCardHeight - 384,
	})

	ref.Generate()
	ref.LogOut()

	ref.Set(maps.Cfg{
		Cnt:  16,
		Seed: uint64(time.Now().UnixNano()),

		MinRadius: 289, // +1
		AvrRadius: 336,
		MaxRadius: 384,

		MaxChecks:  64,
		Iterations: 64,
	})
	ref.Generate()
	ref.LogOut()

	ref.Set(maps.Cfg{
		Cnt:  8,
		Seed: uint64(time.Now().UnixNano()),

		MinRadius: 49, // +1
		AvrRadius: 96,
		MaxRadius: 96,

		MaxChecks:  128,
		Iterations: 128,
	})
	ref.Generate()
	ref.LogOut()

	ref.Render(maps.Point{
		X: 384 / 2,
		Y: 384 / 2,
	}, dc)
}
