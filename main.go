package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type dbManager interface {
	connect(c Credential) *sql.DB
}

// User struct
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

// Credential struct
type Credential struct {
	Username string
	Password string
	Host     string
	Dbname   string
}

var Providers map[string]Credential

type MysqlFactory struct{}

func (m *MysqlFactory) connect(c Credential) *sql.DB {
	db, err := sql.Open("mysql", c.Username+":"+c.Password+"@"+c.Host+"/"+c.Dbname)

	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

type PostgresqlFactory struct{}

func (m *PostgresqlFactory) connect(c Credential) *sql.DB {
	db, err := sql.Open("postgresql", "user="+c.Username+" dbname="+c.Dbname+" sslmode=verify-full")
	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

var db = initDB("mysql")

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

func initDB(driver string) *sql.DB {
	var db *sql.DB

	var providers = map[string]Credential{
		"mysql":      Credential{Username: "root", Password: "", Host: "", Dbname: "db_gorestfulapi"},
		"postgresql": Credential{Username: "root", Password: "", Host: "", Dbname: "db_gorestfulapi"},
	}

	switch driver {
	case "mysql":
		db = connect(&MysqlFactory{}, providers[driver])
	case "postgresql":
		db = connect(&PostgresqlFactory{}, providers[driver])
	}
	return db
}

func connect(manager dbManager, credential Credential) *sql.DB {
	return manager.connect(credential)
}
