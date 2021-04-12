package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JimmyMcBride/elden-hub/db"
	"github.com/JimmyMcBride/elden-hub/repository"
	"github.com/JimmyMcBride/elden-hub/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
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

func MyFunk(str string) string {
	return str
}

func main() {
	bcrypt := utils.NewBCrypt(16)

	// Hash password using the salt
	hashedPassword := bcrypt.HashPassword("hello")

	fmt.Println("Password Hash:", hashedPassword)

	// Check if passed password matches the original password by hashing it
	// with the original password's salt and check if the hashes match
	fmt.Println("Password Match:", bcrypt.DoPasswordsMatch(hashedPassword, "hello"))

	engine := html.New("./views", ".gohtml")

	db := db.NewConnection()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   300 * time.Hour,
		CacheControl: true,
	}))

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		UsersRepository := repository.NewUsersRepository(db)
		user, err := UsersRepository.GetByEmail("mcbride967@gmail.com")
		if err != nil {
			fmt.Println("There was a problem getting your users table.")
			panic(err)
		}

		links := []string{"About", "Services", "Clients", "Contact"}

		return c.Render("index", fiber.Map{
			"Title":        "Hello, world!",
			"User":         user,
			"Links":        links,
			"SidebarOpen":  true,
			"MyFunk":       MyFunk("Hey!"),
			"Dependencies": []string{"https://cdn.jsdelivr.net/npm/vue@2/dist/vue.js"},
		}, "layouts/main")
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":4000"))
}
