package samla

import (
	"encoding/json"
	"log"
)

type store interface {
	Get(id string) (box, bool)
	Set(id string, b box) error
	Delete(id string) bool
}

func encode(in interface{}) string {
	out, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func decode(in string, dest interface{}) {
	err := json.Unmarshal([]byte(in), dest)
	if err != nil {
		log.Fatal(err)
	}
}
