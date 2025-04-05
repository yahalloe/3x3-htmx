package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	port := ":3000"
	staticDir := http.Dir("./src")

	// Serve static files under /src
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public")))) // Serves images
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(staticDir)))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Custom 404 for unknown paths
			http.NotFound(w, r)
			return
		}

		// Serve the index.html file when the page loads or when hx-get is triggered
		if r.Header.Get("HX-Request") == "true" {
			// If it's an HTMX request, you might want to serve only a specific part of the page
			// For example, just the content inside the main div:
			http.ServeFile(w, r, filepath.Join("src", "index.html"))
		} else {
			// Regular request
			http.ServeFile(w, r, filepath.Join("src", "index.html"))
		}
	})

	http.HandleFunc("/edgerunners", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, filepath.Join("src", "edgerunners.html"))
	})

	http.HandleFunc("/romance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, filepath.Join("src", "romance.html"))
	})

	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		// Create a new request to the target URL
		targetURL := "https://myanimelist.net/anime/12189/Hyouka"
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Set headers to mimic a browser request
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

		// Copy the status code and body
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
