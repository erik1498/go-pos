package usecase

import (
	"context"
	"encoding/json"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log"

	"github.com/google/uuid"
)

type categoryUsecase struct {
	aRepo domain.AuditLogRepository
	cRepo domain.CategoryRepository
}

func GetCategoryUsecase(aRepo domain.AuditLogRepository, cRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		aRepo: aRepo,
		cRepo: cRepo,
	}
}

func (cUsecase *categoryUsecase) GetAll(ctx context.Context, opts domain.QueryOptions) ([]model.Category, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return cUsecase.cRepo.GetAll(ctx, cleanOpts)
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

	category, err := cUsecase.cRepo.Create(ctx, category)
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
			if err := cUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return category, nil
}

func (cUsecase *categoryUsecase) GetByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	return cUsecase.cRepo.GetByID(ctx, id)
}

func (cUsecase *categoryUsecase) UpdateCategoryByID(ctx context.Context, id uuid.UUID, req model.CategoryRequest) (model.Category, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	oldCategory, err := cUsecase.cRepo.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, err
	}

	category := model.Category{
		ID:        id,
		Name:      req.Name,
		UpdatedBy: actorID,
	}

	updatedCategory, err := cUsecase.cRepo.UpdateCategoryByID(ctx, id, category)
	if err != nil {
		return model.Category{}, err
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldCategory)
		newValuesJSON, _ := json.Marshal(updatedCategory)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "UPDATE",
			Entity:    "categories",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			if err := cUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return updatedCategory, nil
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
