package database

import (
	"bytes"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os/user"
	"strings"
)

// DB BlotDB的管理类
type DB struct {
	filepath string
	bucket   []byte
	conn     *bolt.DB
}

type KeyValueStore interface {
	Close() error
	RemoveBucket(bucketName string) (err error)
	Add(bucketName string, val []byte) (id uint64, err error)
	Select(bucketName string) error
	RemoveID(bucketName string, id []byte) error
	RemoveVal(bucketName string, val []byte) (err error)
	SelectVal(bucketName string, val []byte) (arr []string, err error)
	RemoveValTransaction(bucketName string, val []byte) (count int, err error)
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

// 创建库管理,并生成Bucket
func (dbase *DB) newConn() error {
	if dbase.conn == nil {
		db, openErr := bolt.Open(dbase.filepath, 0600, nil)
		if openErr != nil {
			return openErr
		}
		dbase.conn = db
		err := dbase.conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(dbase.bucket)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Close 关库数据库
func (dbase *DB) Close() error {
	return dbase.conn.Close()
}

// RemoveBucket 移除Bucket
func (dbase *DB) RemoveBucket(bucketName string) (err error) {
	err = dbase.conn.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
	return err
}

// Add 组Bucket增加值
func (dbase *DB) Add(bucketName string, val []byte) (id uint64, err error) {
	err = dbase.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		id, _ = b.NextSequence() // sequence uint64
		bBuf := fmt.Sprintf("%d", id)
		return b.Put([]byte(bBuf), val)
	})
	return
}

// Select 遍历
func (dbase *DB) Select(bucketName string) error {
	err := dbase.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", string(k), v)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// RemoveID 移除指定Bucket中指定ID
func (dbase *DB) RemoveID(bucketName string, id []byte) error {
	err := dbase.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(id)
	})
	return err
}

// RemoveVal 移除指定Bucket中指定Val
func (dbase *DB) RemoveVal(bucketName string, val []byte) (err error) {
	var arrID []string
	arrID = make([]string, 1)
	err = dbase.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, string(v))
			if bytes.Compare(v, val) == 0 {
				arrID = append(arrID, string(k))
			}
		}
		return nil
	})
	err = dbase.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for _, v := range arrID {
			err := b.Delete([]byte(v))
			if err != nil {
				return err
			}
			fmt.Println("Del k:", v)
		}
		return nil
	})
	return err
}

// SelectVal 查找指定值
func (dbase *DB) SelectVal(bucketName string, val []byte) (arr []string, err error) {
	arr = make([]string, 0, 1)
	err = dbase.conn.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Compare(v, val) == 0 {
				arr = append(arr, string(k))
			}
		}
		return nil
	})
	return arr, err
}

// RemoveValTransaction 在事务中，移除指定Bucket中指定Val
func (dbase *DB) RemoveValTransaction(bucketName string, val []byte) (count int, err error) {
	arrID, err1 := dbase.SelectVal(bucketName, val)
	if err1 != nil {
		return 0, err1
	}
	count = len(arrID)
	if count == 0 {
		return count, nil
	}
	tx, err1 := dbase.conn.Begin(true)
	if err1 != nil {
		return count, err1
	}
	b := tx.Bucket([]byte(bucketName))
	for _, v := range arrID {
		if err = b.Delete([]byte(v)); err != nil {
			fmt.Printf("删除ID(%s)失败! 执行回滚. err:%s \r\n", v, err)
			err := tx.Rollback()
			if err != nil {
				return 0, err
			}
			return
		}
		fmt.Println("删除ID(", v, ")成功!")
	}
	err = tx.Commit()
	return
}
