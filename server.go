package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := ":443"
	certFile := "/etc/letsencrypt/live/yahallo.tech/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/yahallo.tech/privkey.pem"

	// Define all valid routes
	validRoutes := map[string]string{
		"/":            "index.html",
		"/edgerunners": "edgerunners.html",
		"/romance":     "romance.html",
	}

	// Serve static files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("src"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Check valid routes
		if htmlFile, ok := validRoutes[r.URL.Path]; ok {
			serveHTML(w, r, filepath.Join("src", htmlFile))
			return
		}

		// Check static files
		if tryServeFile(w, r, "public", r.URL.Path) || tryServeFile(w, r, "src", r.URL.Path) {
			return
		}

		// HTMX-aware 404 handling
		handle404(w, r)
	})

	log.Printf("Starting production server on port %s", port)
	if err := http.ListenAndServeTLS(port, certFile, keyFile, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request, path string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, path)
}

func tryServeFile(w http.ResponseWriter, r *http.Request, dir, path string) bool {
	fullPath := filepath.Join(dir, path)
	if _, err := os.Stat(fullPath); err == nil {
		http.ServeFile(w, r, fullPath)
		return true
	}
	return false
}

func handle404(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Retarget", "#htmx-target")
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, filepath.Join("src", "404.html"))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, filepath.Join("src", "404.html"))
}
