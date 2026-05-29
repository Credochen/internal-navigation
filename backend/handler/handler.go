package handler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"internal-navigation/backend/config"
	"internal-navigation/backend/model"
	"internal-navigation/backend/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo   *repository.Repo
	cfg    *config.Config
	online sync.Map
}

func New(repo *repository.Repo, cfg *config.Config) *Handler {
	h := &Handler{repo: repo, cfg: cfg}
	go h.healthCheckLoop()
	return h
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/navs", h.GetNavs)
		api.GET("/categories", h.ListCategories)
	}

	admin := api.Group("/admin")
	admin.Use(h.adminAuth())
	{
		admin.POST("/categories", h.CreateCategory)
		admin.PUT("/categories/:id", h.UpdateCategory)
		admin.DELETE("/categories/:id", h.DeleteCategory)

		admin.POST("/navs", h.CreateNav)
		admin.PUT("/navs/:id", h.UpdateNav)
		admin.DELETE("/navs/:id", h.DeleteNav)

		admin.POST("/import", h.Import)
		admin.GET("/export", h.Export)
	}
}

func (h *Handler) adminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		if token != h.cfg.AdminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func (h *Handler) GetNavs(c *gin.Context) {
	tree, err := h.repo.GetNavTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range tree {
		for j := range tree[i].Items {
			if v, ok := h.online.Load(tree[i].Items[j].ID); ok {
				tree[i].Items[j].Online = v.(bool)
			}
		}
	}
	c.JSON(http.StatusOK, tree)
}

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.repo.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var cat model.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.CreateCategory(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var cat model.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := parseUint(c.Param("id"))
	if err := h.repo.UpdateCategory(id, &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := h.repo.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CreateNav(c *gin.Context) {
	var nav model.Navigation
	if err := c.ShouldBindJSON(&nav); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.CreateNavigation(&nav); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nav)
}

func (h *Handler) UpdateNav(c *gin.Context) {
	var nav model.Navigation
	if err := c.ShouldBindJSON(&nav); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := parseUint(c.Param("id"))
	if err := h.repo.UpdateNavigation(id, &nav); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteNav(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := h.repo.DeleteNavigation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type importExportData struct {
	Categories   []model.Category   `json:"categories"`
	Navigations []model.Navigation `json:"navigations"`
}

func (h *Handler) Import(c *gin.Context) {
	var data importExportData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range data.Categories {
		data.Categories[i].ID = 0
		if err := h.repo.CreateCategory(&data.Categories[i]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	for i := range data.Navigations {
		data.Navigations[i].ID = 0
		if err := h.repo.CreateNavigation(&data.Navigations[i]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Export(c *gin.Context) {
	cats, err := h.repo.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	navs, err := h.repo.ListNavigations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, importExportData{Categories: cats, Navigations: navs})
}

func (h *Handler) healthCheckLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	h.checkAll()
	for range ticker.C {
		h.checkAll()
	}
}

func (h *Handler) checkAll() {
	navs, err := h.repo.ListNavigations()
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	for _, nav := range navs {
		wg.Add(1)
		go func(id uint, url string) {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}
			req, _ := http.NewRequest("HEAD", url, nil)
			resp, err := client.Do(req)
			online := err == nil && resp.StatusCode < 500
			if resp != nil {
				resp.Body.Close()
			}
			h.online.Store(id, online)
		}(nav.ID, nav.URL)
	}
	wg.Wait()
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func ServeStaticOrIndex(r *gin.Engine, staticFS http.FileSystem) {
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		f, err := staticFS.Open(path)
		if err == nil {
			stat, _ := f.Stat()
			if stat != nil && !stat.IsDir() {
				c.FileFromFS(path, staticFS)
				return
			}
			f.Close()
		}
		index, err := staticFS.Open("/index.html")
		if err != nil {
			c.String(http.StatusNotFound, "not found")
			return
		}
		defer index.Close()
		stat, _ := index.Stat()
		c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", index, nil)
	})
}
