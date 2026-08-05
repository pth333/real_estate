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
	// Build map ID → con trỏ tới node riêng trên heap.
	// Mỗi node copy ra từ categories[i] bằng biến mới (c) để con trỏ ổn định,
	// tránh lỗi copy value làm mất children (fix: con bị append sau khi root đã copy vào roots).
	nodes := make(map[int64]*model.Category, len(categories))
	var roots []*model.Category
	for i := range categories {
		c := categories[i]
		nodes[c.ID] = &c
		if c.ParentID == nil {
			roots = append(roots, &c)
		}
	}

	// Gắn children vào parent qua nodes
	for _, c := range categories {
		if c.ParentID == nil {
			continue
		}
		if parent, ok := nodes[*c.ParentID]; ok {
			parent.Children = append(parent.Children, nodes[c.ID])
		}
	}

	// Deref các root để trả về []model.Category
	result := make([]model.Category, 0, len(roots))
	for _, r := range roots {
		result = append(result, *r)
	}
	return result
}
