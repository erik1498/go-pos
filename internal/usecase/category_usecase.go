package usecase

import (
	"context"
	"encoding/json"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/middleware"
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

func (cUsecase *categoryUsecase) Create(ctx context.Context, req model.CategoryRequest) (model.Category, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return model.Category{}, domain.ErrIdempotencyKeyDuplicate
	}

	category := model.Category{
		ID:             uuid.Must(uuid.NewV7()),
		Name:           req.Name,
		CreatedBy:      actorID,
		UpdatedBy:      actorID,
		IdempotencyKey: idemKey,
	}

	category, err := cUsecase.cRepo.Create(category)
	if err != nil {
		return model.Category{}, err
	}

	if metaValid {
		newValueJSON, _ := json.Marshal(category)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "categories",
			EntityID:  category.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			cUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}

	return category, nil
}

func (cUsecase *categoryUsecase) GetByID(id uuid.UUID) (model.Category, error) {
	return cUsecase.cRepo.GetByID(id)
}

func (cUsecase *categoryUsecase) UpdateCategoryByID(ctx context.Context, id uuid.UUID, req model.CategoryRequest) (model.Category, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	category := model.Category{
		ID:        id,
		Name:      req.Name,
		UpdatedBy: actorID,
	}

	return cUsecase.cRepo.UpdateCategoryByID(id, category)
}

func (cUsecase *categoryUsecase) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	err := cUsecase.cRepo.DeleteCategoryByID(ctx, id, actorID)
	if err != nil {
		return err
	}

	if metaValid {
		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "categories",
			EntityID:  id.String(),
			OldValues: "{}",
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData model.AuditLog) {
			cUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}
	return nil
}
