package samla

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Samla interface {
	StoreAs(obj interface{}, typ string) error
	LoadID(obj interface{}, id string) error
	DeleteID(id string) error
}

type FieldType int8

const (
	String FieldType = iota + 1
	Float
	Int
	Bool
)

type Field struct {
	Type FieldType
	Get  func(obj interface{}) interface{}
	Set  func(obj interface{}, val interface{})
}

type Fields map[string]Field

type LinkType int8

const (
	Pointer LinkType = iota + 1
	PointerSlice
)

type Link struct {
	Kind      LinkType
	Type      string
	Get       func(obj interface{}) interface{}
	Set       func(obj interface{}, val interface{})
	ToSlice   func(obj interface{}) []interface{}
	FromSlice func(slice []interface{}) interface{}
}

type Links map[string]Link

type Type struct {
	Name      string
	New       func() interface{}
	Version   int
	Fields    Fields
	Reference func(obj interface{}) *Reference
	Links     Links
}

func NewDB() *DB {
	return &DB{
		types: map[string]Type{},
		boxes: map[string]box{},
	}
}

type DB struct {
	types map[string]Type
	boxes map[string]box
}

type box struct {
	fields map[string]boxField
	links  map[string][]string
}

type boxField struct {
	value interface{}
}

type StoreInfo struct {
	ID string
}

type Reference struct {
	id  string
	obj interface{}
}

func (r *Reference) ID() string {
	return r.id
}

var shittyID int

func newShittyID(typ string) string {
	shittyID++
	return typ + "/" + strconv.Itoa(shittyID)
}

func getShittyType(id string) string {
	return id[:strings.Index(id, "/")]
}

func (d *DB) StoreAsWith(in interface{}, typ string, links ...string) (StoreInfo, error) {
	info, err := d.StoreAs(in, typ)
	if err != nil {
		return info, err
	}
	t := d.types[typ]
	b := d.boxes[info.ID]
	b.links = map[string][]string{}
	for _, e := range links {
		l := t.Links[e]
		switch l.Kind {
		case Pointer:
			o := l.Get(in)
			inf, _ := d.StoreAs(o, l.Type)
			b.links[e] = []string{inf.ID}
		case PointerSlice:
			o := l.Get(in)
			arr := l.ToSlice(o)
			links := make([]string, len(arr))
			for i, el := range arr {
				info, err := d.StoreAs(el, l.Type)
				if err != nil {
					log.Println(err)
				}
				links[i] = info.ID
			}
			b.links[e] = links
		}
	}
	d.boxes[info.ID] = b
	return info, nil
}

func (d *DB) StoreAs(in interface{}, typ string) (StoreInfo, error) {
	t := d.types[typ]
	b := box{
		fields: map[string]boxField{},
	}
	ref := t.Reference(in)
	id := ref.id
	if id == "" {
		id = newShittyID(typ)
	}
	for name, e := range t.Fields {
		switch e.Type {
		case String:
			s := e.Get(in)
			f := boxField{value: s}
			b.fields[name] = f
		}
	}
	d.boxes[id] = b
	ref.id = id
	ref.obj = in
	// log.Printf("stored %#v %+v\n", id, b)
	return StoreInfo{ID: id}, nil
}

func (d *DB) LoadID(out interface{}, id string) error {
	typ := getShittyType(id)
	t := d.types[typ]
	b, ok := d.boxes[id]
	if !ok {
		return fmt.Errorf("samla: id %s not found", id)
	}
	ref := t.Reference(out)
	ref.id = id
	ref.obj = out
	for name, e := range t.Fields {
		switch e.Type {
		case String:
			e.Set(out, b.fields[name].value)
		}
	}
	return nil
}

func (d *DB) LoadIDWith(out interface{}, id string, links ...string) error {
	err := d.LoadID(out, id)
	if err != nil {
		return err
	}
	typ := getShittyType(id)
	t := d.types[typ]
	b := d.boxes[id]
	for _, e := range links {
		l := t.Links[e]
		switch l.Kind {
		case Pointer:
			ids := b.links[e]
			if len(ids) > 0 {
				linkType := d.types[l.Type]
				o := linkType.New()
				d.LoadID(o, ids[0])
				l.Set(out, o)
			}
		case PointerSlice:
			ids := b.links[e]
			arr := make([]interface{}, len(ids))
			linkType := d.types[l.Type]
			for i := range ids {
				o := linkType.New()
				err := d.LoadID(o, ids[i])
				if err != nil {
					log.Println(err)
				}
				arr[i] = o
			}
			l.Set(out, l.FromSlice(arr))
		}
	}
	return nil
}

func (d *DB) DeleteID(id string) error {
	delete(d.boxes, id)
	return nil
}

func (d *DB) RegisterTypes(types ...Type) error {
	for _, t := range types {
		d.types[t.Name] = t
	}
	return nil
}
