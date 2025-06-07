package handlers

import "net/http"

func InitRouter() *http.ServeMux {
	router := http.NewServeMux()

	// NOTE: init routes here

	return router
}
