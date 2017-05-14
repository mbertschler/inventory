package db

import (
	"encoding/binary"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbPath = "/tmp/inventory.db"

	PartsBucket = "Parts"
	TypesBucket = "Types"
)

var (
	db *bolt.DB

	ErrBucketNotExist = errors.New("bucket does not exist")
)

type InventoryStorer interface {
	Encode() []byte
	ID() uint64
	SetID(id uint64)
}

func Open() error {
	var err error
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		var err error

		_, err = tx.CreateBucketIfNotExists([]byte(TypesBucket))
		if err != nil {
			log.Fatal(err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(PartsBucket))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return err
}

func Close() error {
	return db.Close()
}

func GetAll(bucket string) map[uint64][]byte {
	m := make(map[uint64][]byte)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotExist
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			m[btoi(k)] = v
		}

		return nil
	})

	return m
}

func Get(bucket string, key uint64) []byte {
	var v []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotExist
		}

		v = b.Get(itob(key))

		return nil
	})

	return v
}

func Create(bucket string, value InventoryStorer) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotExist
		}

		id, _ := b.NextSequence()
		value.SetID(id)

		return b.Put(itob(value.ID()), value.Encode())
	})
}

func Update(bucket string, value InventoryStorer) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotExist
		}

		return b.Put(itob(value.ID()), value.Encode())
	})
}

func Delete(bucket string, value InventoryStorer) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotExist
		}

		return b.Delete(itob(value.ID()))
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// btoi returns a uint64 big endianrepresentation of b.
func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
