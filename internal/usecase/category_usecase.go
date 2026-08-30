package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log/slog"

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

	categories, totalItems, err := cUsecase.cRepo.GetAll(ctx, cleanOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("[usecase][category_usecase][GetAll] failed fetch from repo: %w", err)
	}

	return categories, totalItems, nil
}

func (cUsecase *categoryUsecase) Create(ctx context.Context, req model.CategoryRequest) (model.Category, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok || idemKey == "" {
		return model.Category{}, fmt.Errorf("[usecase][category][Create] idempotency key required: %w", domain.ErrIdempotencyRequired)
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
		return model.Category{}, fmt.Errorf("[usecase][category_usecase][Create] failed create from repo: %w", err)
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
				slog.Error("[usecase][category_usecase][Create] failed to record audit log",
					slog.String("error", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}

	return category, nil
}

func (cUsecase *categoryUsecase) GetByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	category, err := cUsecase.cRepo.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, fmt.Errorf("[usecase][category_usecase][GetByID] failed fetch from repo: %w", err)
	}
	return category, nil
}

func (cUsecase *categoryUsecase) UpdateCategoryByID(ctx context.Context, id uuid.UUID, req model.CategoryRequest) (model.Category, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	oldCategory, err := cUsecase.cRepo.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, fmt.Errorf("[usecase][category_usecase][UpdateCategoryByID] failed fetch from repo: %w", err)
	}

	category := model.Category{
		ID:        id,
		Name:      req.Name,
		UpdatedBy: actorID,
	}

	updatedCategory, err := cUsecase.cRepo.UpdateCategoryByID(ctx, id, category)
	if err != nil {
		return model.Category{}, fmt.Errorf("[usecase][category_usecase][UpdateCategoryByID] failed update from repo: %w", err)
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
				slog.Error("[usecase][category_usecase][UpdateCategoryByID] failed to record audit log",
					slog.String("error", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
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

	oldCategory, err := cUsecase.cRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("[usecase][category_usecase][DeleteCategoryByID] failed fetch from repo: %w", err)
	}

	err = cUsecase.cRepo.DeleteCategoryByID(ctx, id, actorID)
	if err != nil {
		return fmt.Errorf("[usecase][category_usecase][DeleteCategoryByID] failed delete from repo: %w", err)
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldCategory)
		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "categories",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData model.AuditLog) {
			if err := cUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][category_usecase][DeleteCategoryByID] failed to record audit log",
					slog.String("error", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}
	return nil
}
