package models

type Project struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	ID       string `json:"code"`
	Score    string `json:"score"`
	Tasks    []Task `json:"tasks"`
}

type Projects []Project
