package samla

import (
	"log"

	"github.com/etcd-io/bbolt"
)

var boxes = []byte("b")

type embedStore struct {
	db *bbolt.DB
}

func newEmbedStore(file string) (*embedStore, error) {
	db, err := bbolt.Open(file, 0644, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boxes)
		return err
	})
	s := embedStore{
		db: db,
	}
	return &s, nil
}

func (s *embedStore) Get(id string) (box, bool) {
	var b box
	var val []byte
	s.db.View(func(tx *bbolt.Tx) error {
		val = tx.Bucket(boxes).Get([]byte(id))
		return nil
	})
	if val == nil {
		return b, false
	}
	decode(string(val), &b)
	// log.Println("embedStore Get", id, b)
	return b, true
}

func (s *embedStore) Set(id string, b box) error {
	// log.Println("embedStore Set", id, b)
	return s.db.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket(boxes).
			Put([]byte(id), []byte(encode(b)))
		return err
	})
}

func (s *embedStore) Delete(id string) bool {
	// log.Println("embedStore Delete", id)
	err := s.db.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket(boxes).Delete([]byte(id))
		return err
	})
	if err == nil {
		return true
	}
	log.Println("Delete error:", err)
	return false
}
