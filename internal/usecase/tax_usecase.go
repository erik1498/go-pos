package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/middleware"
	"go-pos/pkg/utils"
	"log/slog"

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
	allowedFields := map[string]bool{
		"name": true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	taxes, totalItems, err := tUsecase.tRepo.GetAll(ctx, cleanOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("[usecase][tax_usecase][GetAll] failed fetch from repo: %w", err)
	}

	return taxes, totalItems, nil
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

	tax, err := tUsecase.tRepo.Create(ctx, domain.Tax{
		ID:             uuid.Must(uuid.NewV7()),
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
		newValuesJSON, errMarshal := json.Marshal(tax)

		if errMarshal != nil {
			slog.Warn("[usecase][tax_usecase][Create] failed to marshal new values",
				slog.String("error_trace", errMarshal.Error()),
			)
			newValuesJSON = []byte("{}")
		}

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "taxes",
			EntityID:  tax.ID.String(),
			OldValues: "null",
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := tUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][tax_usecase][Create] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
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
		newValuesJSON, errMarshal := json.Marshal(updatedTax)

		if errMarshal != nil {
			slog.Warn("[usecase][tax_usecase][Update] failed to marshal new values",
				slog.String("error_trace", errMarshal.Error()),
			)
			newValuesJSON = []byte("{}")
		}

		auditLog := domain.AuditLog{
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

		go func(logData domain.AuditLog) {
			if err := tUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][category_usecase][UpdateByID] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
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

	oldTax, err := tUsecase.tRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("[usecase][tax_usecase][DeletedByID] failed fetch from repo: %w", err)
	}

	err = tUsecase.tRepo.DeleteByID(ctx, id, actorID)
	if err != nil {
		return err
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldTax)
		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "taxes",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData domain.AuditLog) {
			if err := tUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][category_usecase][DeleteByID] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}

	return nil
}
