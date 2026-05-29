package main

import (
	"embed"
	"io/fs"
	"net/http"

	"internal-navigation/backend/config"
	"internal-navigation/backend/handler"
	"internal-navigation/backend/model"
	"internal-navigation/backend/repository"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var distFS embed.FS

func main() {
	cfg := config.Load()

	if err := model.InitDB(cfg.DBPath); err != nil {
		panic(err)
	}

	repo := repository.New(model.DB)
	h := handler.New(repo, cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	h.Register(r)

	staticSub, _ := fs.Sub(distFS, "dist")
	handler.ServeStaticOrIndex(r, http.FS(staticSub))

	r.Run(":" + cfg.Port)
}
