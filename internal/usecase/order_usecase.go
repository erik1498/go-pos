package usecase

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type orderUsecase struct {
	oRepo domain.OrderRepository
	mRepo domain.MemberRepository
	pRepo domain.ProductRepository
}

func GetOrderUsecase(oRepo domain.OrderRepository, mRepo domain.MemberRepository, pRepo domain.ProductRepository) domain.OrderUsecase {
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

func (oUsecase *orderUsecase) Create(req model.CreateOrderRequest) (model.Order, error) {
	orderID := uuid.Must(uuid.NewV7())
	orderNo := fmt.Sprintf("ORD-%s-%s", time.Now().Format("060102"), strings.ToUpper(orderID.String()[:4]))
	order := model.Order{
		ID:            orderID,
		OrderNo:       orderNo,
		MemberID:      req.MemberID,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: model.PaymentStatusPending,
		Items:         []model.OrderItem{},
	}

	totalQty := decimal.NewFromInt(0)
	subTotal := decimal.NewFromInt(0)
	totalTax := decimal.NewFromInt(0)

	for _, reqItem := range req.Items {
		product, err := oUsecase.pRepo.GetByIDWithTaxes(reqItem.ProductID)
		if err != nil {
			return model.Order{}, domain.ErrProductNotFound
		}

		if product.Stock.LessThan(reqItem.Qty) {
			return model.Order{}, domain.ErrProductStockIsNotEnough
		}

		itemSubTotal := product.Price.Mul(reqItem.Qty)
		orderItemID := uuid.Must(uuid.NewV7())

		orderItem := model.OrderItem{
			ID:           orderItemID,
			OrderID:      orderID,
			ProductID:    reqItem.ProductID,
			ProductName:  product.Name,
			BasePrice:    product.Price,
			Qty:          reqItem.Qty,
			SubTotal:     itemSubTotal,
			AppliedTaxes: []model.OrderItemTax{},
		}

		itemTotalTax := decimal.NewFromInt(0)
		for _, tax := range product.Taxes {
			if !tax.IsActive {
				continue
			}

			rateDecimal := tax.Rate.Div(decimal.NewFromInt(100))
			taxAmount := itemSubTotal.Mul(rateDecimal)

			taxID := tax.ID

			orderItemTax := model.OrderItemTax{
				ID:          taxID,
				OrderItemID: orderItemID,
				TaxID:       &taxID,
				TaxName:     tax.Nama,
				TaxRate:     tax.Rate,
				TaxAmount:   taxAmount,
			}

			orderItem.AppliedTaxes = append(orderItem.AppliedTaxes, orderItemTax)

			itemTotalTax = itemTotalTax.Add(taxAmount)
		}

		order.Items = append(order.Items, orderItem)

		totalQty = totalQty.Add(reqItem.Qty)
		subTotal = subTotal.Add(itemSubTotal)
		totalTax = totalTax.Add(itemTotalTax)
	}

	order.TotalQty = totalQty
	order.SubTotal = subTotal
	order.TotalTax = totalTax
	order.GrandTotal = subTotal.Add(totalTax)

	if order.PaymentMethod == model.PaymentMethodCash {
		order.PaymentStatus = model.PaymentStatusPaid
		now := time.Now()
		order.PaidAt = &now
	}

	order, err := oUsecase.oRepo.Create(order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
