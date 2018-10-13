package parts

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/etcd-io/bbolt"
)

var (
	db *bolt.DB

	partsBucket = []byte("p")
)

// Ref stores a da id in the database
type Ref struct {
	id []byte
}

func (r Ref) ID() string {
	return keyToBase64(r.id)
}

// Part is the central item in the inventory.
type Part struct {
	Ref
	Code     string
	Name     string
	Location string
	Parent   string

	Supplier string
	Price    int
	Delivery int

	Values map[string]string
}

// SetupDB opens the database.
func SetupDB(path string) error {
	var err error
	db, err = bolt.Open(path, 0644, nil)
	if err != nil {
		return err
	}
	return initDB()
}

func initDB() error {
	rand.Seed(time.Now().Unix())
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(partsBucket)
		return err
	})
}

// All returns all Parts from the database.
func All() ([]*Part, error) {
	var all []*Part
	err := db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(partsBucket).ForEach(func(k, v []byte) error {
			p, err := unpackPart(k, v)
			if err != nil {
				return err
			}
			all = append(all, p)
			return nil
		})
	})
	return all, err
}

// keyFromID returns an 8-byte big endian representation of i.
func idToKey(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// keyToID returns an uint64 from an 8-byte big endian key.
func keyToID(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// base64 returns the base64 representation of in.
func keyToBase64(in []byte) string {
	return base64.RawURLEncoding.EncodeToString(in)
}

// Base64ToID parses an ID out of a base64 string.
func base64ToKey(in string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(in)
}

func unpackPart(key, value []byte) (*Part, error) {
	var p Part
	err := json.Unmarshal(value, &p)
	if err != nil {
		return nil, err
	}
	p.Ref.id = key
	return &p, nil
}

func packPart(p *Part) (key, value []byte, err error) {
	val, err := json.Marshal(p)
	return p.Ref.id, val, err
}

func randomKey() []byte {
	u := rand.Uint64()
	return idToKey(u)
}

// Add adds a new Part to the database.
func Add(name string) (*Part, error) {
	newPart := Part{Name: name}
	_, v, err := packPart(&newPart)
	k := randomKey()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(partsBucket)
		for {
			if b.Get(k) == nil {
				return b.Put(k, v)
			}
			// new key if it already exists
			k = randomKey()
		}
	})
	newPart.Ref.id = k
	return &newPart, err
}

// Reset clears the list that contains all parts.
func Reset() error {
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(partsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(partsBucket)
		return err
	})
}

// ByID loads a Part by an ID.
func ByID(id string) (*Part, error) {
	k, err := base64ToKey(id)
	if err != nil {
		return nil, err
	}
	var v []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(partsBucket)
		v = b.Get(k)
		return nil
	})
	if err != nil {
		return nil, err
	}
	p, err := unpackPart(k, v)
	return p, err
}

// DeleteByID deletes a Part that is already in the DB.
func DeleteByID(id string) error {
	k, err := base64ToKey(id)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(partsBucket)
		return b.Delete(k)
	})
}

// Store saves a Part that is already in the DB.
func Store(p *Part) error {
	k, v, err := packPart(p)
	if err != nil {
		return err
	}
	if k == nil {
		return errors.New("store: new part has no key yet")
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(partsBucket)
		return b.Put(k, v)
	})
}
