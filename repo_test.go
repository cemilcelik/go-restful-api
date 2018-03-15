package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var userSvcMck = &UserService{&RepositoryMock{}}

// var userSvcMck = &UserService{&RepositoryMock{
// 	FindByIDFunc: func(id string) (User, error) {
// 		var user User
// 		row := db.QueryRow("select id, name, surname, email from users where id = ?;", id)
// 		err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)

// 		return user, err
// 	},
// }}

func TestFindByID(t *testing.T) {
	clearTable()
	populateTable()

	user, _ := userSvcMck.repo.findByID("1")

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "Doe", user.Surname)
	assert.Equal(t, "john@doe.com", user.Email)
}
