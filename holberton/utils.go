package holberton

import (
	"github.com/gocolly/colly"
	"net/http"
	"net/http/httptest"
)

func (h *Holberton) setHtml(html, path string) {
	h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}

func (h *Holberton) newServer() {
	h.mux = http.NewServeMux()
	h.ts = httptest.NewServer(h.mux)
	h.collector = colly.NewCollector()
}
