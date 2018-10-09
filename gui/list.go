package gui

import (
	"encoding/json"
	"log"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

func init() {
	// setup guiapi function
	guiapi.DefaultHandler.Functions["newPart"] = newPart
}

func newPart(args json.RawMessage) (*guiapi.Result, error) {
	type input struct {
		Name string
	}
	var in input
	err := json.Unmarshal(args, &in)
	if err != nil {
		return nil, err
	}
	err = parts.Add(in.Name)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", listBlock())
}

func listBlock() html.Block {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	list := html.Blocks{}
	for _, p := range parts {
		list.Add(html.Div(
			html.Class("item"), html.Text(p.Name)),
		)
	}
	return html.Div(html.Class("ui list"),
		list,
		newPartForm,
	)
}

var newPartForm = html.Div(html.Class("ui fluid action input"),
	html.Input(html.Type("text").Attr("placeholder", "Add").Name("Name").Class("ga-new-part")),
	html.Div(html.Class("ui button").Attr("onclick", "sendForm('newPart','.ga-new-part')"), html.Text("Add")),
)
