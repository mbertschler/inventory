package main

import (
	"log"
	"testing"

	"git.exahome.net/tools/inventory/parts"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestDB(t *testing.T) {
	err := parts.Add("Part 1")
	if err != nil {
		log.Println(err)
	}
	err = parts.Add("Part 2")
	if err != nil {
		log.Println(err)
	}
	all, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	for i, p := range all {
		if i == 0 && p.Name != "Part 1" {
			t.Error("wrong data at ", i)
		}
		if i == 1 && p.Name != "Part 2" {
			t.Error("wrong data at ", i)
		}
	}
}
