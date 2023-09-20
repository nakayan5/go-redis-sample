package repository

import (
	"fmt"
	"context"
	// "go-redis-sample/models"
	"time"
	// "log"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db    *sqlx.DB
	cache *Cache[User]
}

type User struct {
	Id              string         `db:"id"`
	Name            string         `db:"name"`
}

func NewUserRepository(
	db *sqlx.DB,
	cacheClient Client,
) *userRepository {
	return &userRepository{
		db: db,
		cache: NewCache[User](
			cacheClient,
			time.Minute,
		),
	}
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id string,
) (User, error) {
	return r.cache.GetOrSet(ctx, id,
		func(ctx context.Context) (User, error) {
			var u = User{}

			err := r.db.Get(&u, "SELECT id, name FROM users WHERE id = ?", id)
			fmt.Println("dbに行ったぞ!!")


			if err != nil {
				fmt.Println("ユーザーは存在しません")
			}
			
			return u, nil
		},
	)
}
