package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Change port if needed (use :8443 for local testing)
	port := ":443" // change back to ":443" in production if running with proper privileges
	certFile := "/etc/letsencrypt/live/yahallo.tech/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/yahallo.tech/privkey.pem"
	staticDir := http.Dir("./src")

	// Log each request for debugging purposes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("src", "index.html"))
		} else {
			http.NotFound(w, r)
		}
	})

	// Serve static files under /public and /src with correct MIME types
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(staticDir)))

	//  Serve additional HTML files
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

	// Reverse Proxy Example with improved logging
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		targetURL := "https://myanimelist.net/anime/12189/Hyouka"
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Mimic a browser request
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to fetch data: %v", err)
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy response headers and body
		w.Header().Set("Access-Control-Allow-Origin", "*")
		contentType := resp.Header.Get("Content-Type")
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error copying response: %v", err)
		}
	})

	log.Printf("Server running at https://localhost%s\n", port)
	// For local testing, you might not have valid TLS certificates.
	// http.ListenAndServe(port, nil)
	err := http.ListenAndServeTLS(port, certFile, keyFile, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
