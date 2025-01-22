package models

import "time"

// Post represents a blog post
type Post struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Content  string    `json:"content"`
	PostTime time.Time `json:"postTime"`
	EditTime time.Time `json:"editTime"`
}

type Response struct {
	Items   interface{} `json:"items"`
	TraceId string      `json:"tx.traceId"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	TraceId string `json:"tx.traceId"`
}
