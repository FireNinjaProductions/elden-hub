package main

import (
	"fmt"
	"log"

	"github.com/JimmyMcBride/elden-hub/db"
	"github.com/JimmyMcBride/elden-hub/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
)

// User is a person with an account on our site.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Users is a list of all the users in the db.
type Users struct {
	Users []User `json:"user"`
}

func main() {
	engine := html.New("./views", ".html")

	db := db.NewConnection()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		UsersRepository := repository.NewUsersRepository(db)
		user, err := UsersRepository.GetByEmail("mcbride967@gmail.com")
		if err != nil {
			fmt.Println("There was a problem getting your users table.")
			panic(err)
		}

		var a [2]string
		a[0] = "Hello"
		a[1] = "World"

		return c.Render("index", fiber.Map{
			"Title":       "Hello, world!",
			"Words":       a,
			"User":        user,
			"Description": "A refuge for hollowed souls.",
		}, "layouts/main")
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":4000"))
}
