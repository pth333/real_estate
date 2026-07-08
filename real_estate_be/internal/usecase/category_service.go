package usecase

import (
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type CategoryService struct {
	repo repo.ICategoryRepository
}

type ICategoryService interface {
	GetAll() ([]model.Category, error)
	BuildCategoriesResponse(categories []model.Category) []model.Category
}

func NewCategoryService(repo repo.ICategoryRepository) ICategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) BuildCategoriesResponse(categories []model.Category) []model.Category {
	// Build map ID → pointer (để sửa trực tiếp trên slice gốc)
	mapCategories := make(map[int64]*model.Category, len(categories))
	for i := range categories {
		mapCategories[categories[i].ID] = &categories[i]
	}

	var roots []model.Category
	for i := range categories {
		if categories[i].ParentID == nil {
			roots = append(roots, categories[i])
		} else {
			if parent, ok := mapCategories[*categories[i].ParentID]; ok {
				parent.Children = append(parent.Children, &categories[i])
			}
		}
	}
	return roots
}
