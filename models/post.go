package models

type Post struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
	Open   bool   `json:"open"`
}
