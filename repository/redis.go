package repository
 
import (
	"encoding/json"
	"fmt"
	"context"
	"go-redis-sample/models"
	"github.com/redis/go-redis/v9"
)
 
var Cache *redis.Client

func Redis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB: 0,
	})
}

func GetUserLis(uuid string) ([]models.User, error) {
	data, err := Cache.Get(context.Background(), uuid).Result()
	if err != nil {
		panic(err)
	}

	fmt.Print(data)

	userList := new([]models.User)
	err = json.Unmarshal([]byte(data), userList)

	if err != nil {
		return nil, err
	}

	return *userList, nil
}