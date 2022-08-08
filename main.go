package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/vecno-io/go-magi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vecno-io/arc-gfx/hexa"
)


const cfgCardWidth = 2048
const cfgCardHeight = 2048

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// SEC-L1-P0.P0
	// SEC-L2-P0.P0:P0.P0
	// SEC-L3-P0.P0:P0.P0:P0.P0

	r.Get("/ARC76/SEC-L1-{l1}", func(w http.ResponseWriter, r *http.Request) {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: 990.0, 
				Y: 990.0,
			}, 
			hexa.Setup{
			Seed: time.Now().UnixNano(),

			Size: 34,
			Count: 18,
			Radius: 12,

			MinRadius: 49, // +1
			AvrRadius: 72,
			MaxRadius: 128,

			MaxChecks: 32,
			Iterations: 256,
		})
		hm.Generate()
		hm.Logout()
		l1 := chi.URLParam(r, "l1")

		buf := sector("GUILD", fmt.Sprintf("L1 %s", l1), hm)
    w.Header().Set("Content-Type", "image/png")
    w.Write(buf.Bytes())
	})
	r.Get("/ARC76/SEC-L2-{l1}:{l2}", func(w http.ResponseWriter, r *http.Request) {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: 990.0, 
				Y: 990.0,
			}, 
			hexa.Setup{
			Seed: time.Now().UnixNano(),

			Size: 34,
			Count: 16,
			Radius: 12,

			MinRadius: 49, // +1
			AvrRadius: 72,
			MaxRadius: 128,

			MaxChecks: 32,
			Iterations: 156,
		})
		hm.Generate()
		hm.Logout()

		l1 := chi.URLParam(r, "l1")
		l2 := chi.URLParam(r, "l2")

		buf := sector("REGION", fmt.Sprintf("L2 %s:%s", l1, l2), hm)
    w.Header().Set("Content-Type", "image/png")
    w.Write(buf.Bytes())
	})
	r.Get("/ARC76/SEC-L3-{l1}:{l2}:{l3}", func(w http.ResponseWriter, r *http.Request) {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: 990.0, 
				Y: 990.0,
			}, 
			hexa.Setup{
			Seed: time.Now().UnixNano(),

			Size: 34,
			Count: 14,
			Radius: 12,

			MinRadius: 49, // +1
			AvrRadius: 72,
			MaxRadius: 128,

			MaxChecks: 32,
			Iterations: 128,
		})
		hm.Generate()
		hm.Logout()

		l1 := chi.URLParam(r, "l1")
		l2 := chi.URLParam(r, "l2")
		l3 := chi.URLParam(r, "l3")

		buf := sector("CLUSTER", fmt.Sprintf("%s:%s:%s", l1, l2, l3), hm)
    w.Header().Set("Content-Type", "image/png")
    w.Write(buf.Bytes())
	})

	http.ListenAndServe(":8080", r)
}

func sector(cat string, token string, hexes *hexa.NeuronMap) *bytes.Buffer  {
	dc := magi.NewContext(cfgCardWidth, cfgCardHeight)
	dc.Push()
	dc.DrawRectangle(0, 0, float64(cfgCardWidth), float64(cfgCardHeight))
	dc.SetRGB(0.0322, 0, 0.0575)
	dc.Fill()
	dc.Pop()

	im, err := magi.LoadPNG("assets/frame-base-2.0.png")
	if err != nil {
		panic(err)
	}
	ha, err := magi.LoadPNG("assets/hex-pointy-64.1.png")
	if err != nil {
		panic(err)
	}
	hb, err := magi.LoadPNG("assets/hex-pointy-64.2.png")
	if err != nil {
		panic(err)
	}
	ps, err := magi.LoadPNG("assets/stars/particle-16.png")
	if err != nil {
		panic(err)
	}
	pm, err := magi.LoadPNG("assets/stars/particle-32.png")
	if err != nil {
		panic(err)
	}
	pl, err := magi.LoadPNG("assets/stars/particle-48.png")
	if err != nil {
		panic(err)
	}
	px, err := magi.LoadPNG("assets/stars/particle-64.png")
	if err != nil {
		panic(err)
	}

	dc.DrawImage(im, 0, 0)

	dc.Push()
	hexes.Render(hexa.Point{
		X: cfgCardWidth,
		Y: cfgCardHeight,
	}, dc, ps, pm, pl, px)
	dc.Pop()

	dc.Push()
	hexes.RenderAll(hexa.Vec2{
		X: cfgCardWidth / 2,
		Y: cfgCardHeight / 2,
	}, dc, ha)
	dc.Pop()

	dc.Push()
	hexes.RenderActive(hexa.Vec2{
		X: cfgCardWidth / 2,
		Y: cfgCardHeight / 2,
	}, dc, hb)
	dc.Pop()

	sx := float64(cfgCardWidth / 2.0)
	sy := float64(cfgCardHeight / 2.0)

	dc.Push()
	if err := dc.LoadFontFace("assets/font/salvar.ttf", 96); err != nil {
		panic(err)
	}
	dc.SetRGBA(1.0, 0.95, 1.0, 0.90)
	// dc.DrawStringAnchored(cat, px, py, 0.5, -19.5) // 64
	dc.DrawStringAnchored(cat, sx, sy, 0.5, -11.6)
	dc.Pop()

	buff := bytes.NewBuffer([]byte{})
	if err := png.Encode(buff, dc.Image()); err != nil {
			log.Fatal(err)
	}

	return buff
}
