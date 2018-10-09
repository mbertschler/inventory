package main

import (
	"log"
	"testing"

	"git.exahome.net/tools/inventory/parts"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestAdd(t *testing.T) {
	err := parts.Reset()
	if err != nil {
		t.Error(err)
	}
	_, err = parts.Add("Part 1")
	if err != nil {
		t.Error(err)
	}

	_, err = parts.Add("Part 2")
	if err != nil {
		t.Error(err)
	}

	all, err := parts.All()
	if err != nil {
		t.Error(err)
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

func TestEdit(t *testing.T) {
	p1, err := parts.Add("Part 1")
	if err != nil {
		t.Error(err)
	}
	log.Printf("p1 %#v \n", p1)

	all, err := parts.All()
	if err != nil {
		t.Error(err)
	}
	t.Log(all)
}
