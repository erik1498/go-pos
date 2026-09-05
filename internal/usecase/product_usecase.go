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

type productUsecase struct {
	pRepo domain.ProductRepository
	cRepo domain.CategoryRepository
	tRepo domain.TaxRepository
	aRepo domain.AuditLogRepository
}

func GetProductUsecase(aRepo domain.AuditLogRepository, pRepo domain.ProductRepository, cRepo domain.CategoryRepository, tRepo domain.TaxRepository) domain.ProductUsecase {
	return &productUsecase{
		aRepo: aRepo,
		pRepo: pRepo,
		cRepo: cRepo,
		tRepo: tRepo,
	}
}

func (pUsecase *productUsecase) GetAll(ctx context.Context, opts domain.QueryOptions) ([]domain.Product, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
		"sku":  true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	products, totalItems, err := pUsecase.pRepo.GetAll(ctx, cleanOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("[usecase][product_usecase][GetAll] failed fetch from repo: %w", err)
	}

	return products, totalItems, nil
}

func (pUsecase *productUsecase) Create(ctx context.Context, req domain.CreateProductParam) (domain.Product, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][Create] idempotency key required: %w", domain.ErrIdempotencyRequired)
	}

	category, err := pUsecase.cRepo.GetByID(ctx, uuid.MustParse(req.CategoryID.String()))
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][Create] failed fetch from repo: %w", err)
	}

	product := domain.Product{
		ID:             uuid.Must(uuid.NewV7()),
		CategoryID:     category.ID,
		Category:       &category,
		Name:           req.Name,
		SKU:            req.SKU,
		Price:          req.Price,
		IdempotencyKey: idemKey,
		CreatedBy:      actorID,
		UpdatedBy:      actorID,
	}

	var taxes []domain.Tax
	for _, item := range req.Taxes {
		tax, err := pUsecase.tRepo.GetByID(ctx, item.ID)
		if err != nil {
			return domain.Product{}, fmt.Errorf("[usecase][product_usecase][Create] failed fetch from repo: %w", err)
		}

		taxes = append(taxes, tax)
	}

	product.Taxes = taxes

	product, err = pUsecase.pRepo.Create(ctx, product)
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][Create] failed create from repo: %w", err)
	}

	if metaValid {
		newValuesJSON, errMarshal := json.Marshal(category)

		if errMarshal != nil {
			slog.Warn("[usecase][product_usecase][Create] failed to marshal new values",
				slog.String("error_trace", errMarshal.Error()),
			)
			newValuesJSON = []byte("{}")
		}

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "products",
			EntityID:  category.ID.String(),
			OldValues: "null",
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := pUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][product_usecase][Create] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}

	return product, nil
}

func (pUsecase *productUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	product, err := pUsecase.pRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][GetByID] failed fetch from repo: %w", err)
	}

	return product, nil
}

func (pUsecase *productUsecase) UpdateByID(ctx context.Context, id uuid.UUID, req domain.UpdateProductParam) (domain.Product, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	category, err := pUsecase.cRepo.GetByID(ctx, uuid.MustParse(req.CategoryID.String()))
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][UpdateByID] failed fetch from repo: %w", err)
	}

	oldProduct, err := pUsecase.pRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][UpdateByID] failed fetch from repo: %w", err)
	}

	product := domain.Product{
		CategoryID: category.ID,
		Category:   &category,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		Taxes:      []domain.Tax{},
	}

	var taxes []domain.Tax
	for _, item := range req.Taxes {
		tax, err := pUsecase.tRepo.GetByID(ctx, item.ID)
		if err != nil {
			return domain.Product{}, fmt.Errorf("[usecase][product_usecase][UpdateByID] failed fetch from repo: %w", err)
		}
		taxes = append(taxes, tax)
	}

	product.Taxes = taxes

	updatedProduct, err := pUsecase.pRepo.UpdateByID(ctx, id, product)
	if err != nil {
		return domain.Product{}, fmt.Errorf("[usecase][product_usecase][UpdateByID] failed update from repo: %w", err)
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldProduct)
		newValuesJSON, errMarshal := json.Marshal(updatedProduct)

		if errMarshal != nil {
			slog.Warn("[usecase][product_usecase][Update] failed to marshal new values",
				slog.String("error_trace", errMarshal.Error()),
			)
			newValuesJSON = []byte("{}")
		}

		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "UPDATE",
			Entity:    "products",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: string(newValuesJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData domain.AuditLog) {
			if err := pUsecase.aRepo.Create(context.Background(), logData); err != nil {
				slog.Error("[usecase][product_usecase][UpdateByID] failed to record audit log",
					slog.String("error_trace", err.Error()),
					slog.String("entity_id", logData.EntityID),
					slog.String("actor_id", logData.ActorID.String()),
				)
			}
		}(auditLog)
	}

	return updatedProduct, nil
}

func (pUsecase *productUsecase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	oldProduct, err := pUsecase.pRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("[usecase][category_usecase][DeleteByID] failed fetch from repo: %w", err)
	}

	if err := pUsecase.pRepo.DeleteByID(ctx, id, actorID); err != nil {
		return fmt.Errorf("[usecase][product_usecase][DeleteByID] failed fetch from repo: %w", err)
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldProduct)
		auditLog := domain.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "products",
			EntityID:  id.String(),
			OldValues: string(oldValuesJSON),
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData domain.AuditLog) {
			slog.Error("[usecase][category_usecase][DeleteByID] failed to record audit log",
				slog.String("error_trace", err.Error()),
				slog.String("entity_id", logData.EntityID),
				slog.String("actor_id", logData.ActorID.String()),
			)
		}(auditLog)
	}
	return nil
}
