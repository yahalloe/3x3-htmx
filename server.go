package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define port and directories
	port := ":443"
	certFile := "/etc/letsencrypt/live/yahallo.tech/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/yahallo.tech/privkey.pem"
	staticDir := http.Dir("./src")

	// Serve static files under /public and /src
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(staticDir)))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("src", "index.html"))
	})

	// Serve other HTML files
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

	// Reverse Proxy Example
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		targetURL := "https://myanimelist.net/anime/12189/Hyouka"
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Mimic a browser request
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Set CORS headers and copy the response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Start HTTPS server
	log.Printf("Server running at https://localhost%s\n", port)

	// Using http.DefaultServeMux for the handler
	log.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
}
