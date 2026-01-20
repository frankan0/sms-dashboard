package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"sms-dashboard/internal/config"
	"sms-dashboard/internal/database"
	"sms-dashboard/internal/router"
	"sms-dashboard/web"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg.DatabasePath)

	r := router.SetupRouter(cfg)

	// Serve Frontend
	distFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	r.NoRoute(func(c *gin.Context) {
		// 如果是 API 请求，返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")

		// 处理根路径或前端路由
		if path == "" || !strings.Contains(path, ".") {
			file, err := distFS.Open("index.html")
			if err != nil {
				c.String(http.StatusNotFound, "index.html not found")
				return
			}
			defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				c.String(http.StatusInternalServerError, "error reading index.html")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", content)
			return
		}

		// 正常提供静态文件 (assets 等)
		f, err := distFS.Open(path)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		f.Close()

		c.FileFromFS(path, http.FS(distFS))
	})

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
