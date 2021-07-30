package kvdb

import (
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	testDB := NewDB("./test.data")
	defer os.Remove("./test.data")

	v, _ := testDB.Get("not-exist")
	t.Log(v)

	testDB.Set("k", "v", 1*time.Minute)

	v, _ = testDB.Get("k")
	t.Log(v)

	ttl, _ := testDB.TTL("k")
	t.Log(int(math.Ceil(ttl.Seconds())))

	time.Sleep(2 * time.Second)

	ttl, _ = testDB.TTL("k")
	t.Log(int(math.Ceil(ttl.Seconds())))

	testDB.Del("k")

	v, _ = testDB.Get("k")
	t.Log(v)
}

func TestNew(t *testing.T) {
	as := assert.New(t)
	defer os.Remove("./test.data")

	os.Remove("./test.data")
	c := NewDB("./test.data")

	t.Run("not found", func(t *testing.T) {
		v, err := c.Get("k")
		as.Equal(nil, err)
		t.Log(v)
		as.Empty(v)
	})

	t.Run("exist set ttl get", func(t *testing.T) {
		as.Nil(c.Set("k", "v", 1*time.Second))

		v, err := c.Get("k")
		as.Nil(err)
		t.Log(v)
		as.Equal("v", v)
	})

	t.Run("expired", func(t *testing.T) {
		as.Nil(c.Set("k", "vf", 1*time.Second))

		time.Sleep(1 * time.Second)

		v, err := c.Get("k")
		as.Equal(nil, err)
		t.Log(v)
		as.Empty(v)
	})

	t.Run("ttl", func(t *testing.T) {
		as.Nil(c.Set("k", "v", 1*time.Second))

		ttl, err := c.TTL("k")
		as.Nil(err)
		t.Log(ttl)
		as.True(ttl <= 1*time.Second && ttl >= 1*time.Second-100*time.Millisecond, ttl)
	})

	t.Run("expire", func(t *testing.T) {
		as.Nil(c.Set("k", "v", 1*time.Second))

		ttl, err := c.TTL("k")
		as.Nil(err)
		t.Log(ttl)
		as.True(ttl <= 1*time.Second && ttl >= 1*time.Second-100*time.Millisecond)

		as.Nil(c.Expire("k", 1*time.Minute))

		ttl, err = c.TTL("k")
		as.Nil(err)
		t.Log(ttl)
		as.True(ttl <= 1*time.Minute && ttl >= 1*time.Minute-100*time.Millisecond)
	})

	t.Run("range", func(t *testing.T) {
		os.Remove("./test")
		c = NewDB("./test")

		for i := 0; i < 10; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, 1*time.Minute), i)
		}

		kvs, err := c.Range()
		as.Nil(err)
		for _, v := range kvs {
			t.Log(v)
			as.Equal(v.Key, v.Val)
		}
		as.Len(kvs, 10)
	})
}
