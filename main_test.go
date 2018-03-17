package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"count":0,"result":null}`, rr.Body.String())
}

func TestNonExistingUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, `{"message":"User not found."}`, rr.Body.String())
}

func TestIndexUser(t *testing.T) {
	clearTable()
	populateTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	var r map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &r)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, float64(3), r["count"])
}

func TestAddUser(t *testing.T) {
	clearTable()
	populateTable()

	params := url.Values{}
	params.Add("name", "John")
	params.Add("surname", "Doe")
	params.Add("email", "john@doe.com")

	req, _ := http.NewRequest("POST", "/user", strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"message":"Record created successfully."}`, rr.Body.String())
}

func TestGetUser(t *testing.T) {
	clearTable()
	populateTable()

	req, _ := http.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	var u User
	json.Unmarshal(rr.Body.Bytes(), &u)

	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "John", u.Name)
	assert.Equal(t, "Doe", u.Surname)
	assert.Equal(t, "john@doe.com", u.Email)
}

func TestEditUser(t *testing.T) {
	clearTable()
	populateTable()

	params := url.Values{}
	params.Add("id", "1")
	params.Add("name", "John2")
	params.Add("surname", "Doe2")
	params.Add("email", "john2@doe2.com")

	req, _ := http.NewRequest("PATCH", "/user/1", strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"message":"Record updated successfully."}`, rr.Body.String())
}

func TestRemoveUser(t *testing.T) {
	clearTable()
	populateTable()

	req, _ := http.NewRequest("DELETE", "/user/3", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	req, _ = http.NewRequest("GET", "/user/3", nil)
	rr = httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func clearTable() {
	app.DB.Exec("DELETE FROM users")
	app.DB.Exec("ALTER SEQUENCE id RESTART WITH 1")
}

func populateTable() {
	users := []User{
		{1, "John", "Doe", "john@doe.com"},
		{2, "Jane", "Doe", "jane@doe.com"},
		{3, "Jack", "Doe", "jack@doe.com"},
	}

	for _, user := range users {
		app.DB.Exec("INSERT INTO users (id, name, surname, email) VALUES (?, ?, ?, ?)", user.ID, user.Name, user.Surname, user.Email)
	}
}
