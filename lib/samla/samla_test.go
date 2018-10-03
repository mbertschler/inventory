package samla

import (
	"testing"
)

// List holds lines
type List struct {
	Name  string
	Main  *Line
	Lines []*Line
}

var listType = Type{
	Name: "List",
	New: func() interface{} {
		return &List{}
	},
	Fields: Fields{
		"Name": {
			Type: String,
			Get: func(obj interface{}) interface{} {
				return obj.(*List).Name
			},
			Set: func(obj interface{}, val interface{}) {
				(obj.(*List)).Name = val.(string)
			},
		},
	},
}

// Line is a line of test
type Line struct {
	Parent *List `json:"-"` // JSON cyclic pointer
	Text   string
}

var lineType = Type{
	Name: "Line",
	New: func() interface{} {
		return &Line{}
	},
	Fields: Fields{
		"Text": {
			Type: String,
			Get: func(obj interface{}) interface{} {
				return obj.(*Line).Text
			},
			Set: func(obj interface{}, val interface{}) {
				(obj.(*Line)).Text = val.(string)
			},
		},
	},
	Links: Links{
		"Parent": {
			Kind: Pointer,
			Type: "List",
			Get: func(obj interface{}) interface{} {
				return obj.(*Line).Parent
			},
			Set: func(obj, val interface{}) {
				(obj.(*Line)).Parent = val.(*List)
			},
		},
	},
}

func TestOne(t *testing.T) {
	db := NewDB()
	err := db.RegisterTypes(lineType)
	if err != nil {
		t.Error(err)
	}
	l := Line{
		Text: "testing",
	}
	info, err := db.StoreAs(&l, "Line")
	if err != nil {
		t.Error(err)
	}
	if info.ID == "" {
		t.Error("StoreAs didn't return an ID")
	}
	l = Line{}
	err = db.LoadID(&l, info.ID)
	if err != nil {
		t.Error(err)
	}
	if l.Text != "testing" {
		t.Errorf("line text not \"testing\" but %#v", l.Text)
	}
}

func TestTwo(t *testing.T) {
	db := NewDB()
	err := db.RegisterTypes(lineType, listType)
	if err != nil {
		t.Error(err)
	}
	list := List{
		Name: "main",
	}
	line := Line{
		Parent: &list,
		Text:   "testing",
	}
	info, err := db.StoreAsWith(&line, "Line", "Parent")
	if err != nil {
		t.Error(err)
	}
	if info.ID == "" {
		t.Error("StoreAs didn't return an ID")
	}
	line = Line{}
	err = db.LoadIDWith(&line, info.ID, "Parent")
	if err != nil {
		t.Error(err)
	}
	if line.Text != "testing" {
		t.Errorf("line text not \"testing\" but %#v", line.Text)
	}
	if line.Parent == nil {
		t.Errorf("line parent is nil")
	} else if line.Parent.Name != "main" {
		t.Errorf("list name not \"main\" but %#v", line.Parent.Name)
	}
}

func TestThree(t *testing.T) {
	db := NewDB()
	err := db.RegisterTypes(lineType, listType)
	if err != nil {
		t.Error(err)
	}
	list := List{
		Name: "main",
	}
	line := Line{
		Parent: &list,
		Text:   "testing",
	}
	info, err := db.StoreAsWith(&line, "Line", "Parent")
	if err != nil {
		t.Error(err)
	}
	if info.ID == "" {
		t.Error("StoreAs didn't return an ID")
	}
	line = Line{}
	err = db.LoadIDWith(&line, info.ID, "Parent")
	if err != nil {
		t.Error(err)
	}
	if line.Text != "testing" {
		t.Errorf("line text not \"testing\" but %#v", line.Text)
	}
	if line.Parent == nil {
		t.Errorf("line parent is nil")
	} else if line.Parent.Name != "main" {
		t.Errorf("list name not \"main\" but %#v", line.Parent.Name)
	}
}
