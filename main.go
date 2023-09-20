package main

import (
	"context"
	"log"
	"fmt"
	"go-redis-sample/models"
	"go-redis-sample/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	repository.Setup()	
	
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("user/:id", getUser)
	app.Post("user", createUser)

	app.Listen(":8080")
}

func createUser(c *fiber.Ctx) error {
	db := repository.Connect()

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(404).SendString(err.Error())
	}

	_, err := db.NamedExec(`INSERT INTO users (id, name) VALUES (:id, :name);`, user)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

    return c.Status(200).SendString("success create user")
}


func getUser(c *fiber.Ctx) error {
	db := repository.Connect()
	client := repository.NewClient()
	context := context.Background()
	id := c.Params("id")

	log.Println(fmt.Sprintf("idは%sです", id))

    userRepository := repository.NewUserRepository(db, *client)

	user, err := userRepository.GetByID(context, id)


	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(user)
}