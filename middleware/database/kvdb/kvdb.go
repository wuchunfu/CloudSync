package kvdb

import (
	"fmt"
	"github.com/wuchunfu/CloudSync/utils/filex"
	bolt "go.etcd.io/bbolt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

// DBType BlotDB的管理类
type DBType struct {
	dbFilePath string
	bucketName []byte
	conn       *bolt.DB
}

type KVInterface interface {
	// Set 设置值 kv 键值对
	Set(kv map[string][]byte) error
	// GetByKey 根据键名获取各自的值, key 键名
	GetByKey(key string) ([]byte, error)
	// GetByKeys 根据键名数组获取各自的值, keys 键名数组
	GetByKeys(keys []string) (map[string][]byte, error)
	// GetAll 获取全部
	GetAll() (map[string][]byte, error)
	// DeleteByKeys 删除键值
	DeleteByKeys(keys []string) error
	// DeleteByKeysTransaction 在事务中，删除键值
	DeleteByKeysTransaction(keys []string) error
	// DeleteBucket 移除Bucket
	DeleteBucket(bucketName string) (err error)
	// Backup 备份数据库文件
	Backup(backupDbFilePath string) error
	// Close 关库数据库
	Close()
}

// NewKvDb 初始化实例
func NewKvDb(dbFilePath string, bucketName string) KVInterface {
	if dbFilePath == "" {
		return &DBType{}
	}
	if strings.HasPrefix(dbFilePath, "~") {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		dbFilePath = u.HomeDir + dbFilePath[1:]
	}
	dirName := path.Dir(dbFilePath)
	if !filex.FilePathExists(dirName) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Errorf("create dir(%s) error: %s", dirName, err)
			return &DBType{}
		}
	}
	return &DBType{
		dbFilePath: dbFilePath,
		bucketName: []byte(bucketName),
	}
}

// newConn 创建库管理, 并生成Bucket
func (dbase *DBType) newConn() error {
	if dbase.conn == nil {
		db, openErr := bolt.Open(dbase.dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if openErr != nil {
			return fmt.Errorf("open db error: %s", openErr)
		}
		err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(dbase.bucketName)
			if err != nil {
				return fmt.Errorf("create bucket error: %s", err)
			}
			return nil
		})
		dbase.conn = db
		return err
	}
	return nil
}

// Set 设置值 kv 键值对
func (dbase *DBType) Set(kv map[string][]byte) error {
	connErr := dbase.newConn()
	if connErr != nil {
		return connErr
	}

	return dbase.conn.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbase.bucketName)
		var err error
		for k, v := range kv {
			err = bucket.Put([]byte(k), v)
			if err != nil {
				return err
			}
		}
		return err
	})
}

// GetByKey 根据键名获取各自的值, key 键名
func (dbase *DBType) GetByKey(key string) ([]byte, error) {
	connErr := dbase.newConn()
	if connErr != nil {
		return nil, connErr
	}

	var value []byte
	err := dbase.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbase.bucketName)
		result := bucket.Get([]byte(key))
		tmp := make([]byte, len(result))
		copy(tmp, result)
		value = tmp
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetByKey error: %s", err)
	}
	return value, nil
}

// GetByKeys 根据键名数组获取各自的值, keys 键名数组
func (dbase *DBType) GetByKeys(keys []string) (map[string][]byte, error) {
	connErr := dbase.newConn()
	if connErr != nil {
		return nil, connErr
	}

	values := make(map[string][]byte)
	err := dbase.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbase.bucketName)
		for _, k := range keys {
			result := bucket.Get([]byte(k))
			if result == nil {
				continue
			}
			tmp := make([]byte, len(result))
			copy(tmp, result)
			values[k] = tmp
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetByKeys error: %s", err)
	}
	return values, nil
}

// GetAll 获取全部
func (dbase *DBType) GetAll() (map[string][]byte, error) {
	connErr := dbase.newConn()
	if connErr != nil {
		return nil, connErr
	}

	values := make(map[string][]byte)
	err := dbase.conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbase.bucketName)
		err := bucket.ForEach(func(k, v []byte) error {
			tmpK := make([]byte, len(k))
			copy(tmpK, k)
			tmpV := make([]byte, len(v))
			copy(tmpV, v)
			values[string(tmpK)] = tmpV
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetAll error: %s", err)
	}
	return values, nil
}

// DeleteByKeys 删除键值
func (dbase *DBType) DeleteByKeys(keys []string) error {
	connErr := dbase.newConn()
	if connErr != nil {
		return connErr
	}

	return dbase.conn.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbase.bucketName)
		var err error
		for _, key := range keys {
			err = bucket.Delete([]byte(key))
			if err != nil {
				return fmt.Errorf("DeleteByKeys error: %s", err)
			}
		}
		return err
	})
}

// DeleteByKeysTransaction 在事务中，删除键值
func (dbase *DBType) DeleteByKeysTransaction(keys []string) error {
	connErr := dbase.newConn()
	if connErr != nil {
		return connErr
	}

	tx, txErr := dbase.conn.Begin(true)
	if txErr != nil {
		return txErr
	}
	bucket := tx.Bucket(dbase.bucketName)
	for _, key := range keys {
		err := bucket.Delete([]byte(key))
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return fmt.Errorf("DeleteByKeysTransaction error: %s", err)
		}
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteBucket 删除Bucket
func (dbase *DBType) DeleteBucket(bucketName string) error {
	connErr := dbase.newConn()
	if connErr != nil {
		return connErr
	}

	err := dbase.conn.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
	if err != nil {
		return fmt.Errorf("DeleteBucket error: %s", err)
	}
	return err
}

// Backup 备份数据库文件
func (dbase *DBType) Backup(backupDbFilePath string) error {
	db, err := bolt.Open(dbase.dbFilePath, 0600, nil)
	if err != nil {
		return fmt.Errorf("open db error: %s", err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile(backupDbFilePath, 0644)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("backup error: %s", err)
	}
	return db.Close()
}

// Close 关库数据库
func (dbase *DBType) Close() {
	func(conn *bolt.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(dbase.conn)
}
