package kvdb

import (
	"encoding/binary"
	"errors"
	"os/user"
	"strings"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

type KV struct {
	// Key
	Key string
	// Value
	Val string
	// TTL(存活时间), 如果设置 ttl = 0 or Persistent, 这个key就会永久不删除
	TTL time.Duration
}

type DB struct {
	filepath string
	bucket   []byte
	bOnce    sync.Once
	conn     *bolt.DB
}

type KeyValueStore interface {
	GetBytes(key string) ([]byte, error)
	Get(key string) (string, error)
	SetBytes(key string, val []byte, ttl time.Duration) error
	Set(key, val string, ttl time.Duration) error
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) error
	Del(key string) error
	Range() ([]*KV, error)
}

func NewDB(filepath string) KeyValueStore {
	if strings.HasPrefix(filepath, "~") {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		filepath = u.HomeDir + filepath[1:]
	}
	return &DB{
		filepath: filepath,
		bucket:   []byte("kv-bucket"),
	}
}

func (dbase *DB) newConn() error {
	if dbase.conn == nil {
		db, err := bolt.Open(dbase.filepath, 0600, nil)
		if err != nil {
			return err
		}
		dbase.conn = db
	}
	return nil
}

func (dbase *DB) GetBytes(key string) ([]byte, error) {
	ttl, result, err := dbase.getWithExpire(key)
	if err != nil {
		return nil, err
	} else if ttl < 0 {
		return nil, nil
	}
	return result, nil
}

func (dbase *DB) Get(key string) (string, error) {
	ttl, result, err := dbase.getWithExpire(key)
	if err != nil {
		return "", err
	} else if ttl < 0 {
		return "", nil
	}
	return string(result), nil
}

func (dbase *DB) SetBytes(key string, val []byte, ttl time.Duration) error {
	if err := dbase.newConn(); err != nil {
		return err
	}
	return dbase.conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(dbase.bucket)
		if err != nil {
			return err
		}
		buf := make([]byte, 8+len(val))
		binary.PutVarint(buf[:8], toMillisecond(ttl))
		copy(buf[8:], val)
		return b.Put([]byte(key), buf)
	})
}

func (dbase *DB) Set(key, val string, ttl time.Duration) error {
	if err := dbase.newConn(); err != nil {
		return err
	}
	return dbase.conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(dbase.bucket)
		if err != nil {
			return err
		}
		buf := make([]byte, 8+len(val))
		binary.PutVarint(buf[:8], toMillisecond(ttl))
		copy(buf[8:], val)
		return b.Put([]byte(key), buf)
	})
}

func (dbase *DB) TTL(key string) (time.Duration, error) {
	ttl, _, err := dbase.getWithExpire(key)
	if err != nil {
		return -1, err
	} else if ttl < -1 {
		return -1, nil
	}
	return time.Duration(ttl) * time.Millisecond, nil
}

func (dbase *DB) Expire(key string, ttl time.Duration) error {
	_, result, err := dbase.getWithExpire(key)
	if err != nil {
		return err
	} else if ttl < -1 {
		return errors.New("key expired")
	}
	return dbase.Set(key, string(result), ttl)
}

func (dbase *DB) Del(key string) error {
	if err := dbase.newConn(); err != nil {
		return err
	}
	return dbase.conn.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(dbase.bucket)
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

func (dbase *DB) Range() ([]*KV, error) {
	if err := dbase.newConn(); err != nil {
		return nil, err
	}
	var kvs []*KV
	if err := dbase.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbase.bucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			expiredAt, err := binaryInt(v[:8])
			if err != nil {
				return err
			}
			ttl := expiredAt - int(time.Now().UnixNano()/int64(1000000))
			if ttl < 0 {
				// 过期了, 删除
				err := dbase.conn.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists(dbase.bucket)
					if err != nil {
						return err
					}
					return b.Delete(k)
				})
				if err != nil {
					return err
				}
				return nil
			}
			kvs = append(kvs, &KV{
				Key: string(k),
				Val: string(v[8:]),
				TTL: time.Duration(ttl) * time.Millisecond,
			})
			return nil
		})
	}); err != nil {
		return nil, err
	}
	return kvs, nil
}

func (dbase *DB) getOriginData(key string) ([]byte, error) {
	if err := dbase.newConn(); err != nil {
		return nil, err
	}
	var result []byte
	if err := dbase.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbase.bucket)
		if b == nil {
			return nil
		}
		result = b.Get([]byte(key))
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (dbase *DB) getWithExpire(key string) (int, []byte, error) {
	result, err := dbase.getOriginData(key)
	if err != nil {
		return -1, nil, nil
	} else if result == nil {
		return -1, nil, nil
	}
	expiredAt, err := binaryInt(result[:8])
	if err != nil {
		return -1, nil, err
	}
	ttl := expiredAt - int(time.Now().UnixNano()/int64(1000000))
	if ttl < 0 {
		// 过期了, 删除
		err := dbase.conn.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(dbase.bucket)
			if err != nil {
				return err
			}
			return b.Delete([]byte(key))
		})
		if err != nil {
			return -1, nil, err
		}
		return -1, nil, err
	}
	return ttl, result[8:], nil
}

func toMillisecond(ttl time.Duration) int64 {
	return time.Now().Add(ttl).UnixNano() / int64(1000000)
}

func binaryInt(buf []byte) (int, error) {
	x, n := binary.Varint(buf)
	if n == 0 {
		return 0, errors.New("buf too small")
	} else if n < 0 {
		return 0, errors.New("value larger than 64 bits (overflow) and -n is the number of bytes read")
	}
	return int(x), nil
}
