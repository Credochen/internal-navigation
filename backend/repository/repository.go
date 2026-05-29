package repository

import (
	"internal-navigation/backend/model"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) ListCategories() ([]model.Category, error) {
	var cats []model.Category
	err := r.db.Order("sort_order asc, id asc").Find(&cats).Error
	return cats, err
}

func (r *Repo) CreateCategory(c *model.Category) error {
	return r.db.Create(c).Error
}

func (r *Repo) UpdateCategory(id uint, c *model.Category) error {
	return r.db.Model(&model.Category{}).Where("id = ?", id).Updates(map[string]any{
		"name":       c.Name,
		"sort_order": c.SortOrder,
	}).Error
}

func (r *Repo) DeleteCategory(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *Repo) ListNavigations() ([]model.Navigation, error) {
	var navs []model.Navigation
	err := r.db.Order("sort_order asc, id asc").Find(&navs).Error
	return navs, err
}

func (r *Repo) CreateNavigation(n *model.Navigation) error {
	return r.db.Create(n).Error
}

func (r *Repo) UpdateNavigation(id uint, n *model.Navigation) error {
	return r.db.Model(&model.Navigation{}).Where("id = ?", id).Updates(map[string]any{
		"category_id": n.CategoryID,
		"title":       n.Title,
		"url":         n.URL,
		"icon":        n.Icon,
		"description": n.Description,
		"sort_order":  n.SortOrder,
	}).Error
}

func (r *Repo) DeleteNavigation(id uint) error {
	return r.db.Delete(&model.Navigation{}, id).Error
}

func (r *Repo) GetNavTree() ([]model.NavTree, error) {
	cats, err := r.ListCategories()
	if err != nil {
		return nil, err
	}
	var all []model.Navigation
	if err := r.db.Order("sort_order asc, id asc").Find(&all).Error; err != nil {
		return nil, err
	}
	m := make(map[uint][]model.NavItem)
	for _, nav := range all {
		m[nav.CategoryID] = append(m[nav.CategoryID], model.NavItem{
			ID:          nav.ID,
			Title:       nav.Title,
			URL:         nav.URL,
			Icon:        nav.Icon,
			Description: nav.Description,
		})
	}
	tree := make([]model.NavTree, 0)
	for _, cat := range cats {
		items := m[cat.ID]
		if items == nil {
			items = []model.NavItem{}
		}
		tree = append(tree, model.NavTree{
			CategoryID:   cat.ID,
			CategoryName: cat.Name,
			Items:        items,
		})
	}
	return tree, nil
}
