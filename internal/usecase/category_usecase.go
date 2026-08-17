package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type categoryUsecase struct {
	cRepo domain.CategoryRepository
}

func GetCategoryUsecase(cRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		cRepo: cRepo,
	}
}

func (pUsecase *categoryUsecase) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return pUsecase.cRepo.GetAll(cleanOpts)
}

func (pUsecase *categoryUsecase) Create(req model.Category) (model.Category, error) {
	category := model.Category{
		ID:   uuid.Must(uuid.NewV7()),
		Name: req.Name,
	}

	return pUsecase.cRepo.Create(category)
}

func (pUsecase *categoryUsecase) GetByID(id uuid.UUID) (model.Category, error) {
	return pUsecase.cRepo.GetByID(id)
}

func (pUsecase *categoryUsecase) UpdateCategoryByID(id uuid.UUID, req model.Category) (model.Category, error) {
	return pUsecase.cRepo.UpdateCategoryByID(id, req)
}

func (pUsecase *categoryUsecase) DeleteCategoryByID(id uuid.UUID) error {
	return pUsecase.cRepo.DeleteCategoryByID(id)
}
