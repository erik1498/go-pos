package pos

import (
	"go-pos/internal/model"
	"go-pos/internal/repository/category"

	"github.com/google/uuid"
)

type posUsecase struct {
	cRepo category.Repository
}

func GetUsecase(cRepo category.Repository) Usecase {
	return &posUsecase{
		cRepo: cRepo,
	}
}

func (pUsecase *posUsecase) GetCategoryList() ([]model.Category, error) {
	return pUsecase.cRepo.GetCategoryList()
}

func (pUsecase *posUsecase) CreateCategory(req model.Category) (model.Category, error) {
	category := model.Category{
		PublicID: uuid.Must(uuid.NewV7()),
		Name:     req.Name,
	}

	return pUsecase.cRepo.CreateCategory(category)
}

func (pUsecase *posUsecase) GetCategoryById(id uuid.UUID) (model.Category, error) {
	return pUsecase.cRepo.GetCategoryById(id)
}

func (pUsecase *posUsecase) UpdateCategoryById(req model.Category, id uuid.UUID) (model.Category, error) {
	return pUsecase.cRepo.UpdateCategoryById(req, id)
}

func (pUsecase *posUsecase) DeleteCategoryByID(id uuid.UUID) error {
	return pUsecase.cRepo.DeleteCategoryByID(id)
}
