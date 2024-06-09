package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var database *sql.DB

func homee(c echo.Context) error {
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
	tmpl, _ := template.ParseFiles("get.html")
	tmpl.Execute(c.Response().Writer, user)

	return c.File("create.html")
}

func createUser(c echo.Context) error {
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
	return c.JSON(http.StatusMethodNotAllowed, map[string]string{"message": "Method Not Allowed"})
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	_, err = database.Exec("DELETE FROM public.fortestdb WHERE id = $1", idInt)
	if err != nil {
		log.Println(err)
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editPage(c echo.Context) error {
	id := c.Param("id")

	row := database.QueryRow("select * from public.fortestdb where id = $1", id)
	user := User{}
	err := row.Scan(&user.Id, &user.Name, &user.Age)
	if err != nil {
		log.Println(err)
	}
	tmpl, _ := template.ParseFiles("update.html")
	tmpl.Execute(c.Response().Writer, user)

	return nil
}

func editUser(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		id := c.FormValue("id")
		name := c.FormValue("name")
		age := c.FormValue("age")

		_, err := database.Exec("update public.fortestdb set name=$1, age=$2 where id = $3",
			name, age, id)
		if err != nil {
			log.Println(err)
		}
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	return c.JSON(http.StatusMethodNotAllowed, map[string]string{"message": "Method Not Allowed"})
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

	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" {
				method := c.FormValue("_method")
				if method != "" {
					c.Request().Method = method
				}
			}
			return next(c)
		}
	}))
	e.GET("/", homee)
	e.POST("/", createUser)
	e.GET("/edit/:id", editPage)
	e.POST("/edit/:id", editUser)
	e.POST("/delete/:id", deleteUser)

	fmt.Println("loading...")
	e.Start(":8181")
}
