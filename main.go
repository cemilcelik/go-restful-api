package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// User struct
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

var db = initDB()

func main() {
	defer db.Close()

	router := gin.Default()

	router.GET("/user/:id", getUser)
	router.GET("/users", indexUser)
	router.POST("/user", addUser)
	router.PATCH("/user/:id", editUser)
	router.DELETE("/user/:id", removeUser)

	router.Run(":3000")
}

func indexUser(c *gin.Context) {
	var (
		user  User
		users []User
	)

	rows, err := db.Query("select id, name, surname, email from users;")
	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)
		users = append(users, user)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"result": users,
		"count":  len(users),
	})
}

func getUser(c *gin.Context) {
	var (
		user   User
		result gin.H
	)

	id := c.Param("id")

	row := db.QueryRow("select id, name, surname, email from users where id = ?;", id)

	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)
	if err != nil {
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": user,
			"count":  1,
		}
	}

	c.JSON(http.StatusOK, result)
}

func addUser(c *gin.Context) {
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	email := c.PostForm("email")

	stmt, err := db.Prepare("insert into users (name, surname, email) values(?, ?, ?);")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(name, surname, email)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record created successfully."),
	})
}

func editUser(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	email := c.PostForm("email")

	stmt, err := db.Prepare("update users set name=?, surname=?, email=? where id=?;")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(name, surname, email, id)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record updated successfully."),
	})
}

func removeUser(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id:", id)
	stmt, err := db.Prepare("delete from users where id=?;")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Record deleted successfully."),
	})
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@/db_gorestfulapi")
	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}
