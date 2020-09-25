package holberton

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/http/httptest"
)

func (h *Holberton) setHtml(html, path string) {
	if _, ok := h.InternalStatus.VisitedURLS[path]; ok {
		fmt.Println("You can't go again to " + path)
		return
	}

	h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
	h.InternalStatus.VisitedURLS[path] = true
}

func (h *Holberton) newServer() {
	h.mux = http.NewServeMux()
	h.ts = httptest.NewServer(h.mux)
	h.collector = colly.NewCollector()
}
