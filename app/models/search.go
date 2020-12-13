package models

type ProjectsResultSearch struct {
	Query string `json:"query"`
	Results []ProjectSearch `json:"results"`
}

type ProjectSearch struct {
	ID string `json:"id"`
	Score float64 `json:"score"`
}
