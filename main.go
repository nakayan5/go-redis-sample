package main

import (
	"go-redis-sample/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// reids
	repository.Setup()
	// repository.Redis()
	
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// app.Get("users/:uuid", getUserList)
	app.Get("user", createUser)

	app.Listen(":8080")
}

// リクエストの際にダミーユーザーを生成して
func createUser(c *fiber.Ctx) error {
	db := repository.Connect()

	result, err := db.Exec("INSERT INTO users (`id`, `name`) VALUES (1, '中村')")

	if err != nil {
		return err
	}

	// return result.LastInsertId()

    return c.JSON(result)
}

// func getUserList(c *fiber.Ctx) error {
// 	uuid := c.Params("uuid")

// 	userList, err := repository.GetUserLis(uuid)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return c.JSON(userList)
// }