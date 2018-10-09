package gui

import (
	"encoding/json"
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
}

func partPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/part/")
	part, err := parts.ByID(id)
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(partBlock(part))
	err = html.Render(page, w)
	if err != nil {
		log.Println(err)
	}
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
	return guiapi.Replace("#container", partBlock(part))
}

func partBlock(p *parts.Part) html.Block {
	return html.Div(nil,
		html.H1(nil, html.Text(p.Name)),
	)
}
