package main

import (
	"go-redis-sample/models"
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

	app.Get("users", getUserList)
	app.Post("user", createUser)
	// app.Get("user/:id")

	app.Listen(":8080")
}

// リクエストの際にダミーユーザーを生成して
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

	// return result.LastInsertId()

    return c.Status(200).SendString("success create user")
}


// TODO: キャッシュがあればそれを返す
func getUserList(c *fiber.Ctx) error {
	db := repository.Connect()

	result, err := db.Exec(`SELECT * FROM users`);

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

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