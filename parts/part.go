package parts

import (
	"log"

	"git.exahome.net/tools/inventory/lib/samla"
)

var db *samla.DB
var allPartsID string

func init() {
	db = samla.NewMemoryDB()
	err := db.RegisterTypes(partType, allPartsType)
	if err != nil {
		log.Fatal(err)
	}
	err = Reset()
	if err != nil {
		log.Fatal(err)
	}
}

// Part is the central item in the inventory.
type Part struct {
	samla.Reference
	Name string
}

var partType = samla.Type{
	Name: "Part",
	New: func() interface{} {
		return &Part{}
	},
	Reference: func(obj interface{}) *samla.Reference {
		return &obj.(*Part).Reference
	},
	Fields: samla.Fields{
		"Name": {
			Type: samla.String,
			Get: func(obj interface{}) interface{} {
				return obj.(*Part).Name
			},
			Set: func(obj interface{}, val interface{}) {
				(obj.(*Part)).Name = val.(string)
			},
		},
	},
}

// AllParts holds all Parts
type AllParts struct {
	samla.Reference
	Parts []*Part
}

var allPartsType = samla.Type{
	Name: "AllParts",
	New: func() interface{} {
		return &AllParts{}
	},
	Reference: func(obj interface{}) *samla.Reference {
		return &obj.(*AllParts).Reference
	},
	Links: samla.Links{
		"Parts": {
			Kind: samla.PointerSlice,
			Type: "Part",
			Get: func(obj interface{}) interface{} {
				return obj.(*AllParts).Parts
			},
			Set: func(obj, val interface{}) {
				(obj.(*AllParts)).Parts = val.([]*Part)
			},
			ToSlice: func(in interface{}) []interface{} {
				arr := in.([]*Part)
				out := make([]interface{}, len(arr))
				for i := range arr {
					out[i] = arr[i]
				}
				return out
			},
			FromSlice: func(in []interface{}) interface{} {
				out := make([]*Part, len(in))
				for i := range in {
					out[i] = in[i].(*Part)
				}
				return out
			},
		},
	},
}

// All returns all Parts from the database.
func All() ([]*Part, error) {
	var all AllParts
	err := db.LoadIDWith(&all, allPartsID, "Parts")
	return all.Parts, err
}

// Add adds a new Part to the database.
func Add(name string) (*Part, error) {
	var all AllParts
	err := db.LoadIDWith(&all, allPartsID, "Parts")
	if err != nil {
		log.Println(err)
	}

	newPart := Part{Name: name}
	all.Parts = append(all.Parts, &newPart)

	_, err = db.StoreAsWith(&all, "AllParts", "Parts")
	return &newPart, err
}

// Reset clears the list that contains all parts.
func Reset() error {
	// TODO actually reset the samla DB
	all := AllParts{}
	info, err := db.StoreAs(&all, "AllParts")
	if err != nil {
		log.Fatal(err)
	}
	allPartsID = info.ID
	return err
}

// ByID loads a Part by an ID.
func ByID(id string) (*Part, error) {
	var p Part
	err := db.LoadID(&p, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &p, nil
}

// DeleteByID deletes a Part that is already in the DB.
func DeleteByID(id string) error {
	return db.DeleteID(id)
}

// Store saves a Part that is already in the DB.
func Store(p *Part) error {
	info, err := db.StoreAs(&p, "Part")
	log.Println("Store:", info)
	return err
}
