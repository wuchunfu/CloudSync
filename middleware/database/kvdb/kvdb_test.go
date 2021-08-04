package kvdb

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewKvDb(t *testing.T) {
	as := assert.New(t)
	defer os.Remove("./test.data")
	defer os.Remove("./testBackup.data")

	c := NewKvDb("./test.data", "testBucket")

	t.Run("not found", func(t *testing.T) {
		v, err := c.GetByKey("k")
		as.Equal(nil, err)
		t.Log(string(v))
		as.Empty(v)
	})

	t.Run("set key and get key", func(t *testing.T) {
		set := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}
		as.Nil(c.Set(set))

		v, err := c.GetByKey("key1")
		as.Nil(err)
		t.Log(string(v))
		as.Equal("value1", string(v))
	})

	t.Run("set keys and get keys", func(t *testing.T) {
		set := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}
		as.Nil(c.Set(set))

		keys := []string{"key1", "key2"}
		kvs, err := c.GetByKeys(keys)
		as.Nil(err)
		for k, v := range kvs {
			t.Log(k, string(v))
		}
	})

	t.Run("set keys and get all", func(t *testing.T) {
		set := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}
		as.Nil(c.Set(set))

		kvs, err := c.GetAll()
		as.Nil(err)
		for k, v := range kvs {
			t.Log(k, string(v))
		}
	})

	t.Run("set keys and delete key", func(t *testing.T) {
		set := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}
		as.Nil(c.Set(set))

		keys := []string{"key1"}
		err := c.DeleteByKeys(keys)
		as.Nil(err)

		v, err := c.GetByKey("key1")
		as.Nil(err)
		t.Log(string(v))
		as.Empty(string(v))
	})

	t.Run("set keys and delete key by transaction", func(t *testing.T) {
		set := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}
		as.Nil(c.Set(set))

		keys := []string{"key1"}
		err := c.DeleteByKeysTransaction(keys)
		as.Nil(err)

		v, err := c.GetByKey("key1")
		as.Nil(err)
		t.Log(string(v))
		as.Empty(string(v))
	})

	t.Run("delete bucket", func(t *testing.T) {
		err := c.DeleteBucket("testBucket")
		as.Nil(err)
	})

	t.Run("delete bucket", func(t *testing.T) {
		err := c.Backup("./testBackup.data")
		as.Nil(err)
	})

}
