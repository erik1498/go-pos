package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type taxUsecase struct {
	tRepo domain.TaxRepository
}

func GetTaxUsecase(tRepo domain.TaxRepository) domain.TaxUsecase {
	return &taxUsecase{
		tRepo: tRepo,
	}
}

func (tUsecase *taxUsecase) GetAll(opts domain.QueryOptions) ([]model.Tax, int64, error) {
	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	allowedFields := map[string]bool{
		"name": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return tUsecase.tRepo.GetAll(cleanOpts)
}

func (tUsecase *taxUsecase) Create(req model.TaxRequest) (model.Tax, error) {
	ID := uuid.Must(uuid.NewV7())

	tax, err := tUsecase.tRepo.Create(model.Tax{
		ID:   ID,
		Name: req.Name,
		Rate: req.Rate,
	})

	if err != nil {
		return model.Tax{}, err
	}

	return tax, nil
}

func (tUsecase *taxUsecase) GetByID(id uuid.UUID) (model.Tax, error) {
	tax, err := tUsecase.tRepo.GetByID(id)
	if err != nil {
		return model.Tax{}, err
	}

	return tax, nil
}

func (tUsecase *taxUsecase) UpdateByID(id uuid.UUID, req model.TaxRequest) (model.Tax, error) {
	tax, err := tUsecase.tRepo.UpdateByID(id, model.Tax{
		Name: req.Name,
		Rate: req.Rate,
	})

	return tax, err
}

func (tUsecase *taxUsecase) DeleteByID(id uuid.UUID) error {
	return tUsecase.tRepo.DeleteByID(id)
}
