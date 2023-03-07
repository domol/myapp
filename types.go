package main

import "fmt"

type Todo struct {
	ID          int64  `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	IsDone      bool   `json:"is_done" xml:"isDone"`
}

func (t Todo) Validate() error {
	if t.Description == "" {
		return fmt.Errorf("unexpected empty description")
	}

	return nil
}
