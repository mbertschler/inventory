package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"git.exahome.net/tools/inventory/gui"
	"git.exahome.net/tools/inventory/lib/guiapi"
	"git.exahome.net/tools/inventory/parts"
)

var root string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gopath := os.Getenv("GOPATH")
	paths := filepath.SplitList(gopath)
	for _, p := range paths {
		project := filepath.Join(p, "src", "git.exahome.net", "tools", "inventory")
		info, err := os.Stat(project)
		if err == nil && info.IsDir() {
			root = project
			break
		}
	}
	if root == "" {
		log.Fatal("couldn't find the project in GOPATH")
	}
	err := parts.Add("Capacitor")
	if err != nil {
		log.Fatal(err)
	}
	err = parts.Add("Resistor")
	if err != nil {
		log.Fatal(err)
	}
}

func fileServer(url string, path ...string) http.Handler {
	dir := http.Dir(filepath.Join(path...))
	return http.StripPrefix(url,
		http.FileServer(dir))
}

func main() {
	log.Println("inventory server at :5080")
	http.Handle("/js/", fileServer("/js/", root, "gui", "js"))
	http.Handle("/guiapi/", guiapi.DefaultHandler)
	http.HandleFunc("/", gui.HandleFunc)
	log.Println(http.ListenAndServe(":5080", nil))
}
