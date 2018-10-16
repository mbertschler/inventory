package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["searchList"] = searchListAction
}

func listPage(w http.ResponseWriter, r *http.Request) {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(listBlock(parts))
	err = html.Render(page, w)
	if err != nil {
		log.Println(err)
	}
}

func searchListAction(args json.RawMessage) (*guiapi.Result, error) {
	var in string
	err := json.Unmarshal(args, &in)
	if err != nil {
		return nil, err
	}
	parts, err := parts.Search(in)
	if err != nil {
		log.Println(err)
	}
	return guiapi.Replace("#partsList", listOnlyBlock(parts))
}

func listBlock(parts []*parts.Part) html.Block {
	return html.Blocks{
		listControls,
		listOnlyBlock(parts),
	}
}

func listOnlyBlock(parts []*parts.Part) html.Block {
	list := html.Blocks{}
	for _, p := range parts {
		link := fmt.Sprintf("/part/%s", p.ID())
		list.Add(html.Elem("tr", html.Attr("onclick", "redirect('"+link+"')"),
			html.Elem("td", nil, html.Text(p.Code)),
			html.Elem("td", nil, html.Text(p.Name)),
			html.Elem("td", nil, html.Text(p.Type)),
			html.Elem("td", nil, html.Text(p.Value)),
			html.Elem("td", nil, html.Text(p.Size)),
			html.Elem("td", nil, html.Text(fmt.Sprint(p.Quantity))),
			html.Elem("td", nil, html.Text(p.Location)),
		))
	}
	return html.Div(html.Id("partsList").Class("ui list"),
		html.Elem("table", html.Class("ui celled compact table"),
			html.Elem("thead", nil,
				html.Elem("tr", nil,
					html.Elem("th", nil, html.Text("Code")),
					html.Elem("th", nil, html.Text("Name")),
					html.Elem("th", nil, html.Text("Type")),
					html.Elem("th", nil, html.Text("Value")),
					html.Elem("th", nil, html.Text("Size")),
					html.Elem("th", nil, html.Text("Quantity")),
					html.Elem("th", nil, html.Text("Location")),
				),
			), html.Elem("tbody", nil,
				list,
			),
		),
	)
}

var listControls = html.Div(html.Class("ui grid"),
	html.Div(html.Class("thirteen wide column"),
		html.Div(html.Class("ui fluid icon input").Attr("oninput", "sendInput('searchList', event)"),
			html.Input(html.Type("text").Attr("placeholder", "Search parts...")),
			html.I(html.Class("search icon")),
		),
	),
	html.Div(html.Class("three wide column right floated"),
		html.Div(html.Class("ui button primary").Attr("onclick", "guiapi('newPart', null)"), html.Text("Add Part")),
	),
)
