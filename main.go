package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"text/template"

	"github.com/go-webpack/webpack"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
)

// Database instance
var db *sql.DB

// Database settings
const (
	host     = "localhost"
	port     = 5432 // Default port
	user     = "postgres"
	password = "password"
	dbname   = "fiber_demo"
)

// User is a person with an account on our site.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Users is a list of all the users in the db.
type Users struct {
	Users []User `json:"user"`
}

// Connect function
func Connect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

func init() {
	// this is because public folder is shared between examples
	webpack.FsPath = "./public/webpack" // /home/wiz/go/src/github.com/JimmyMcBride/elden-hub/
}

func viewHelpers() template.FuncMap {
	return template.FuncMap{
		"asset": webpack.AssetHelper,
	}
}

func main() {
	engine := html.New("./views", ".html")
	for name, method := range viewHelpers() {
		engine.AddFunc(name, method)
	}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	isDev := flag.Bool("dev", false, "development mode")
	flag.Parse()

	webpack.IgnoreMissing = true

	webpack.Init(*isDev)

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT * FROM users order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Users{}

		for rows.Next() {
			user := User{}
			if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
				return err // Exit if we get an error
			}

			// Append User to Users
			result.Users = append(result.Users, user)
		}

		var a [2]string
		a[0] = "Hello"
		a[1] = "World"

		return c.Render("index", fiber.Map{
			"Title": "Hello, world!",
			"Words": a,
			"Users": result.Users,
		}, "layouts/main")
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":4000"))
}
