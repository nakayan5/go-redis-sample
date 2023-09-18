package repository

import (
	"fmt"
	"context"
	"go-redis-sample/models"
	"time"
	"log"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db    *sqlx.DB
	cache *Cache[*models.User]
}

func NewUserRepository(
	db *sqlx.DB,
	cacheClient Client,
) *userRepository {
	return &userRepository{
		db: db,
		cache: NewCache[*models.User](
			cacheClient,
			time.Minute,
		),
	}
}

// TODO: どこでどうやって呼び出す？
func (r *userRepository) GetByID(
	ctx context.Context,
	id int,
) (*models.User, error) {
	return r.cache.GetOrSet(ctx, fmt.Sprintf("%d", id),
		func(ctx context.Context) (*models.User, error) {
			var u *models.User
			// r.dbを使ってDBから取得する

			query := `SELECT * FROM users WHERE id = :id`
			rows, err := r.db.NamedQuery(query, map[string]interface{}{"id": "1"})

			if err != nil {
				fmt.Println("ユーザーは存在しません")
			}

			err = rows.Scan(&u)

			log.Fatal(err)
			
			fmt.Println("dbに行ったぞ!!")

			return u, nil
		},
	)
}

// redisにsetするkeyを生成する
func makeUserKey(id int) string {
	return fmt.Sprintf("user:id:%d", id)
}