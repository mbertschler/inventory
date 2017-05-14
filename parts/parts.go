package parts

import (
	"errors"
	"inventory/db"
	"inventory/types"
	"log"
	"sort"
	"unicode"

	"fmt"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var (
	ErrDBEmpty = errors.New("Database is empty, please run database init.")

	ids map[uint64]interface{}
)

// Part describes a part
type Part struct {
	Id        uint64
	Name      string
	Type      uint64
	Value     string
	Location  string
	Datasheet string
	Stock     uint64
}

func (p Part) String() string {
	t := types.ByID(p.Type)
	return fmt.Sprintf(`%s
	Type:      %s
	Value:     %s
	Location:  %s
	Datasheet: %s
	Stock:     %d`,
		p.Name,
		t.Name,
		p.Value,
		p.Location,
		p.Datasheet,
		p.Stock,
	)
}

// Encode returns a binary representation of a part
func (t *Part) Encode() []byte {
	b, err := msgpack.Marshal(&t)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// ID returns the id of a part
func (t *Part) ID() uint64 {
	return t.Id
}

// SetID sets the id of a part
func (t *Part) SetID(id uint64) {
	t.Id = id
}

// Add adds a new type of part to the list
func Add(partName, partType, partValue, partLocation, partDatasheet string) {
	t := types.ByName(partType)
	p := Part{
		Name:      partName,
		Type:      t.Id,
		Value:     partValue,
		Location:  partLocation,
		Datasheet: partDatasheet,
	}
	err := db.Create(db.PartsBucket, &p)

	if err != nil {
		log.Fatal(err)
	}
}

// Rm removes a part from the list
func Rm(name string) {
	p := ByName(name)

	err := db.Delete(db.PartsBucket, &p)

	if err != nil {
		log.Fatal(err)
	}

}

// ByName returns a part by its name
func ByName(name string) Part {
	for _, e := range List() {
		if e.Name == name {
			return e
		}
	}
	return Part{}
}

// List returns a sorted slice of parts
func List() []Part {
	var s []Part
	m := db.GetAll(db.PartsBucket)
	for _, e := range m {
		t := Part{}
		err := msgpack.Unmarshal(e, &t)
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, t)
	}

	sort.Sort(PartSorter(s))

	return s
}

type PartSorter []Part

func (t PartSorter) Len() int {
	return len(t)
}
func (t PartSorter) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t PartSorter) Less(i, j int) bool {
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
