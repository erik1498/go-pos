package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllOrder(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	orders, totalItems, err := h.oUsecase.GetAll(opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, orders, meta)
}

func (h *handler) CreateOrder(c echo.Context) error {
	var req model.CreateOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	order, err := h.oUsecase.Create(req)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNoIsAlreadyRegistered) {
			return response.ErrBadRequest(c, domain.ErrOrderNoIsAlreadyRegistered.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, order)
}
