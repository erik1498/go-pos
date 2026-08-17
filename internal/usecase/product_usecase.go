package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type productUsecase struct {
	pRepo domain.ProductRepository
	cRepo domain.CategoryRepository
}

func GetProductUsecase(pRepo domain.ProductRepository, cRepo domain.CategoryRepository) domain.ProductUsecase {
	return &productUsecase{
		pRepo: pRepo,
		cRepo: cRepo,
	}
}

func (pUsecase *productUsecase) GetAll(opts domain.QueryOptions) ([]model.Product, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
		"sku":  true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return pUsecase.pRepo.GetAll(cleanOpts)
}

func (pUsecase *productUsecase) Create(req model.ProductRequest) (model.Product, error) {
	category, err := pUsecase.cRepo.GetByID(uuid.MustParse(req.CategoryID))

	if err != nil {
		return model.Product{}, err
	}

	product := model.Product{
		ID:         uuid.Must(uuid.NewV7()),
		CategoryID: category.ID,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
	}

	return pUsecase.pRepo.Create(product)
}

func (pUsecase *productUsecase) GetByID(id uuid.UUID) (model.Product, error) {
	product, err := pUsecase.pRepo.GetByID(id)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (pUsecase *productUsecase) UpdateByID(id uuid.UUID, req model.ProductRequest) (model.Product, error) {
	category, err := pUsecase.cRepo.GetByID(uuid.MustParse(req.CategoryID))
	if err != nil {
		return model.Product{}, err
	}

	product := model.Product{
		CategoryID: category.ID,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
	}

	product, err = pUsecase.pRepo.UpdateByID(id, product)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (pUsecase *productUsecase) DeleteByID(id uuid.UUID) error {
	return pUsecase.pRepo.DeleteByID(id)
}
