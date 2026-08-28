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

func (pUsecase *productUsecase) GetAll(ctx context.Context, opts domain.QueryOptions) ([]model.Product, int64, error) {
	allowedFields := map[string]bool{
		"name": true,
		"sku":  true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return pUsecase.pRepo.GetAll(ctx, cleanOpts)
}

func (pUsecase *productUsecase) Create(ctx context.Context, req model.ProductRequest) (model.Product, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	idemKey, ok := ctx.Value(middleware.IdempotencyKeyCtx).(string)
	if !ok && idemKey == "" {
		return model.Product{}, domain.ErrIdempotencyKeyDuplicate
	}

	category, err := pUsecase.cRepo.GetByID(ctx, uuid.MustParse(req.CategoryID))
	if err != nil {
		return model.Product{}, err
	}

	product := model.Product{
		ID:             uuid.Must(uuid.NewV7()),
		CategoryID:     category.ID,
		Name:           req.Name,
		SKU:            req.SKU,
		Price:          req.Price,
		IdempotencyKey: idemKey,
		CreatedBy:      actorID,
		UpdatedBy:      actorID,
	}

	var taxes []model.Tax
	for _, item := range req.Tax {
		tax, err := pUsecase.tRepo.GetByID(ctx, item.ID)
		if err != nil {
			return model.Product{}, err
		}

		taxes = append(taxes, tax)
	}

	product.Taxes = taxes

	product, err = pUsecase.pRepo.Create(ctx, product)
	if err != nil {
		return model.Product{}, err
	}

	if metaValid {
		newValueJSON, _ := json.Marshal(category)

		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "CREATE",
			Entity:    "products",
			EntityID:  category.ID.String(),
			OldValues: "null",
			NewValues: string(newValueJSON),
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}

		go func(logData model.AuditLog) {
			if err := pUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
			}
		}(auditLog)
	}

	return product, nil
}

func (pUsecase *productUsecase) GetByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	product, err := pUsecase.pRepo.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (pUsecase *productUsecase) UpdateByID(ctx context.Context, id uuid.UUID, req model.ProductRequest) (model.Product, error) {
	meta, metaValid := ctx.Value(middleware.AuditMetaKey).(middleware.AuditMeta)
	var actorID uuid.UUID
	if metaValid {
		actorID = uuid.MustParse(meta.UserID)
	}

	category, err := pUsecase.cRepo.GetByID(ctx, uuid.MustParse(req.CategoryID))
	if err != nil {
		return model.Product{}, err
	}

	oldProduct, err := pUsecase.pRepo.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, err
	}

	product := model.Product{
		CategoryID: category.ID,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		Taxes:      []model.Tax{},
	}

	var taxes []model.Tax
	for _, item := range req.Tax {
		tax, err := pUsecase.tRepo.GetByID(ctx, item.ID)
		if err != nil {
			return model.Product{}, domain.ErrTaxNotFound
		}
		taxes = append(taxes, tax)
	}

	product.Taxes = taxes

	updatedProduct, err := pUsecase.pRepo.UpdateByID(ctx, id, product)
	if err != nil {
		return model.Product{}, err
	}

	if metaValid {
		oldValuesJSON, _ := json.Marshal(oldProduct)
		newValuesJSON, _ := json.Marshal(updatedProduct)

		auditLog := model.AuditLog{
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

		go func(logData model.AuditLog) {
			if err := pUsecase.aRepo.Create(context.Background(), logData); err != nil {
				log.Printf("AUDIT LOG: RECORD AUDIT LOG FAILED, ERR : %v", err)
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

	err := pUsecase.pRepo.DeleteByID(ctx, id, actorID)

	if err != nil {
		return err
	}

	if metaValid {
		auditLog := model.AuditLog{
			ActorID:   actorID,
			ActorRole: meta.Role,
			Action:    "DELETE",
			Entity:    "products",
			EntityID:  id.String(),
			OldValues: "{}",
			NewValues: "null",
			IPAddress: meta.IPAddress,
			UserAgent: meta.UserAgent,
		}
		go func(logData model.AuditLog) {
			pUsecase.aRepo.Create(context.Background(), logData)
		}(auditLog)
	}
	return nil
}
