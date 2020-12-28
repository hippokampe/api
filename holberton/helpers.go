package holberton

import (
	"fmt"
	"net/http"
	"time"
)

func (hbtn *Holberton) setHtml(html string) string {
	newPath := fmt.Sprintf("/%d", time.Now().UnixNano())
	hbtn.mux.HandleFunc(newPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	return hbtn.ts.URL + newPath
}
