package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type DummyRepo struct {
	listCallCount   int
	listReturnValue []Todo
	listReturnError error
}

func (dr DummyRepo) list() ([]Todo, error) {
	dr.listCallCount++
	return dr.listReturnValue, dr.listReturnError
}

func (dr DummyRepo) get(id int64) (Todo, error) {
	return Todo{}, nil
}

func (dr DummyRepo) create(todo Todo) (Todo, error) {
	return todo, nil
}

func (dr DummyRepo) delete(id int64) error {
	return errors.New("")
}

func (dr DummyRepo) update(id int64, todo Todo) error {
	return nil
}

func TestList(t *testing.T) {
	listReturnValue := [3]Todo{
		{
			ID:          1,
			Description: "asdasdasd",
			IsDone:      false,
		},
		{
			ID:          2,
			Description: "2asdasdasd",
			IsDone:      true,
		},
		{
			ID:          3,
			Description: "3asdasdasd",
			IsDone:      false,
		},
	}
	dr := DummyRepo{
		listReturnValue: listReturnValue[:],
		listReturnError: nil,
	}
	userJSON := `{"name":"Jon Snow","email":"jon@labstack.com"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	if assert.NoError(t, h.list(c)) {
		assert.Equal(t, rec.Code, http.StatusOK, "Status should be OK.")
		var data [3]Todo
		json.Unmarshal(rec.Body.Bytes(), &data)
		assert.Equal(t, data, listReturnValue, "Response data should match.")
	}
}
