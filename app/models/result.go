package models

type ResultChecker struct {
	ID     string  `json:"id"`
	Checks []Check `json:"checks"`
}

type Check struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
