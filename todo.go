package main

type Todo struct {
	ID          int64  `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	IsDone      bool   `json:"is_done" xml:"isDone"`
}
