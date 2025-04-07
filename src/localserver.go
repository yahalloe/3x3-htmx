package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := ":3000"
	staticDir := http.Dir("./src")

	// Serve static files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(staticDir)))

	// Handle all other routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Special case: root path
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("src", "index.html"))
			return
		}

		// List of valid HTML routes (without .html extension)
		validHTMLRoutes := map[string]bool{
			"/edgerunners": true,
			"/romance":     true,
			// Add other valid paths here
		}

		// Handle valid HTML routes
		if validHTMLRoutes[r.URL.Path] {
			http.ServeFile(w, r, filepath.Join("src", r.URL.Path+".html"))
			return
		}

		// Check if file exists in /public or /src
		if fileExists(filepath.Join("public", r.URL.Path)) || fileExists(filepath.Join("src", r.URL.Path)) {
			http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
			return
		}

		// Handle HTMX requests differently
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Retarget", "#htmx-target")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, filepath.Join("src", "404.html"))
			return
		}

		// Regular 404 for full page loads
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, filepath.Join("src", "404.html"))
	})

	log.Printf("Server running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
