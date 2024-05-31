package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var database *sql.DB

func homeHandler(c echo.Context) error {
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

func main() {
	conn := "user=postgres password=meda13 dbname=fortestdb sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	e := echo.New()
	e.GET("/", homeHandler)
	e.POST("/", homeHandler)

	fmt.Println("loading...")
	e.Start(":8181")
}
