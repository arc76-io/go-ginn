package main

import (
	"fmt"
	"net/http"

	"github.com/vecno-io/arc-gfx/sector"
	"github.com/vecno-io/go-magi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SEC-L1-P0.P0
// SEC-L2-P0.P0:P0.P0
// SEC-L3-P0.P0:P0.P0:P0.P0
	
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/ARC76/SEC-L1-{l1}", func(w http.ResponseWriter, r *http.Request) {
		l1 := chi.URLParam(r, "l1")

		dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
		tk := fmt.Sprintf("L1 %s", l1)
		sc := sector.CreateGuild(tk)
		bf := sc.Render(dc)

    w.Header().Set("Content-Type", "image/png")
    w.Write(bf.Bytes())
	})
	r.Get("/ARC76/SEC-L2-{l1}:{l2}", func(w http.ResponseWriter, r *http.Request) {
		l1 := chi.URLParam(r, "l1")
		l2 := chi.URLParam(r, "l2")

		dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
		tk := fmt.Sprintf("L2 %s:%s", l1, l2)
		sc := sector.CreateRegion(tk)
		bf := sc.Render(dc)

    w.Header().Set("Content-Type", "image/png")
    w.Write(bf.Bytes())
	})
	r.Get("/ARC76/SEC-L3-{l1}:{l2}:{l3}", func(w http.ResponseWriter, r *http.Request) {
		l1 := chi.URLParam(r, "l1")
		l2 := chi.URLParam(r, "l2")
		l3 := chi.URLParam(r, "l3")

		dc := magi.NewContext(sector.CardWidth, sector.CardHeight)
		tk := fmt.Sprintf("%s:%s:%s", l1, l2, l3)
		sc := sector.CreateCluster(tk)
		bf := sc.Render(dc)

    w.Header().Set("Content-Type", "image/png")
    w.Write(bf.Bytes())
	})

	http.ListenAndServe(":8080", r)
}

