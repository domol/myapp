package main

import "fmt"

type Todo struct {
	ID          int64  `json:"id" xml:"id"`
	Description string `json:"description" xml:"description" valid:"required,maxstringlength(100)"`
	IsDone      bool   `json:"is_done" xml:"isDone"`
}

func (t Todo) Validate() error {
	if t.Description == "" {
		return fmt.Errorf("unexpected empty description")
	}

	return nil
}
