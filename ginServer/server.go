// production server

package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start the HTTP redirect server on port 80
	go func() {
		httpRouter := gin.Default()
		httpRouter.GET("/*path", func(c *gin.Context) {
			target := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, target)
		})
		if err := httpRouter.Run(":80"); err != nil {
			log.Fatalf("HTTP redirect server failed: %v", err)
		}
	}()

	// HTTPS server setup
	e := gin.Default()

	// Load templates
	e.LoadHTMLGlob("../src/*.html")

	// Serve static files
	e.Static("/src", "../src")
	e.Static("/public", "../public")

	// Home route
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Define valid HTMX-style routes
	htmxRoutes := map[string]bool{
		"/edgerunners": true,
		"/romance":     true,
		"/madeinabyss": true,
		"/op":          true,
	}

	// Dynamic route handling
	e.GET("/:page", func(c *gin.Context) {
		page := "/" + c.Param("page")
		if htmxRoutes[page] {
			path := filepath.Join("../src", c.Param("page")+".html")
			c.File(path)
			return
		}
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// TLS certs from Let's Encrypt
	certFile := "/etc/letsencrypt/live/yahallo.tech/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/yahallo.tech/privkey.pem"

	log.Println("HTTPS server running on :443")
	if err := e.RunTLS(":443", certFile, keyFile); err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
