package models

type Projects []Project

type Project struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Score    string `json:"score"`
}
