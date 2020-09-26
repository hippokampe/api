package holberton

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"net/http/httptest"
	"time"
)

func (h *Holberton) setHtml(html, path string) string {
	newPath := fmt.Sprintf("%s%d", path, time.Now().UnixNano())
	fmt.Println(newPath)
	h.mux.HandleFunc(newPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	return h.ts.URL + newPath
}

func (h *Holberton) newServer() {
	h.mux = http.NewServeMux()
	h.ts = httptest.NewServer(h.mux)
	h.collector = colly.NewCollector()
}

