package gui

import (
	"fmt"
	"log"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/parts"
)

func listBlock() html.Block {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	list := html.Blocks{}
	for _, p := range parts {
		link := fmt.Sprintf("/part/%s", p.ID())
		list.Add(html.Elem("tr", nil,
			html.Elem("td", nil, html.Text(p.Code)),
			html.Elem("td", nil, html.Text(p.Name)),
			html.Elem("td", nil, html.Text(p.Type)),
			html.Elem("td", nil, html.Text(p.Value)),
			html.Elem("td", nil, html.Text(p.Size)),
			html.Elem("td", nil, html.Text(fmt.Sprint(p.Quantity))),
			html.Elem("td", nil, html.Text(p.Location)),
			html.Elem("td", nil, html.A(html.Href(link), html.Text("View Part"))),
		))
	}
	return html.Div(html.Class("ui list"),
		newPartForm,
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
					html.Elem("th", nil, html.Text("View Part")),
				),
			), html.Elem("tbody", nil,
				list,
			),
		),
	)
}

var newPartForm = html.Div(nil,
	html.Div(html.Class("ui button primary").Attr("onclick", "guiapi('newPart', null)"), html.Text("Add Part")),
)
