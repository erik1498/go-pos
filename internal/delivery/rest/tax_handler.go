package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllTax(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	taxes, totalItems, err := h.tUsecase.GetAll(opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, taxes, meta)
}

func (h *handler) CreateTax(c echo.Context) error {
	var req model.TaxRequest

	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	tax, err := h.tUsecase.Create(req)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, tax)
}

func (h *handler) GetTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	tax, err := h.tUsecase.GetByID(ID)
	if err != nil {
		if errors.Is(err, domain.ErrTaxNotFound) {
			return response.ErrNotFound(c, domain.ErrTaxNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, tax)
}

func (h *handler) UpdateTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.TaxRequest
	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	tax, err := h.tUsecase.UpdateByID(ID, req)
	if err != nil {
		if errors.Is(err, domain.ErrTaxNotFound) {
			return response.ErrNotFound(c, domain.ErrTaxNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, tax)
}

func (h *handler) DeleteTaxByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	if err := h.tUsecase.DeleteByID(ID); err != nil {
		if errors.Is(err, domain.ErrTaxNotFound) {
			return response.ErrNotFound(c, domain.ErrTaxNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}
