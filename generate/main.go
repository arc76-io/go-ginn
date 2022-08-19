package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vecno-io/arc-gfx/sector"

	"github.com/vecno-io/go-magi"
	imaging "github.com/vecno-io/go-magin"
)

// SEC-L1-P0.P0
// SEC-L2-P0.P0:P0.P0
// SEC-L3-P0.P0:P0.P0:P0.P0

type JsonData struct {
	Media					string	`json:"media"`
	Media_hash		string	`json:"media_hash"`
	Refrence			string	`json:"refrence"`
	Refrence_hash	string	`json:"refrence_hash"`
}

type JsonRarity struct {
	Base	map[uint32][]string	`json:"base"`
	Core	map[uint32][]string	`json:"core"`
	Micro	map[uint32][]string	`json:"micro"`
	Nano	map[uint32][]string	`json:"nano"`
	Active	map[uint32][]string	`json:"active"`
}

var rarity = &JsonRarity {
	Base: make(map[uint32][]string),
	Core: make(map[uint32][]string),
	Micro: make(map[uint32][]string),
	Nano: make(map[uint32][]string),
	Active: make(map[uint32][]string),
} 

func main() {
	// // SEC-L1-P0.P0
	l1 := tileKey(0, 0)
	title := fmt.Sprintf("SEC-L1-%s", l1)
	fmt.Println(title)

	dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
	sc := sector.CreateGuild(title)
	sc.Render(dc)
	dc.SavePNG(fmt.Sprintf("./out/%s.png", title))
	scaleImage(title)

	// TODO IPFS Uploads
	writeJson(title, sc)
	writeData(title)
	addStats(sc)

	buildRegion(l1)
	writeStats()
}

func hashFile(file string) string {
  f, err := os.Open(file)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  h := sha256.New()
  if _, err := io.Copy(h, f); err != nil {
    log.Fatal(err)
  }

  return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func scaleImage(title string) {
	path := fmt.Sprintf("./out/%s.png", title)
	srcImage, err := imaging.Open(path)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	dst1024 := imaging.Resize(srcImage, 1024, 1024, imaging.Lanczos)
	path = fmt.Sprintf("./out/1024/%s.png", title)
	err = imaging.Save(dst1024, path)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	dst512 := imaging.Resize(srcImage, 512, 512, imaging.Lanczos)
	path = fmt.Sprintf("./out/512/%s.png", title)
	err = imaging.Save(dst512, path)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func writeData(title string) {
	data := JsonData{
		Media: fmt.Sprintf("/%s.png", title),
		Media_hash: hashFile(fmt.Sprintf("./out/%s.json", title)),
		Refrence: fmt.Sprintf("/%s.json", title),
		Refrence_hash: hashFile(fmt.Sprintf("./out/%s.png", title)),
	}
	js, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fl, err := os.Create(fmt.Sprintf("./out/%s.meta", title))
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	if _, err := fl.Write(js); err != nil {
		panic(err)
	}
	fl.Sync()
}

func writeJson(title string, sc *sector.Sector) {
	fl, err := os.Create(fmt.Sprintf("./out/%s.json", title))
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	if _, err := fl.Write(sc.Json()); err != nil {
		panic(err)
	}
	fl.Sync()
}

func writeStats() {
	fl, err := os.Create("./out/rarity.json")
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	js, err := json.MarshalIndent(rarity, "", " ")
	if err != nil {
		panic(err)
	}
	if _, err := fl.Write(js); err != nil {
		panic(err)
	}
	fl.Sync()
}

func buildRegion(l1 string) {
	// SEC-L2-P0.P0:PX.PX
	for q := int32(-2); q <= 2; q++ {
		r1 := max(-2, -q - 2);
		r2 := min(2, -q + 2);
		for r := r1; r <= r2; r++ {
			l2 := tileKey(q, r)
			title := fmt.Sprintf("SEC-L2-%s:%s", l1, l2)
			fmt.Println(title)

			dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
			sc := sector.CreateRegion(title)
			sc.Render(dc)
			dc.SavePNG(fmt.Sprintf("./out/%s.png", title))
			scaleImage(title)

			// TODO IPFS Uploads
			writeJson(title, sc)
			writeData(title)
			addStats(sc)

			buildCluster(l1, l2)
		}
	}
}

func buildCluster(l1, l2 string) {
	// SEC-L3-P0.P0:PX.PX:PX.PX
	for q := int32(-3); q <= 3; q++ {
		r1 := max(-3, -q - 3);
		r2 := min(3, -q + 3);
		for r := r1; r <= r2; r++ {
			l3 := tileKey(q, r)
			title := fmt.Sprintf("SEC-L3-%s:%s:%s", l1, l2, l3)
			fmt.Println(title)

			dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
			sc := sector.CreateCluster(title)
			sc.Render(dc)
			dc.SavePNG(fmt.Sprintf("./out/%s.png", title))
			scaleImage(title)

			// TODO IPFS Uploads
			writeJson(title, sc)
			writeData(title)
			addStats(sc)
		}
	}
}

func addStats(sc *sector.Sector) {
	data := sc.Data()

	rarity.Base[data.Stats.Base] = append(rarity.Base[data.Stats.Base], data.Token)
	rarity.Core[data.Stats.Core] = append(rarity.Core[data.Stats.Core], data.Token)
	rarity.Micro[data.Stats.Micro] = append(rarity.Micro[data.Stats.Micro], data.Token)
	rarity.Nano[data.Stats.Nano] = append(rarity.Nano[data.Stats.Nano], data.Token)
	rarity.Active[data.Stats.Active] = append(rarity.Active[data.Stats.Active], data.Token)
}

func abs(x int32) int32 {
	if x > 0 {
		return x;
	}
	return -x;
}

func max(a, b int32) int32 {
	if a > b {
			return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
			return a
	}
	return b
}

func tileKey(q, r int32) string {
	qs := "P"
	rs := "P"
	if q < 0 {
		qs = "N"
	}
	if r < 0 {
		rs = "N"
	}
	return fmt.Sprintf("%s%d.%s%d", qs, abs(q), rs, abs(r))
}