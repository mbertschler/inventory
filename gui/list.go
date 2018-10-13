package gui

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["newPart"] = newPartAction
}

func newPartAction(args json.RawMessage) (*guiapi.Result, error) {
	type input struct {
		Name string
	}
	var in input
	err := json.Unmarshal(args, &in)
	if err != nil {
		return nil, err
	}
	p, err := parts.Add(in.Name)
	if err != nil {
		return nil, err
	}
	return guiapi.Replace("#container", editPartBlock(p))
}

func listBlock() html.Block {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	list := html.Blocks{}
	for _, p := range parts {
		link := fmt.Sprintf("/part/%s", p.ID())
		list.Add(html.A(html.Href(link),
			html.Div(html.Class("item"),
				html.Text(p.Name)),
		))
	}
	return html.Div(html.Class("ui list"),
		list,
		newPartForm,
	)
}

var newPartForm = html.Div(html.Class("ui fluid action input"),
	html.Input(html.Type("text").Attr("placeholder", "Add").Name("Name").Class("ga-new-part")),
	html.Div(html.Class("ui button").Attr("onclick", "sendForm('newPart', '.ga-new-part')"), html.Text("Add")),
)
