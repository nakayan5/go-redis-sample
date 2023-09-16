package main
 
import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis-sample/models"
	"log"
	"time"
 
	"github.com/bxcodec/faker/v3"
	"github.com/redis/go-redis/v9"
)

type SomeStructWithTags struct {
	Email     string `faker:"email"`
	Name      string `faker:"name"`
	UUID      string `faker:"uuid_digit"`
	AccountID int    `faker:"oneof: 15, 27, 61"`
}

func SetupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		// docker-compose.ymlに指定したservice名+port
		Addr: "redis:6379",
		DB:   0,
	})
}

func main() {
	a := SomeStructWithTags{}
	err := faker.FakeData(&a)

	if err != nil {
		panic(err)
	}

	var userList []models.User
	for i := 0; i < 1000; i++ {
		err = faker.FakeData(&a)

		if err != nil {
			panic(err)
		}

		userList = append(userList, models.User{
			AccountID: a.AccountID,
			Name: a.Name,
			Email: a.Email,
		})
	}

	// Redisに格納するため、シリアライズ
	serialize, err := json.Marshal(&userList)
	if err != nil {
		panic(err)
	}

	// UUIDはデータにアクセスするために必要
	UUID := a.UUID
	fmt.Println("UUID: ", UUID)

	// Reedisに接続
	c := context.Background()
	r := SetupRedis()

	// Redisにデータを格納
	err = r.Set(c, UUID, serialize, time.Hour*24).Err()
	if err != nil {
		panic(err)
	}

	log.Println("complete")
}