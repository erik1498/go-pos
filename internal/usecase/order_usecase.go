package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type orderUsecase struct {
	oRepo domain.OrderRepository
	mRepo domain.MemberRepository
}

func GetOrderUsecase(oRepo domain.OrderRepository, mRepo domain.MemberRepository) domain.OrderUsecase {
	return &orderUsecase{
		oRepo: oRepo,
		mRepo: mRepo,
	}
}

func (oUsecase *orderUsecase) GetAll(opts domain.QueryOptions) ([]model.Order, int64, error) {
	allowedSorts := map[string]bool{
		"order_no":   true,
		"created_at": true,
	}

	allowedFilter := map[string]bool{
		"order_no": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFilter, allowedSorts, "created_at")

	return oUsecase.oRepo.GetAll(cleanOpts)
}

func (oUsecase *orderUsecase) Create(req model.OrderRequest) (model.Order, error) {
	member, err := oUsecase.mRepo.GetByID(uuid.MustParse(req.MemberID))
	if err != nil {
		return model.Order{}, err
	}

	order := model.Order{
		OrderNo:       req.OrderNo,
		MemberID:      member.ID,
		PaymentMethod: req.PaymentMethod,
	}

	order, err = oUsecase.oRepo.Create(order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
