package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindByID(t *testing.T) {
	clearTable()
	populateTable()

	user, _ := UserSvcMck.repo.findByID("1")

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "Doe", user.Surname)
	assert.Equal(t, "john@doe.com", user.Email)
}
