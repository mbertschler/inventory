package gui

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mbertschler/blocks/html"
)

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	page := mainLayout(html.Text("hihi"))
	err := html.Render(page, w)
	if err != nil {
		log.Println(err)
	}
}

func mainLayout(content html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		layoutHead(),
		layoutBody(content),
	}
}

func layoutHead() html.Block {
	return html.Head(nil,
		html.Meta(html.Charset("utf-8")),
		html.Meta(html.Attr{{Key: "http-equiv", Value: "X-UA-Compatible"}}.Content("IE=edge,chome=1")),
		html.Meta(html.Name("viewport").Content("width=device-width, initial-scale=1.0, maximum-scale=1.0")),
		html.Meta(html.Name("apple-mobile-web-app-capable").Content("yes")),
		html.Title(nil,
			html.Text("Inventory"),
		),
		html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/semantic-ui@2.4.0/dist/semantic.min.css")),
	)
}

func addRefreshQuery(in string) string {
	return fmt.Sprint(in, "?q=", time.Now().Unix())
}

func layoutBody(content html.Block) html.Block {
	return html.Body(nil,
		html.H1(html.Class("ui center aligned header").Styles("padding:32px 0 16px"),
			html.Text("Inventory Manager")),
		html.Div(html.Class("ui container").Id("container"),
			content,
		),
		html.Script(html.Src("https://cdn.jsdelivr.net/npm/jquery@3.3.1/dist/jquery.min.js")),
		html.Script(html.Src("https://cdn.jsdelivr.net/npm/semantic-ui@2.4.0/dist/semantic.min.js")),
		html.Script(html.Src(addRefreshQuery("/js/app.js"))),
	)
}
