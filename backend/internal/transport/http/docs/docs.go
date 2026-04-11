package docs

import (
	"embed"
	"net/http"
)

//go:embed openapi.yaml swagger.html
var files embed.FS

func OpenAPIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		http.ServeFileFS(w, r, files, "openapi.yaml")
	})
}

func SwaggerUIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFileFS(w, r, files, "swagger.html")
	})
}
