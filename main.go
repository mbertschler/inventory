package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/mbertschler/inventory/gui"
	"github.com/mbertschler/inventory/lib/guiapi"
	"github.com/mbertschler/inventory/parts"
)

//go:embed gui/js
var guiFS embed.FS

func fileServer(url string, path ...string) http.Handler {
	dir := http.Dir(filepath.Join(path...))
	return http.StripPrefix(url,
		http.FileServer(dir))
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("first argument needs to be a database file path")
	}
	err := parts.SetupDB(args[0])
	if err != nil {
		log.Fatal(err)
	}

	static, _ := fs.Sub(guiFS, "gui/js")

	log.Println("inventory server at :5080")
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.FS(static))))
	http.Handle("/guiapi/", guiapi.DefaultHandler)
	http.Handle("/", gui.Router())
	log.Println(http.ListenAndServe(":5080", nil))
}
