package models

type TasksBasic []TaskBasic
type Checks []Check

type TaskBasic struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Position string `json:"position"`
	Type     string `json:"type"`
}

type TaskChecker struct {
	TaskBasic
	ResultDisplay ResultDisplay `json:"result_display,omitempty"`
}

type ResultDisplay struct {
	Checks Checks `json:"checks,omitempty"`
	Done   bool   `json:"done"`
}

type Check struct {
	Title      string `json:"title"`
	Passed     bool   `json:"passed"`
	CheckLabel string `json:"check_label"`
}
