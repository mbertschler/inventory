package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["viewPart"] = viewPartAction
	guiapi.DefaultHandler.Functions["editPart"] = editPartAction
	guiapi.DefaultHandler.Functions["savePart"] = savePartAction
	guiapi.DefaultHandler.Functions["deletePart"] = deletePartAction
}

func partPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/part/")
	part, err := parts.ByID(id)
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(viewPartBlock(part))
	err = html.Render(page, w)
	if err != nil {
		log.Println(err)
	}
}

func editPartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", editPartBlock(part))
}

func editPartBlock(p *parts.Part) html.Block {
	cancelAction := fmt.Sprintf("guiapi('viewPart', '%s')", p.ID())
	saveAction := "sendForm('savePart', '.ga-edit-part')"
	return html.Div(nil,
		html.Div(nil,
			html.Button(html.Class("ui button").
				Attr("onclick", cancelAction),
				html.Text("Cancel"),
			),
			html.Button(html.Class("ui green button").
				Attr("onclick", saveAction),
				html.Text("Save"),
			),
		),
		html.Div(html.Class("ui form"),
			html.Input(html.Type("hidden").Name("ID").Value(p.ID()).Class("ga-edit-part")),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Name")),
				html.Input(html.Type("Text").Name("Name").Value(p.Name).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Code")),
				html.Input(html.Type("Text").Name("Code").Value(p.Code).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Location")),
				html.Input(html.Type("Text").Name("Location").Value(p.Location).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Parent")),
				html.Input(html.Type("Text").Name("Parent").Value(p.Parent).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Supplier")),
				html.Input(html.Type("Text").Name("Supplier").Value(p.Supplier).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Price (in cents)")),
				html.Input(html.Type("Number").Name("Price").Value(p.Price).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Delivery (in cents)")),
				html.Input(html.Type("Number").Name("Delivery").Value(p.Delivery).Class("ga-edit-part")),
			),
		),
	)
}

func savePartAction(args json.RawMessage) (*guiapi.Result, error) {
	type input struct {
		ID       string
		Name     string
		Code     string
		Location string
		Parent   string
		Supplier string
		Price    string
		Delivery string
	}
	var in input
	err := json.Unmarshal(args, &in)
	if err != nil {
		return nil, err
	}
	p, err := parts.ByID(in.ID)
	if err != nil {
		return nil, err
	}
	p.Name = in.Name
	p.Code = in.Code
	p.Location = in.Location
	p.Parent = in.Parent
	p.Supplier = in.Supplier
	price, err := strconv.Atoi(in.Price)
	if err != nil {
		return nil, err
	}
	p.Price = price
	delivery, err := strconv.Atoi(in.Delivery)
	if err != nil {
		return nil, err
	}
	p.Delivery = delivery
	err = parts.Store(p)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", viewPartBlock(p))
}

func deletePartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	err = parts.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Redirect("/")
}

func viewPartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", viewPartBlock(part))
}

func viewPartBlock(p *parts.Part) html.Block {
	editAction := fmt.Sprintf("guiapi('editPart', '%s')", p.ID())
	deleteAction := fmt.Sprintf("guiapi('deletePart', '%s')", p.ID())

	var rows html.Blocks
	r := func(k, v string) html.Block {
		return html.Elem("tr", nil,
			html.Elem("td", nil, html.Text(k)),
			html.Elem("td", nil, html.Text(v)),
		)
	}
	rows.Add(r("Code", p.Code))
	rows.Add(r("Location", p.Location))
	rows.Add(r("Parent", p.Parent))
	rows.Add(r("Supplier", p.Supplier))
	rows.Add(r("Price", fmt.Sprintf("%.2f€", float64(p.Price)/100)))
	rows.Add(r("Delivery", fmt.Sprintf("%.2f€", float64(p.Delivery)/100)))

	return html.Div(nil,
		html.Div(nil,
			html.A(html.Href("/"),
				html.Button(html.Class("ui button"),
					html.Text("< List"),
				),
			),
			html.Button(html.Class("ui button").
				Attr("onclick", editAction),
				html.Text("Edit"),
			),
			html.Button(html.Class("ui red button").
				Attr("onclick", deleteAction),
				html.Text("Delete"),
			),
		),
		html.H1(nil, html.Text(p.Name)),
		html.Elem("table", html.Class("ui celled table"),
			html.Elem("tbody", nil,
				rows,
			),
		),
	)
}
