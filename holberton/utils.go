package holberton

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
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

func cleanString(str string) string {
	return strings.Trim(str, "\t\n ")
}

func searchTitleTask(task *goquery.Selection) (h4, span *goquery.Selection) {
	h4 = task.Find("h4.task")
	span = task.Find("h4 > span")

	return h4, span
}

func searchTaskDone(task *goquery.Selection) bool {
	hbtn := task.Find("button.student_task_done.btn.btn-default")
	classHbtn, _ := hbtn.Attr("class")
	status := strings.SplitAfter(classHbtn, "default ")[1]

	return status == "yes"
}

func parseTitleTask(h4, span *goquery.Selection) (title, class string) {
	title = strings.Replace(h4.Text(), span.Text(), "", 1)
	title = strings.Trim(title, "\t\n ")
	class = strings.Trim(span.Text(), "\t\n ")

	if class[0] == '#' { //Ex, #advanced to advanced
		class = class[1:]
	}

	return title, class
}

func parseCheck(title string) (typeCheck string, status bool) {
	title = strings.ToLower(title)

	typeCheck, status = parseCheckTitle(title)
	switch typeCheck {
	case "correct output of your code":
		return "output code", status
	case "requirement":
		return "requirement", status
	case "efficiency":
		return "efficiency", status
	case "correct answer":
		return "text answer", status
	}

	return "unknown", false
}

func parseCheckTitle(title string) (string, bool) {
	parts := strings.Split(title, "-")

	typeCheck := cleanString(parts[0])
	statusLiteral := cleanString(parts[1])

	return typeCheck, statusLiteral == "success"
}
