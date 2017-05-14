package parts

import (
	"errors"
	"math/rand"
	"time"
	"unicode"
)

var (
	ErrDBEmpty = errors.New("Database is empty, please run database init.")

	ids map[uint64]interface{}
)

func newID() uint64 {
	rand.Seed(time.Now().UnixNano())

	for {
		id := rand.Uint64()
		_, ok := ids[id]
		if !ok {
			return id
		}
	}
}

// Part describes a part
type Part struct {
	ID        uint64
	Name      string
	Type      string
	Value     string
	Location  string
	Datasheet string
}

type PartsByName []*Part

func (p PartsByName) Len() int {
	return len(p)
}
func (p PartsByName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PartsByName) Less(i, j int) bool {
	iRunes := []rune(p[i].Name)
	jRunes := []rune(p[j].Name)

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
