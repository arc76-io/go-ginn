package sector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"log"

	"github.com/vecno-io/arc-gfx/hexa"
	xor "github.com/vecno-io/arc-gfx/shifts"
	"github.com/vecno-io/go-magi"
)

const cfgHexWidth = 990.0
const cfgHexHeight = 990.0

const CardWidth = 2048
const CardHeight = 2048

type Sector struct {
	title string
	token string
	hexmap *hexa.NeuronMap
}

type JsonSector struct {
	Token		string			`json:"token"`
	Stats 	*hexa.Stats `json:"stats"`
}

func CreateGuild(token_id string) *Sector {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: cfgHexWidth, 
				Y: cfgHexHeight,
			}, 
			hexa.Setup{
			Seed: xor.GetSeed(token_id),

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
		//hm.Logout()

		return &Sector{
			title: "GUILDED",
			token: token_id,
			hexmap: hm,
		}
}

func CreateRegion(token_id string) *Sector {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: cfgHexWidth, 
				Y: cfgHexHeight,
			}, 
			hexa.Setup{
			Seed: xor.GetSeed(token_id),

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
		//hm.Logout()

		return &Sector{
			title: "REGION",
			token: token_id,
			hexmap: hm,
		}
}

func CreateCluster(token_id string) *Sector {
		hm := hexa.NewNeuronMap(hexa.Vec2{
				X: cfgHexWidth, 
				Y: cfgHexHeight,
			}, 
			hexa.Setup{
			Seed: xor.GetSeed(token_id),

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
		//hm.Logout()

		return &Sector{
			title: "CLUSTER",
			token: token_id,
			hexmap: hm,
		}
}

func (ref *Sector) Json() []byte  {
	data := JsonSector{
		Token: ref.token,
		Stats: ref.hexmap.Stats(),
	}
	js, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return js
}

func (ref *Sector) Render(dc *magi.Context) *bytes.Buffer  {
	dc.Push()
	dc.DrawRectangle(0, 0, float64(CardWidth), float64(CardHeight))
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
	pn, err := magi.LoadPNG("assets/stars/particle-12.png")
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
	ref.hexmap.Render(hexa.Point{
		X: CardWidth,
		Y: CardHeight,
	}, dc, pn, ps, pm, pl, px)
	dc.Pop()

	dc.Push()
	ref.hexmap.RenderAll(hexa.Vec2{
		X: CardWidth / 2,
		Y: CardHeight / 2,
	}, dc, ha)
	dc.Pop()

	dc.Push()
	ref.hexmap.RenderActive(hexa.Vec2{
		X: CardWidth / 2,
		Y: CardHeight / 2,
	}, dc, hb)
	dc.Pop()

	sx := float64(CardWidth / 2.0)
	sy := float64(CardHeight / 2.0)

	ss := float64(CardWidth / 5.0)

	dc.Push()
	if err := dc.LoadFontFace("assets/font/salvar.ttf", 96); err != nil {
		panic(err)
	}
	dc.SetRGBA(1.0, 0.95, 1.0, 0.90)
	// dc.DrawStringAnchored(cat, px, py, 0.5, -19.5) // 64
	dc.DrawStringAnchored(ref.title, sx, sy, 0.5, -11.725)
	dc.Pop()

	// Draw attrib tags
	dc.Push()
	if err := dc.LoadFontFace("assets/font/hacked.ttf", 32); err != nil {
		panic(err)
	}	
	dc.SetRGBA(1.0, 0.95, 1.0, 0.80)
	dc.DrawStringAnchored("Base", ss * 1, sy, 0.5, 35.6)
	dc.DrawStringAnchored("Core", ss * 2, sy, 0.5, 35.6)
	dc.DrawStringAnchored("Micro", ss * 3, sy, 0.5, 35.6)
	dc.DrawStringAnchored("Nano", ss * 4, sy, 0.5, 35.6)
	dc.Pop()

	stats := ref.hexmap.Stats()
	// Draw attrib values
	dc.Push()
	if err := dc.LoadFontFace("assets/font/hacked.ttf", 48); err != nil {
		panic(err)
	}	
	dc.SetRGBA(1.0, 0.95, 1.0, 0.80)
	dc.DrawStringAnchored(fmt.Sprintf("%03d", stats.Base), ss * 1, sy, 0.5, 25.1)
	dc.DrawStringAnchored(fmt.Sprintf("%03d", stats.Core), ss * 2, sy, 0.5, 25.1)
	dc.DrawStringAnchored(fmt.Sprintf("%03d", stats.Micro), ss * 3, sy, 0.5, 25.1)
	dc.DrawStringAnchored(fmt.Sprintf("%03d", stats.Nano), ss * 4, sy, 0.5, 25.1)
	dc.Pop()

	buff := bytes.NewBuffer([]byte{})
	if err := png.Encode(buff, dc.Image()); err != nil {
			log.Fatal(err)
	}

	return buff
}