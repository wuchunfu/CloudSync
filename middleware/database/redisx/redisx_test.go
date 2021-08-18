package redisx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedis(t *testing.T) {
	as := assert.New(t)
	redisClient := NewRedisClient("127.0.0.1:6379", "123456", 0)

	t.Run("set get", func(t *testing.T) {
		err := redisClient.SetKey("key", "value", 1000)
		as.Nil(err)

		value, err := redisClient.GetValue("key")
		as.Nil(err)
		t.Log(value)
		as.Equal("value", value)
	})

	t.Run("delete one", func(t *testing.T) {
		err := redisClient.DeleteOneKey("key")
		as.Nil(err)

		value, err := redisClient.GetValue("key")
		t.Log(value)
		as.Empty(value)
	})

	t.Run("delete all", func(t *testing.T) {
		err := redisClient.DeleteOneKey("key")
		as.Nil(err)

		value, err := redisClient.GetValue("key")
		t.Log(value)
		as.Empty(value)
	})
}
