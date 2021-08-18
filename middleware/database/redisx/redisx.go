package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	ctx = context.Background()
)

type ClientType struct {
	Client *redis.Client
}

func NewRedisClient(address string, pass string, DB int) *ClientType {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass, // no password set if need
		DB:       DB,   // use default DB
	})
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println(err)
	}
	log.Println(pong)
	clientType := &ClientType{
		Client: client,
	}
	return clientType
}

func (client ClientType) SetKey(key string, value interface{}, duration int) error {
	err := client.Client.Set(ctx, key, value, time.Duration(duration)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client ClientType) GetValue(key string) (string, error) {
	value, err := client.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (client ClientType) DeleteOneKey(key string) error {
	err := client.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client ClientType) DeleteAllKeys(key ...string) error {
	for i := range key {
		iter := client.Client.Scan(ctx, 0, key[i]+"*", 0).Iterator()
		for iter.Next(ctx) {
			err := client.Client.Del(ctx, iter.Val()).Err()
			if err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}
	return nil
}
