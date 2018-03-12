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
	connect() *sql.DB
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

type UserService struct {
	repo UserRepository
}

type Repository interface {
	findById(id string) (User, error)
	save(name string, surname string, email string) error
	update(id string, name string, surname string, email string) error
	delete(id string) error
	getAll() ([]User, error)
}

type UserRepository struct{}

func (m *UserRepository) findById(id string) (User, error) {
	var user User
	row := db.QueryRow("select id, name, surname, email from users where id = ?;", id)
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)

	return user, err
}

func (m *UserRepository) save(name string, surname string, email string) (err error) {
	stmt, err := db.Prepare("insert into users (name, surname, email) values(?, ?, ?);")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, surname, email)
	return
}

func (m *UserRepository) update(id string, name string, surname string, email string) (err error) {
	stmt, err := db.Prepare("update users set name=?, surname=?, email=? where id=?;")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, surname, email, id)
	return
}

func (m *UserRepository) delete(id string) (err error) {
	stmt, err := db.Prepare("delete from users where id=?;")
	if err != nil {
		return
	}

	_, err = stmt.Exec(id)
	return
}

func (m *UserRepository) getAll() ([]User, error) {
	var user User
	var users []User

	rows, err := db.Query("select id, name, surname, email from users;")
	if err != nil {
		return users, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

type MysqlFactory struct {
	c Credential
}

func (m *MysqlFactory) connect() *sql.DB {
	db, err := sql.Open("mysql", m.c.Username+":"+m.c.Password+"@"+m.c.Host+"/"+m.c.Dbname)

	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

type PostgresqlFactory struct {
	c Credential
}

func (m *PostgresqlFactory) connect() *sql.DB {
	db, err := sql.Open("postgresql", "user="+m.c.Username+" dbname="+m.c.Dbname+" sslmode=verify-full")
	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

var providers = getProviders()
var db = initDB("mysql")
var userSvc = initUserService()

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

	users, err := userSvc.repo.getAll()
	if err != nil {
		fmt.Println(err.Error())
	}

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

	user, err := userSvc.repo.findById(id)
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

	err := userSvc.repo.save(name, surname, email)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record created successfully."),
	})
}

func editUser(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	email := c.PostForm("email")

	err := userSvc.repo.update(id, name, surname, email)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record updated successfully."),
	})
}

func removeUser(c *gin.Context) {
	id := c.Param("id")

	err := userSvc.repo.delete(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Record deleted successfully."),
	})
}

func initDB(driver string) *sql.DB {
	var db *sql.DB

	switch driver {
	case "mysql":
		db = connect(&MysqlFactory{providers[driver]})
	case "postgresql":
		db = connect(&PostgresqlFactory{providers[driver]})
	}
	return db
}

func initUserService() *UserService {
	return &UserService{UserRepository{}}
}

func connect(manager dbManager) *sql.DB {
	return manager.connect()
}

func getProviders() map[string]Credential {
	return map[string]Credential{
		"mysql":      Credential{Username: "root", Password: "", Host: "", Dbname: "db_gorestfulapi"},
		"postgresql": Credential{Username: "root", Password: "", Host: "", Dbname: "db_gorestfulapi"},
	}
}
