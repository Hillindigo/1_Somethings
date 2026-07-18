package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
		PoolSize: 100,
	})

	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Printf("Redis ping response: %s\n", pong)
	return nil
}

func setUser() error {
	ctx := context.Background()
	userKey := "user:1"
	// 定义要设置的字段和值
	fields := map[string]interface{}{
		"name": "aa",
		"age":  18,
	}
	// 逐个字段设置
	for field, value := range fields {
		err := rdb.HSet(ctx, userKey, field, value).Err()
		if err != nil {
			return fmt.Errorf("set field %s failed: %w", field, err)
		}
	}
	fmt.Println("User data saved successfully")
	return nil
}

func getUser() error {
	ctx := context.Background()
	userKey := "user:1"

	user, err := rdb.HGetAll(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("get user failed: %w", err)
	}

	if len(user) == 0 {
		return fmt.Errorf("user does not exist")
	}

	fmt.Printf("User data: %+v\n", user)
	return nil
}

// watchDemo 在key值不变的情况下将其值+1
func watchDemo() {
	key := "watch_count"
	ctx := context.Background() // 定义上下文

	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		// Get 需传 ctx
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// Pipelined 需传 ctx
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			time.Sleep(2 * time.Second)
			// Set 需传 ctx
			pipe.Set(ctx, key, n+1, 0)
			return nil
		})
		return err
	}, key) // 将键名作为可变参数传递

	if err != nil {
		fmt.Println("tx exec failed, err:", err)
		return
	}
	fmt.Println("tx exec success")
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed: %v\n", err)
		return
	}
	fmt.Println("Connect redis success")
	defer rdb.Close()

	if err := setUser(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := getUser(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	watchDemo()
}
