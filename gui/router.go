package gui

import (
	"net/http"
)

// Router handles all the GUI routes
func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/part/", partPage)
	return router
}
