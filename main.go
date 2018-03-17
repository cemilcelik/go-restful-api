package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var providers = map[string]Credential{
	"mysql": Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Dbname:   os.Getenv("DB_NAME")},
	"postgresql": Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Dbname:   os.Getenv("DB_NAME")},
}

var app App

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
	repo Repository
}

type Repository interface {
	findByID(id string) (User, error)
	save(name string, surname string, email string) error
	update(id string, name string, surname string, email string) error
	delete(id string) error
	getAll() ([]User, error)
}

type UserRepository struct{}

func (m *UserRepository) findByID(id string) (User, error) {
	var user User
	row := app.DB.QueryRow("select id, name, surname, email from users where id = ?;", id)
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email)

	return user, err
}

func (m *UserRepository) save(name, surname, email string) (err error) {
	stmt, err := app.DB.Prepare("insert into users (name, surname, email) values(?, ?, ?);")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, surname, email)
	return
}

func (m *UserRepository) update(id, name, surname, email string) (err error) {
	stmt, err := app.DB.Prepare("update users set name=?, surname=?, email=? where id=?;")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, surname, email, id)
	return
}

func (m *UserRepository) delete(id string) (err error) {
	stmt, err := app.DB.Prepare("delete from users where id=?;")
	if err != nil {
		return
	}

	_, err = stmt.Exec(id)
	return
}

func (m *UserRepository) getAll() ([]User, error) {
	var user User
	var users []User

	rows, err := app.DB.Query("select id, name, surname, email from users;")
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
	dsn := m.c.Username + ":" + m.c.Password + "@" + m.c.Host + "/" + m.c.Dbname
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

type PostgresqlFactory struct {
	c Credential
}

func (p *PostgresqlFactory) connect() *sql.DB {
	dsn := "user=" + p.c.Username + " dbname=" + p.c.Dbname + " sslmode=verify-full"
	db, err := sql.Open("postgresql", dsn)
	if err != nil {
		fmt.Print(err.Error())
	}
	return db
}

type App struct {
	Router  *gin.Engine
	DB      *sql.DB
	UserSvc *UserService
}

func (a *App) Init() {
	a.Router = initRouter()
	a.DB = initDB("mysql")
	a.UserSvc = &UserService{&UserRepository{}}
}

func (a *App) Run(port string) {
	err := a.Router.Run(port)

	if err != nil {
		panic(err)
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func main() {
	app = App{}

	app.Init()

	defer app.Close()

	app.Run(":3000")
}

func indexUserHandler(c *gin.Context) {

	users, err := app.UserSvc.repo.getAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"result": users,
		"count":  len(users),
	})
}

func getUserHandler(c *gin.Context) {
	var (
		user User
	)

	id := c.Param("id")

	user, err := app.UserSvc.repo.findByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprint("User not found."),
		})
	} else {
		c.JSON(http.StatusOK, user)
	}

}

func addUserHandler(c *gin.Context) {
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	email := c.PostForm("email")

	err := app.UserSvc.repo.save(name, surname, email)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record created successfully."),
	})
}

func editUserHandler(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	email := c.PostForm("email")

	err := app.UserSvc.repo.update(id, name, surname, email)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprint("Record updated successfully."),
	})
}

func removeUserHandler(c *gin.Context) {
	id := c.Param("id")

	err := app.UserSvc.repo.delete(id)
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
	return &UserService{&UserRepository{}}
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/user/:id", getUserHandler)
	router.GET("/users", indexUserHandler)
	router.POST("/user", addUserHandler)
	router.PATCH("/user/:id", editUserHandler)
	router.DELETE("/user/:id", removeUserHandler)

	return router
}

func connect(manager dbManager) *sql.DB {
	return manager.connect()
}
