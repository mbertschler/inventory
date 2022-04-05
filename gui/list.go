package gui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"

	"github.com/mbertschler/inventory/lib/guiapi"
	"github.com/mbertschler/inventory/parts"
)

func init() {
	// setup guiapi action
	guiapi.DefaultHandler.Functions["searchList"] = searchListAction
	guiapi.DefaultHandler.Functions["startScan"] = startScanAction
	guiapi.DefaultHandler.Functions["stopScan"] = stopScanAction
	guiapi.DefaultHandler.Functions["scanCode"] = scanCodeAction
	guiapi.DefaultHandler.Functions["clearScan"] = clearScanAction
}

func listPage(w http.ResponseWriter, r *http.Request) {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	page := mainLayout(listBlock(parts, false))
	err = html.Render(w, page)
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

func startScanAction(args json.RawMessage) (*guiapi.Result, error) {
	res, err := guiapi.Replace("#listControls", listControls(true))
	if err != nil {
		return res, err
	}
	res.JS = append(res.JS, guiapi.JSCall{Name: "startScan"})
	return res, nil
}

func stopScanAction(args json.RawMessage) (*guiapi.Result, error) {
	res, err := guiapi.Replace("#listControls", listControls(false))
	if err != nil {
		return res, err
	}
	res.JS = append(res.JS, guiapi.JSCall{Name: "stopScan"})
	return res, nil
}

func scanCodeAction(args json.RawMessage) (*guiapi.Result, error) {
	var code string
	err := json.Unmarshal(args, &code)
	if err != nil {
		return nil, err
	}
	parts, err := parts.Search(code)
	if err != nil {
		log.Println(err)
	}
	if len(parts) == 0 {
		return guiapi.Replace("#container", editPartBlock(nil, code))
	}
	if len(parts) == 1 {
		return guiapi.Replace("#container", viewPartBlock(parts[0]))
	}
	blocks := html.Blocks{
		html.Div(html.Class("ui teal message grid").Attr("onclick", "guiapi('clearScan', null)"),
			html.I(html.Class("close icon")),
			html.Strong(nil, html.Text("Scanned Code:")),
			html.Text(code),
			html.Br(),
			html.Text("click to clear"),
		),
		listOnlyBlock(parts),
	}
	return guiapi.Replace("#partsList", blocks)
}

func clearScanAction(args json.RawMessage) (*guiapi.Result, error) {
	parts, err := parts.All()
	if err != nil {
		log.Println(err)
	}
	return guiapi.Replace("#partsList", listOnlyBlock(parts))
}

func listBlock(parts []*parts.Part, scanning bool) html.Block {
	return html.Blocks{
		html.Div(html.Class("ui grid").Id("listControls"),
			listControls(scanning),
		),
		html.Div(html.Id("partsList").Class("ui list"),
			listOnlyBlock(parts),
		),
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
	return html.Elem("table", html.Class("ui celled compact table"),
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
	)
}

func listControls(scanning bool) html.Block {
	var scanButton, video html.Block
	if scanning {
		scanButton = html.Div(html.Class("ui button orange").Attr("onclick", "guiapi('stopScan', null)"),
			html.Text("Stop Scanner"))
		video = html.Elem("video", html.Id("scanVideo").Styles("height:240px;margin-bottom:14px"))
	} else {
		scanButton = html.Div(html.Class("ui button yellow").Attr("onclick", "guiapi('startScan', null)"),
			html.Text("Start Scanner"))
	}
	return html.Blocks{
		html.Div(html.Class("ten wide column"),
			html.Div(html.Class("ui fluid icon input").Attr("oninput", "sendInput('searchList', event)"),
				html.Input(html.Type("text").Id("autofocus").Attr("placeholder", "Search parts...")),
				html.I(html.Class("search icon")),
			),
		),
		html.Div(html.Class("six wide column right floated"),
			html.Div(html.Class("ui button primary").Attr("onclick", "guiapi('newPart', null)"), html.Text("Add Part")),
			scanButton,
		),
		video,
	}
}
