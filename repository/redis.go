package repository
 
import (
	"fmt"
	"context"
	"time"
	"log"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Client struct {
	client *redis.Client
}

type Cache[T any] struct {
	client     Client
	expiration time.Duration
	sfg        *singleflight.Group
}

func NewCache[T any](client Client, expiration time.Duration) *Cache[T] {
	return &Cache[T]{
		client:     client,
		expiration: expiration,
		sfg:        &singleflight.Group{},
	}
}

func NewClient() *Client {
	client := redis.NewClient(&redis.Options{
		// Password:     pass,
		Addr:         "redis:6379",
		DB: 0,
		PoolSize:     20,
		MinIdleConns: 10,
	})

	return &Client{
		client: client,
	}
}

// キャッシュがあれば取得する、なければセットする
func (c *Cache[T]) GetOrSet(
	ctx context.Context,
	key string, // ユーザーのkey
	callback func(context.Context) (T, error), // キャッシュがなければDBにインサートする
) (T, error) {
	// singleflightでリクエストをまとめる
	res, err, _ := c.sfg.Do(key, func() (any, error) {
		// キャッシュから取得
		bytes, exist, err := c.client.Get(ctx, key)
		if err != nil {
			log.Println(err.Error())
		}
		if exist {
			return bytes, nil
		}
		// キャッシュがなければcallbackを実行
		t, err := callback(ctx)
		if err != nil {
			return nil, err
		}
		bytes, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		// キャッシュに保存
		err = c.client.Set(ctx, key, bytes, c.expiration)
		if err != nil {
			log.Println(err.Error())
		}
		return bytes, nil
	})

	var value T
	if err != nil {
		return value, err
	}

	bytes, ok := res.([]byte)
	if !ok {
		// 実装上、起きることはないはず
		return value, fmt.Errorf("failed to get from cache: invalid type %T", res)
	}
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}

// キャッシュを取得する
func(c *Client) Get(	
	ctx context.Context,
	key string,
)([]byte, bool, error) {
	bytes, err := c.client.Get(ctx, key).Bytes()
	// キャッシュが存在しない場合
	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to get from redis: %w", err)
	}

	// キャッシュが存在する場合
	return bytes, true, nil
}

// redisにvalueをsetする
func (c *Client) Set(
	ctx context.Context,
	key string,
	bytes []byte,
	expiration time.Duration,
) error {
	err := c.client.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set to redis: %w", err)
	}
	return nil
}