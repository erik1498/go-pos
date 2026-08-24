package usecase

import (
	"context"
	"encoding/json"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type categoryUsecase struct {
	cRepo domain.CategoryRepository
	aRepo domain.AuditLogRepository
}

func GetCategoryUsecase(cRepo domain.CategoryRepository, aRepo domain.AuditLogRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		cRepo: cRepo,
		aRepo: aRepo,
	}
}

func (cUsecase *categoryUsecase) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return cUsecase.cRepo.GetAll(cleanOpts)
}

func (cUsecase *categoryUsecase) Create(ctx context.Context, req model.Category) (model.Category, error) {
	category := model.Category{
		ID:   uuid.Must(uuid.NewV7()),
		Name: req.Name,
	}

	createdCategory, err := cUsecase.cRepo.Create(category)
	if err != nil {
		return model.Category{}, err
	}

	meta, ok := ctx.Value(rest.AuditMetaKey).(rest.AuditMeta)
	if ok {
		newValueJSON, _ := json.Marshal(createdCategory)

		auditLog := model.AuditLog{
			ActorID:   uuid.MustParse(meta.UserID),
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "categories",
			EntityID:  createdCategory.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			cUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}
	return createdCategory, nil
}

func (cUsecase *categoryUsecase) GetByID(id uuid.UUID) (model.Category, error) {
	return cUsecase.cRepo.GetByID(id)
}

func (cUsecase *categoryUsecase) UpdateCategoryByID(id uuid.UUID, req model.Category) (model.Category, error) {
	return cUsecase.cRepo.UpdateCategoryByID(id, req)
}

func (cUsecase *categoryUsecase) DeleteCategoryByID(id uuid.UUID) error {
	return cUsecase.cRepo.DeleteCategoryByID(id)
}
