package types

import (
	"inventory/db"
	"log"
	"sort"
	"unicode"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

// Type describes a category of parts
type Type struct {
	Id   uint64
	Name string
}

func (t Type) String() string {
	return t.Name
}

// Encode returns a binary representation of a type
func (t *Type) Encode() []byte {
	b, err := msgpack.Marshal(&t)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// ID returns the id of a type
func (t *Type) ID() uint64 {
	return t.Id
}

// SetID sets the id of a type
func (t *Type) SetID(id uint64) {
	t.Id = id
}

// Add adds a new type of part to the list
func Add(name string) {
	t := Type{Name: name}
	err := db.Create(db.TypesBucket, &t)

	if err != nil {
		log.Fatal(err)
	}
}

// Rm removes a type of part from the list
func Rm(name string) {
	t := ByName(name)

	err := db.Delete(db.TypesBucket, &t)

	if err != nil {
		log.Fatal(err)
	}

}

// ByName returns a type by its name
func ByName(name string) Type {
	for _, e := range List() {
		if e.Name == name {
			return e
		}
	}
	return Type{}
}

// ByID returns a type by its name
func ByID(id uint64) Type {
	b := db.Get(db.TypesBucket, id)
	t := Type{}
	err := msgpack.Unmarshal(b, &t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// List returns a sorted slice of part types
func List() []Type {
	var s []Type
	m := db.GetAll(db.TypesBucket)
	for _, e := range m {
		t := Type{}
		err := msgpack.Unmarshal(e, &t)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, t)
	}

	sort.Sort(TypeSorter(s))

	return s
}

type TypeSorter []Type

func (t TypeSorter) Len() int {
	return len(t)
}
func (t TypeSorter) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TypeSorter) Less(i, j int) bool {
	iRunes := []rune(t[i].Name)
	jRunes := []rune(t[j].Name)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}
	return false
}
