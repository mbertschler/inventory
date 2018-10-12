package samla

import (
	"log"

	"github.com/tidwall/buntdb"
)

type memoryStore struct {
	db *buntdb.DB
}

func newMemoryStore() (*memoryStore, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}
	s := memoryStore{
		db: db,
	}
	return &s, nil
}

func (s *memoryStore) Get(id string) (box, bool) {

	var b box
	var val string
	err := s.db.View(func(tx *buntdb.Tx) error {
		var err error
		val, err = tx.Get(id)
		return err
	})
	if err != nil {
		log.Println("Get error:", err)
		return b, false
	}
	decode(val, &b)
	// log.Println("memoryStore Get", id, b)
	return b, true
}

func (s *memoryStore) Set(id string, b box) error {
	// log.Println("memoryStore Set", id, b)
	return s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(id, encode(b), nil)
		return err
	})
}

func (s *memoryStore) Delete(id string) bool {
	// log.Println("memoryStore Delete", id)
	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})
	if err == nil {
		return true
	}
	log.Println("Delete error:", err)
	return false
}
