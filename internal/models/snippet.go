package models

type Snippet struct {
	CreateTime float64 `json:"timestamp"`
	Role       string  `json:"role"`
	Content    string  `json:"content"`
}
