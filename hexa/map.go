package hexa

import (
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/vecno-io/go-magi"
)

type Point struct {
	X, Y int32
}

type Neuron struct {
	point Point
	value float32
}

type Setup struct {
	Seed int64

	Size int32
	Count int32
	Radius int32

	MinRadius float64
	AvrRadius float64
	MaxRadius float64

	MaxChecks  uint32
	Iterations uint32
}

type NeuronMap struct {
	rng *rand.Rand

	setup  *Setup
	layout *Layout
	
	size Point
	origin Vec2

	data []Neuron
	
	root  *Neuron
	list  []*Neuron
	slice [][]*Neuron

	core  Axial
	hexed map[Axial]bool
	store map[Axial]float32
}

func (r Point) Vec2() Vec2 {
	return Vec2{
		X: float64(r.X),
		Y: float64(r.Y),
	}
}

func NewNeuronMap(center Vec2, cfg Setup) *NeuronMap {
	layout := NewLayout(Vec2{float64(cfg.Size), float64(cfg.Size)}, center, LayoutPointy)
	scale := layout.HexSize()
	hexes := abs(-cfg.Radius-cfg.Radius) + 1
	width := int32(scale.X * float64(hexes) - float64(cfg.Size / 4))
	height := int32(scale.Y * float64(hexes) - float64(cfg.Size / 4))

	layout.Origin = Vec2{
		X: float64(width) * 0.5,
		Y: float64(height) * 0.5,
	}
	ref := &NeuronMap{
		rng:  rand.New(rand.NewSource(int64(cfg.Seed))),

		setup: &cfg,
		layout: layout,

		size: Point {
			X: width,
			Y: height,
		},
		origin: center,

		data:  make([]Neuron, width*height),
		list:  make([]*Neuron, 0, cfg.Iterations),
		slice: make([][]*Neuron, width),

		hexed: make(map[Axial]bool),
		store: make(map[Axial]float32),
	}

	// 2. Initialized all map hexes
	radius := ref.setup.Radius
	for q := -radius; q <= radius; q++ {
		r1 := max(-radius, -q - radius);
		r2 := min(radius, -q + radius);
		for r := r1; r <= r2; r++ {
			hex := Axial{q, r}
			ref.store[hex] = 0;
			if q == 0 && r == 0 {
				ref.core = hex
			}
		}
	}
	// 2. Initialized all map points
	for x := range ref.slice {
		ref.slice[x] = make([]*Neuron, height)
		for y := range ref.slice[x] {
			idx := (int(height) * x) + y
			ref.slice[x][y] = &ref.data[idx]
			ref.data[idx].point.X = int32(x)
			ref.data[idx].point.Y = int32(y)
		}
	}
	// 3. Select a random root node
	rx := int(width / 8);
	ry := int(height / 8);
	rx = ref.rng.Intn(rx) - (rx/2)
	ry = ref.rng.Intn(ry) - (ry/2)

	ref.root = ref.slice[int(width/2) + rx][int(height/2) + ry]
	ref.list = append(ref.list, ref.root)
	ref.root.value = 2.2

	return ref;
}

func (ref *NeuronMap) Logout() {
	fmt.Printf("Size Map: %d / %d\n", ref.size.X, ref.size.Y)

	fmt.Printf("Node Count: %d\n", len(ref.list))
	fmt.Printf("Total Hexes: %d\n", len(ref.store))
	fmt.Printf("Active Hexes: %d\n", len(ref.hexed))
}

func (ref *NeuronMap) Generate() {
	for i := 1; i < int(ref.setup.Iterations); i++ {
		var rn int = 0
		var ds float64 = 0
		var na *Neuron = nil
		var nb *Neuron = nil
		// 1.  Get random & nearest node (&& distance)
		// 2.  While min radius > distance between (&& max checks)
		// 2.1  goto: 1
		for (ref.setup.MinRadius > ds) && int(ref.setup.MaxChecks) > rn {
			rn++
			x := ref.rng.Intn(int(ref.size.X))
			y := ref.rng.Intn(int(ref.size.Y))
			na = ref.slice[x][y]
			nb, ds = ref.findNearest(na)
		}
		// 3.  max radius < distance between
		// 3.1  change to max distance
		if ref.setup.MaxRadius < ds {
			xc := float64(na.point.X-nb.point.X) / ds
			yc := float64(na.point.Y-nb.point.Y) / ds
			xi := int(ref.setup.AvrRadius*xc + float64(nb.point.X))
			yi := int(ref.setup.AvrRadius*yc + float64(nb.point.Y))
			na = ref.slice[xi][yi]
		}
		// 4.  min radius < distance between
		// 4.1  check ifnode is inside the map
		// 4.1.1  add node and link up
		if ref.isInRange(na.point) && ref.setup.MinRadius <= ds {
			ref.list = append(ref.list, na)
			na.value = 0.3
		}
	}
	// 5   For wanted active points
	// 5.1  get random node from list
	// 5.2  validate distance for spread
	rn := 0
	tn := 0
	ln := len(ref.list)
	er := ref.setup.AvrRadius * 2.0
	cn := min(ref.setup.Count, int32(ref.setup.Iterations / 2))
	for rn < int(cn) && tn < int(ref.setup.MaxChecks) {
		ix := ref.rng.Intn(ln)
		na := ref.list[ix]
		_, ds := ref.findNearestActive(na)
		if er < ds && na.value < 0.7 {
			na.value = 0.8
			rn++
		}
		tn++
	}
	rn = 0
	tn = 0
	for rn < int(cn) && tn < int(ref.setup.MaxChecks) {
		ix := ref.rng.Intn(ln)
		na := ref.list[ix]
		_, ds := ref.findNearestActive(na)
		if er < ds && na.value < 1.4 {
			na.value = 1.5
			rn++
		}
		tn++
	}

	// 5   For active nodes
	// 5.1  activate hex
	// 5.2  center on hex
	for _, n := range ref.list {
		if (n.value >= 0.7) {
			h := ref.layout.HexFor(n.point.Vec2())
			n.point = ref.layout.CenterFor(h).Point();
			ref.hexed[h] = true;
		}
	}
}

func (ref *NeuronMap) Render(size Point, dc *magi.Context, ps, pm, pl, px image.Image) {
	for _, n := range ref.list {
		dc.Push()
		x := n.point.X + ((size.X - ref.size.X) / 2)
		y := n.point.Y + ((size.Y - ref.size.Y) / 2)
		if (n.value < 0.7) {
			//dc.DrawImage(ps, int(x), int(y))
			dc.DrawImageAnchored(ps, int(x), int(y), 0.5, 0.5)
		} else if (n.value < 1.4) {
			//dc.DrawImage(pm, int(x), int(y))
			dc.DrawImageAnchored(pm, int(x), int(y), 0.5, 0.5)
		} else if (n.value < 2.1) {
			//dc.DrawImage(pl, int(x), int(y))
			dc.DrawImageAnchored(pl, int(x), int(y), 0.5, 0.5)
		} else {
			//dc.DrawImage(px, int(x), int(y))
			dc.DrawImageAnchored(px, int(x), int(y), 0.5, 0.5)
		}

		dc.Pop()
	}
}

func (ref *NeuronMap) RenderAll(c Vec2, dc *magi.Context, pc image.Image) {
	dc.Push()
	// For all nodes in store render image
	for k, _ := range ref.store {
		p := ref.layout.CenterFor(k)
		px := p.X + (float64(2048 - ref.size.X) * 0.5 - ref.layout.Radius.X)
		py := p.Y + (float64(2048 - ref.size.Y) * 0.5 - ref.layout.Radius.Y)
		dc.DrawImage(pc, int(px), int(py))
	}
	dc.Pop()
}

func (ref *NeuronMap) RenderActive(c Vec2, dc *magi.Context, pc image.Image) {
	dc.Push()
	// For all active nodes render image
	for k, _ := range ref.hexed {
		p := ref.layout.CenterFor(k)
		px := p.X + (float64(2048 - ref.size.X) * 0.5 - ref.layout.Radius.X)
		py := p.Y + (float64(2048 - ref.size.Y) * 0.5 - ref.layout.Radius.Y)
		dc.DrawImage(pc, int(px), int(py))
	}
	dc.Pop()
}

func (ref *NeuronMap) findNearestActive(to *Neuron) (*Neuron, float64) {
	nc := ref.root
	dc := ref.getDistance(nc, to)
	for _, n := range ref.list {
		if n.value >= 0.5 {
			d := ref.getDistance(n, to)
			if d < dc {
				dc, nc = d, n
			}
		}
	}
	return nc, dc
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

func (ref *NeuronMap) addNeuron(n *Neuron) {

}

func (ref *NeuronMap) getDistance(from *Neuron, to *Neuron) float64 {
	dx := float64(from.point.X - to.point.X)
	dy := float64(from.point.Y - to.point.Y)
	return float64(math.Sqrt((dx * dx) + (dy * dy)))
}

func (ref *NeuronMap) isInRange(pnt Point) bool {
	h := ref.layout.HexFor(pnt.Vec2())
	return h.Distance(Axial{0, 0}) < int32(ref.setup.Radius)
	// return true;
}

// func (ref *NeuronMap) shiftToTarget(hex, target Axial) Axial {
// 	dir := target.Direction(hex).Delta();
// 	return target.Add(multiply(dir, 1+ref.rng.Int31n(2)).Axial())
// }