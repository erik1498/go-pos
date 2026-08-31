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

type taxUsecase struct {
	aRepo domain.AuditLogRepository
	tRepo domain.TaxRepository
}

func GetTaxUsecase(aRepo domain.AuditLogRepository, tRepo domain.TaxRepository) domain.TaxUsecase {
	return &taxUsecase{
		aRepo: aRepo,
		tRepo: tRepo,
	}
}

func (tUsecase *taxUsecase) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Tax, int64, error) {
	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	allowedFields := map[string]bool{
		"name": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return tUsecase.tRepo.GetAll(ctx, cleanOpts)
}

func (tUsecase *taxUsecase) Create(ctx context.Context, req domain.CreateTaxParam) (domain.Tax, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return domain.Tax{}, domain.ErrIdempotencyKeyDuplicate
	}

	ID := uuid.Must(uuid.NewV7())

	tax, err := tUsecase.tRepo.Create(ctx, domain.Tax{
		ID:             ID,
		Name:           req.Name,
		Rate:           req.Rate,
		IdempotencyKey: idemKey,
		CreatedBy:      actorID,
		UpdatedBy:      actorID,
	})

	if err != nil {
		return domain.Tax{}, err
	}

	if metaValid {
		newValueJSON, _ := json.Marshal(tax)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "taxes",
			EntityID:  tax.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			if err := tUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return tax, nil
}

func (tUsecase *taxUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Tax, error) {
	tax, err := tUsecase.tRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tax{}, err
	}

	return tax, nil
}

func (tUsecase *taxUsecase) UpdateByID(ctx context.Context, id uuid.UUID, req domain.UpdateTaxParam) (domain.Tax, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	oldTax, err := tUsecase.tRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Tax{}, err
	}

	tax := domain.Tax{
		Name:      req.Name,
		Rate:      req.Rate,
		UpdatedBy: actorID,
	}

	updatedTax, err := tUsecase.tRepo.UpdateByID(ctx, id, tax)
	if err != nil {
		return domain.Tax{}, err
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldTax)
		newValuesJSON, _ := json.Marshal(updatedTax)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "UPDATE",
			Entity:    "taxes",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			if err := tUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return updatedTax, err
}

func (tUsecase *taxUsecase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	err := tUsecase.tRepo.DeleteByID(ctx, id, actorID)
	if err != nil {
		return err
	}

	if metaValid {
		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "taxes",
			EntityID:  id.String(),
			OldValues: "{}",
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData model.AuditLog) {
			tUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}

	return nil
}
