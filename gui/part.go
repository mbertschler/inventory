package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbertschler/blocks/html"

	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["viewPart"] = viewPartAction
	// guiapi.DefaultHandler.Functions["editPart"] = editPartAction
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

// func editPartAction(args json.RawMessage) (*guiapi.Result, error) {
// 	var id string
// 	err := json.Unmarshal(args, &id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	part, err := parts.ByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return guiapi.Replace("#container", partBlock(part))
// }

func deletePartAction(args json.RawMessage) (*guiapi.Result, error) {
	var id string
	err := json.Unmarshal(args, &id)
	if err != nil {
		return nil, err
	}
	log.Println("deleting part", id)
	err = parts.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return guiapi.Redirect("/")
}

func deletePartBlock(p *parts.Part) html.Block {
	editAction := fmt.Sprintf("guiapi('editPart', '%s')", p.ID())
	deleteAction := fmt.Sprintf("guiapi('deletePart', '%s')", p.ID())
	return html.Div(nil,
		html.Div(nil,
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
	)
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
	return html.Div(nil,
		html.Div(nil,
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
	)
}
