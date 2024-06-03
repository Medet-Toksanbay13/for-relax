package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var database *sql.DB

func home(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		id := c.FormValue("id")
		name := c.FormValue("name")
		age := c.FormValue("age")

		_, err := database.Exec("INSERT INTO public.fortestdb (id, name, age) VALUES ($1, $2, $3)", id, name, age)
		if err != nil {
			log.Println(err)
		}
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	return c.File("hello.html")
}

func secPage(c echo.Context) error {
	var u User
	rows, err := database.Query("SELECT * FROM public.fortestdb")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	user := []User{}

	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			panic(err)
		}
		user = append(user, u)
	}
	tmpl, _ := template.ParseFiles("second.html")
	return tmpl.Execute(c.Response().Writer, user)
}

func main() {
	conn := "user=postgres password=meda13 dbname=fortestdb sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	e := echo.New()
	e.GET("/", home)
	e.POST("/", home)
	e.GET("/users", secPage)

	fmt.Println("loading...")
	e.Start(":8181")
}
