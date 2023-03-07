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
	"github.com/stretchr/testify/require"
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
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &Handler{todoRepo: dr}

	require.NoError(t, h.list(c))
	assert.Equal(t, http.StatusOK, rec.Code, "Status should be OK.")
	var data []Todo
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, data, listReturnValue, "Response data should match.")
}

func TestListRepositoryError(t *testing.T) {
	dr := DummyRepo{
		listReturnValue: nil,
		listReturnError: errors.New(""),
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.list(c))
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code, "Status should be 500.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "Error connecting to the database."}, data, "Response data should match.")
}

func TestCreateSuccess(t *testing.T) {
	returnedTodo := Todo{
		ID:          1,
		Description: "asd",
		IsDone:      false,
	}
	dr := DummyRepo{
		createReturnValue: returnedTodo,
		createReturnError: nil,
	}
	todoJSON := `{"description":"asd"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.create(c))
	assert.Equal(t, http.StatusCreated, rec.Code, "Status should be Created.")
	var data Todo
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, data, returnedTodo, "Response data should match.")
}

func TestCreateParseError(t *testing.T) {
	dr := DummyRepo{
		createReturnValue: Todo{},
		createReturnError: nil,
	}
	todoJSON := `{"description": 1}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.create(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, "Status should be BadRequest.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "Error parsing data."}, data, "Response data should match.")
}
func TestCreateDBError(t *testing.T) {
	dr := DummyRepo{
		createReturnValue: Todo{},
		createReturnError: errors.New(""),
	}
	todoJSON := `{"description": "asd"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.create(c))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, "Status should be 500.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "Database error."}, data, "Response data should match.")
}
func TestCreateValidationFail(t *testing.T) {
	dr := DummyRepo{
		getReturnError: errors.New(""),
	}
	todoJSON := `{"description": "", "is_done": true, "id": 1}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.create(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, "Status should be 500.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "description: non zero value required"}, data, "Response data should match.")
}
func TestDeleteSuccess(t *testing.T) {
	dr := DummyRepo{
		deleteReturnError: nil,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.delete(c))
	require.Equal(t, http.StatusNoContent, rec.Code, "Status should be NoContent.")
}
func TestDeleteWrongParam(t *testing.T) {
	dr := DummyRepo{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/asd", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("asd")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.delete(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, "Status should be BadRequest.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "ID must be integer."}, data, "Error detail must match.")
}
func TestDeleteDBError(t *testing.T) {
	dr := DummyRepo{
		deleteReturnError: errors.New(""),
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.delete(c))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, "Status should be BadRequest.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "DB error response."}, data, "Error detail must match.")

}

func TestRetrieveSuccess(t *testing.T) {
	returnTodo := Todo{
		ID:          1,
		Description: "asd",
		IsDone:      true,
	}
	dr := DummyRepo{
		getReturnValue: returnTodo,
		getReturnError: nil,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos/1", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.retrieve(c))
	require.Equal(t, http.StatusOK, rec.Code, "Status should be OK.")
	var data Todo
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, returnTodo, data, "Data must match.")
}
func TestRetrieveNotFound(t *testing.T) {
	dr := DummyRepo{
		getReturnValue: Todo{},
		getReturnError: ErrNotFound,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos/1", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.retrieve(c))
	require.Equal(t, http.StatusNotFound, rec.Code, "Status should be NotFound.")
}
func TestRetrieveWrongParam(t *testing.T) {
	dr := DummyRepo{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos/asd", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("asd")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.retrieve(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, "Status should be BadRequest.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "ID must be integer."}, data, "Error detail must match.")
}
func TestRetrieveDBError(t *testing.T) {
	dr := DummyRepo{
		getReturnError: errors.New(""),
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", strings.NewReader(``))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.retrieve(c))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, "Status should be 500.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "Failed to fetch object."}, data, "Error detail must match.")
}

func TestUpdateSuccess(t *testing.T) {
	returnedTodo := Todo{
		ID:          1,
		Description: "ASD",
		IsDone:      false,
	}
	dr := DummyRepo{
		updateReturnError: nil,
		getReturnValue:    returnedTodo,
	}
	todoJSON := `{"description":"ASD", "is_done": true}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/todos/1", strings.NewReader(todoJSON))
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

func TestUpdateValidationFail(t *testing.T) {
	dr := DummyRepo{
		getReturnError: errors.New(""),
	}
	todoJSON := `{"description": "", "is_done": true, "id": 1}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := &Handler{todoRepo: dr}

	require.NoError(t, h.update(c))
	require.Equal(t, http.StatusBadRequest, rec.Code, "Status should be BadRequest.")
	var data ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &data)
	assert.Equal(t, ErrorResponse{Detail: "description: non zero value required"}, data, "Response data should match.")
}
