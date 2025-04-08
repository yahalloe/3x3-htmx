// local server

package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()

	// HTML templates in ../src
	e.LoadHTMLGlob("../src/*.html")

	// Serve static files
	e.Static("/src", "../src")       // for style.css, script.js, etc.
	e.Static("/public", "../public") // for images

	// Serve the index.html on root route
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	htmxRoutes := map[string]bool{
		"/edgerunners": true,
		"/romance":     true,
		"/madeinabyss": true,
		"/op":          true,
		"/aboutme":     true,
	}

	// Catch routes and serve files if valid
	e.GET("/:page", func(c *gin.Context) {
		page := "/" + c.Param("page")
		if htmxRoutes[page] {
			path := filepath.Join("../src", c.Param("page")+".html")
			c.File(path)
			return
		}
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	e.Run(":8080") // or whatever port you want
}
