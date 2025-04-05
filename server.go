package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define port and directories
	port := ":3000"
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

	// Load certificates for HTTPS
	cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/yahallo.tech/fullchain.pem", "/etc/letsencrypt/live/yahallo.tech/privkey.pem")
	if err != nil {
		log.Fatalf("Error loading certificates: %v", err)
	}

	// Create a custom HTTP server with TLS (without HTTP/2)
	server := &http.Server{
		Addr:    port,
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	// Start HTTPS server (without HTTP/2)
	log.Printf("Server running at https://localhost%s\n", port)
	log.Fatal(server.ListenAndServeTLS(":80", ""))
}
