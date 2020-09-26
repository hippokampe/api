package models

type Project struct {
	Category string `json:"category,omitempty"`
	Title    string `json:"title,omitempty"`
	ID       string `json:"code,omitempty"`
	Score    string `json:"score,omitempty"`
	Tasks    []Task `json:"tasks,omitempty"`
}

type Projects []Project
