package models

type Snippet struct {
	CreateTime string `json:"timestamp"`
	Role       string `json:"role"`
	Content    string `json:"content"`
}
