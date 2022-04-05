package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mbertschler/blocks/html"

	"github.com/mbertschler/inventory/lib/guiapi"
	"github.com/mbertschler/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["viewPart"] = viewPartAction
	guiapi.DefaultHandler.Functions["newPart"] = newPartAction
	guiapi.DefaultHandler.Functions["editPart"] = editPartAction
	guiapi.DefaultHandler.Functions["savePart"] = savePartAction
	guiapi.DefaultHandler.Functions["deletePart"] = deletePartAction
	guiapi.DefaultHandler.Functions["checkout"] = checkoutAction
}

func partPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/part/")
	part, err := parts.ByID(id)
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(viewPartBlock(part))
	err = html.Render(w, page)
	if err != nil {
		log.Println(err)
	}
}

func newPartAction(_ json.RawMessage) (*guiapi.Result, error) {
	return guiapi.Replace("#container", editPartBlock(nil, ""))
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
	return guiapi.Replace("#container", editPartBlock(part, ""))
}

func editPartBlock(p *parts.Part, code string) html.Block {
	isNew := false
	if p == nil {
		isNew = true
		p = &parts.Part{}
	}
	if code != "" {
		p.Code = code
	}
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
			html.Input(html.Type("hidden").Name("New").Value(isNew).Class("ga-edit-part")),
			html.Input(html.Type("hidden").Name("ID").Value(p.ID()).Class("ga-edit-part")),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Code")),
				html.Input(html.Type("Text").Name("Code").Value(p.Code).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Name")),
				html.Input(html.Type("Text").Name("Name").Value(p.Name).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Type (R, C, npn Mosfet)")),
				html.Input(html.Type("Text").Name("Type").Value(p.Type).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Value (format: 5.8k, 10, 5u)")),
				html.Input(html.Type("Text").Name("Value").Value(p.Value).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Size")),
				html.Input(html.Type("Text").Name("Size").Value(p.Size).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Serial")),
				html.Input(html.Type("Text").Name("Serial").Value(p.Serial).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Quantity")),
				html.Input(html.Type("Number").Name("Quantity").Value(p.Quantity).Class("ga-edit-part")),
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
				html.Label(nil, html.Text("Price (in €, format: 12.3456)")),
				html.Input(html.Type("Number").Attr("step", "0.01").Name("Price").Value(float64(p.Price)/10000).Class("ga-edit-part")),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Delivery (in €, format: 12.3456)")),
				html.Input(html.Type("Number").Name("Delivery").Value(float64(p.Delivery)/10000).Class("ga-edit-part")),
			),
		),
	)
}

func savePartAction(args json.RawMessage) (*guiapi.Result, error) {
	type input struct {
		ID       string
		New      string
		Name     string
		Type     string
		Value    string
		Serial   string
		Code     string
		Size     string
		Quantity string
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
	var p *parts.Part
	if in.New == "true" {
		p, err = parts.Create()
	} else {
		p, err = parts.ByID(in.ID)
	}
	if err != nil {
		return nil, err
	}
	p.Name = in.Name
	p.Type = in.Type
	p.Value = in.Value
	p.Code = in.Code
	p.Serial = in.Serial
	p.Size = in.Size
	p.Location = in.Location
	quant, err := strconv.Atoi(in.Quantity)
	if err != nil {
		return nil, err
	}
	p.Quantity = quant
	p.Parent = in.Parent
	p.Supplier = in.Supplier
	price, err := strconv.ParseFloat(in.Price, 64)
	if err != nil {
		return nil, err
	}
	p.Price = int(price * 10000)
	delivery, err := strconv.ParseFloat(in.Delivery, 64)
	if err != nil {
		return nil, err
	}
	p.Delivery = int(delivery * 10000)

	// generate name for some parts
	switch p.Type {
	case "R", "L", "C", "CP":
		p.Name = p.Type + p.Size + "_" + p.Value
	}

	// store part
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
	if part == nil {
		return guiapi.Redirect("/")
	}
	return guiapi.Replace("#container", viewPartBlock(part))
}

func checkoutAction(args json.RawMessage) (*guiapi.Result, error) {
	var data struct {
		ID       string
		Quantity string
	}
	err := json.Unmarshal(args, &data)
	if err != nil {
		return nil, err
	}
	quant, err := strconv.Atoi(data.Quantity)
	if err != nil {
		return nil, err
	}
	part, err := parts.ByID(data.ID)
	if err != nil {
		return nil, err
	}

	part.Quantity -= quant

	err = parts.Store(part)
	if err != nil {
		return nil, err
	}

	if part == nil {
		return guiapi.Redirect("/")
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
	rows.Add(r("Type", p.Type))
	rows.Add(r("Value", p.Value))
	rows.Add(r("Size", p.Size))
	rows.Add(r("Serial", p.Serial))
	rows.Add(r("Quantity", fmt.Sprint(p.Quantity)))
	rows.Add(r("Location", p.Location))
	rows.Add(r("Parent", p.Parent))
	rows.Add(r("Supplier", p.Supplier))
	rows.Add(r("Price", fmt.Sprintf("%.4f€", float64(p.Price)/10000)))
	rows.Add(r("Delivery", fmt.Sprintf("%.2f€", float64(p.Delivery)/10000)))

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
		html.Div(html.Class("ui form"),
			html.Input(html.Type("hidden").Name("ID").Value(p.ID()).Class("ga-checkout")),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Inventory")),
				html.Input(html.Type("Text").Value(p.Quantity).Attr("disabled", true)),
			),
			html.Div(html.Class("field"),
				html.Label(nil, html.Text("Checkout")),
				html.Input(html.Type("Text").Name("Quantity").Value("0").Class("ga-checkout")),
			),
			html.Button(html.Class("ui yellow button").
				Attr("onclick", "sendForm('checkout', '.ga-checkout')"),
				html.Text("Checkout Parts"),
			),
		),
		html.Elem("table", html.Class("ui celled table"),
			html.Elem("tbody", nil,
				rows,
			),
		),
	)
}
