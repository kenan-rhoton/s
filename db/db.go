package db

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

type DB struct {
	name string
	bolt *bolt.DB
}

func UseDatabase(name string) (*DB, error) {
	db := &DB{}
	db.name = name
	err := db.Create()
	return db, err
}

func (db *DB) Create() error {
	var err error
	db.bolt, err = bolt.Open(db.name+".db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte("ID"))
		if e != nil {
			return fmt.Errorf("create bucket: %s", e)
		}
		_, e = tx.CreateBucketIfNotExists([]byte("Messages"))
		if e != nil {
			return fmt.Errorf("create bucket: %s", e)
		}
		return nil
	})

	return err
}

func (db *DB) Close() error {
	return db.bolt.Close()
}

func (db *DB) Reset() error {
	err := db.Close()
	if err != nil {
		return err
	}
	err = os.Remove(db.name + ".db")
	if err != nil {
		return err
	}
	db.Create()
	return err

}

func (db *DB) SaveID(user string, key []byte) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("ID"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte(user), key)
		return err
	})
}

func (db *DB) GetKey(user string) ([]byte, error) {
	var res []byte
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ID"))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		key := b.Get([]byte(user))
		if key == nil {
			return fmt.Errorf("invalid target")
		}
		res = make([]byte, len(key), (cap(key)+1)*2)
		copy(res, key)
		return nil
	})
	return res, err
}

func (db *DB) SaveMessage(user string, msg []byte) error {
	_, err := db.GetKey(user)
	if err != nil {
		return err
	}
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Messages"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		u, err := b.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		num, _ := u.NextSequence()
		err = u.Put(itob(num), msg)
		return err
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func (db *DB) GetMessages(user string) ([][]byte, error) {
	res := make([][]byte, 0)
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Messages"))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		u := b.Bucket([]byte(user))
		if u == nil {
			return fmt.Errorf("bucket does not exist")
		}
		c := u.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			temp := make([]byte, len(v), (cap(v)+1)*2)
			copy(temp, v)
			res = append(res, temp)
		}
		return nil
	})
	return res, err
}
