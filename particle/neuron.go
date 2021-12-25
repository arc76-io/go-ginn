package particle

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/vecno-io/go-magi"
)

type Neuron struct {
	value float32
	point Point
}

type NeuronMap struct {
	cfg  Cfg
	size Size

	val []Neuron
	rng *rand.Rand

	root  *Neuron
	list  []*Neuron
	slice [][]*Neuron
}

func NewNeuronMap(cfg Cfg, size Size) *NeuronMap {
	ref := &NeuronMap{
		cfg:  cfg,
		size: size,
		rng:  rand.New(rand.NewSource(int64(cfg.Seed))),
		val:  make([]Neuron, size.Width*size.Height),

		list:  make([]*Neuron, 0, cfg.Iterations),
		slice: make([][]*Neuron, size.Width),
	}
	for x := range ref.slice {
		ref.slice[x] = make([]*Neuron, size.Height)
		for y := range ref.slice[x] {
			idx := (int(size.Height) * x) + y
			ref.slice[x][y] = &ref.val[idx]
			ref.val[idx].point.X = x
			ref.val[idx].point.Y = y
		}
	}

	// Note: These values are based on the map
	// size and require trail and error to work.

	// Map Size: 1024
	// rx := ref.rand.Intn(128) - 64
	// ry := ref.rand.Intn(128) - 64

	// Map Size: 2048
	rx := ref.rng.Intn(256) - 128
	ry := ref.rng.Intn(256) - 128

	cx := int(size.Width/2) + rx
	cy := int(size.Height/2) + ry

	ref.root = ref.slice[cx][cy]
	ref.root.value = 1.0
	ref.list = append(ref.list, ref.root)
	return ref
}

func (ref *NeuronMap) Set(cfg Cfg) {
	ref.cfg = cfg
}

func (ref *NeuronMap) Generate() {
	for i := 1; i < int(ref.cfg.Iterations); i++ {
		var rn int = 0
		var ds float64 = 0
		var na *Neuron = nil
		var nb *Neuron = nil
		// 1.  Get random & nearest node (&& distance)
		// 2.  While min radius > distance between (&& max checks)
		// 2.1  goto: 1
		for (ref.cfg.MinRadius > ds) && int(ref.cfg.MaxChecks) > rn {
			rn++
			x := ref.rng.Intn(int(ref.size.Width))
			y := ref.rng.Intn(int(ref.size.Height))
			na = ref.slice[x][y]
			nb, ds = ref.findNearest(na)
		}
		// 3.  max radius < distance between
		// 3.1  change to max distance
		if ref.cfg.MaxRadius < ds {
			xc := float64(na.point.X-nb.point.X) / ds
			yc := float64(na.point.Y-nb.point.Y) / ds
			xi := int(ref.cfg.AvrRadius*xc + float64(nb.point.X))
			yi := int(ref.cfg.AvrRadius*yc + float64(nb.point.Y))
			na = ref.slice[xi][yi]
		}
		// 4.  min radius < distance between
		// 4.1  add node and link up
		if ref.cfg.MinRadius <= ds {
			ref.list = append(ref.list, na)
			na.value = 0.3
		}
	}
	// 5   For wanted active points
	// 5.1  get random node from list
	// 5.2  validate distance for spread
	rn := 0
	ln := len(ref.list)
	er := ref.cfg.AvrRadius * 2.0
	for rn < int(ref.cfg.Cnt) {
		ix := ref.rng.Intn(ln)
		na := ref.list[ix]
		_, ds := ref.findNearestActive(na)
		if er < ds && na.value < 0.7 {
			na.value = 0.8
			rn++
		}
	}
}

func (ref *NeuronMap) Render(c Point, dc *magi.Context) {
	for _, n := range ref.list {
		dc.Push()

		ns := 16.0
		cv := uint8(n.value * 64)

		px := float64(n.point.X + c.X)
		py := float64(n.point.Y + c.Y)

		gx := px + (ns * 0.5)
		gy := py + (ns * 0.5)

		grad := magi.NewRadialGradient(gx, gy, 4, gx, gy, ns*0.5)
		grad.AddColorStop(0, color.RGBA{cv, cv, cv + 8, 92})
		//grad.AddColorStop(0, color.RGBA{24, 0, 40, 32})
		grad.AddColorStop(1, color.RGBA{0, 0, 0, 0})
		dc.SetFillStyle(grad)
		dc.DrawRectangle(px, py, ns, ns)
		dc.Fill()

		dc.DrawPoint(px+(ns*0.5), py+(ns*0.5), 4)
		dc.SetRGBA(0.4, 0.4, 0.6, float64(0.3*n.value))
		dc.Fill()

		dc.DrawPoint(px+(ns*0.5), py+(ns*0.5), 2)
		dc.SetRGBA(0.8, 0.8, 1.0, float64(0.4*n.value))
		dc.Fill()

		dc.Pop()
	}
}

func (ref *NeuronMap) findNearest(to *Neuron) (*Neuron, float64) {
	nc := ref.root
	dc := ref.getDistance(nc, to)
	for _, n := range ref.list {
		d := ref.getDistance(n, to)
		if d < dc {
			dc, nc = d, n
		}
	}
	return nc, dc
}

func (ref *NeuronMap) findNearestActive(to *Neuron) (*Neuron, float64) {
	nc := ref.root
	dc := ref.getDistance(nc, to)
	for _, n := range ref.list {
		if n.value >= 1.0 {
			d := ref.getDistance(n, to)
			if d < dc {
				dc, nc = d, n
			}
		}
	}
	return nc, dc
}

func (ref *NeuronMap) getDistance(from *Neuron, to *Neuron) float64 {
	dx := float64(from.point.X - to.point.X)
	dy := float64(from.point.Y - to.point.Y)
	return float64(math.Sqrt((dx * dx) + (dy * dy)))
}

func (ref *NeuronMap) LogOut() {
	fmt.Printf("Cnt in list: %d\n", len(ref.list))
	fmt.Printf("Rows in slice: %d\n", len(ref.slice))
	fmt.Printf("Columns in slice: %d\n", len(ref.slice[0]))
}
