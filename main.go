package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"git.exahome.net/tools/inventory/gui"
	"git.exahome.net/tools/inventory/lib/guiapi"
)

var root string

func init() {
	gopath := os.Getenv("GOPATH")
	paths := filepath.SplitList(gopath)
	for _, p := range paths {
		project := filepath.Join(p, "src", "git.exahome.net", "tools", "inventory")
		info, err := os.Stat(project)
		if err == nil && info.IsDir() {
			root = project
			return
		}
	}
	log.Fatal("couldn't find the project in GOPATH")
}

func fileServer(url string, path ...string) http.Handler {
	dir := http.Dir(filepath.Join(path...))
	return http.StripPrefix(url,
		http.FileServer(dir))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("inventory server at :5080")
	http.Handle("/js/", fileServer("/js/", root, "gui", "js"))
	http.Handle("/guiapi/", guiapi.DefaultHandler)
	http.HandleFunc("/", gui.HandleFunc)
	log.Println(http.ListenAndServe(":5080", nil))
}
