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

func (h *handler) GetAllMember(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	members, totalItems, err := h.mUsecase.GetAll(opts)

	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, members, meta)
}

func (h *handler) CreateMember(c echo.Context) error {
	var req model.MemberRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	member, err := h.mUsecase.Create(req)
	if err != nil {
		if errors.Is(err, domain.ErrPhoneNumberRequired) {
			return response.ErrValidation(c, domain.ErrPhoneNumberRequired.Error(), nil)
		}
		return response.ErrInternalServer(c, err.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, member)
}

func (h *handler) GetMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	member, err := h.mUsecase.GetByID(ID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return response.ErrNotFound(c, domain.ErrMemberNotFound.Error())
		}
		return response.ErrInternalServer(c, err.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, member)
}

func (h *handler) UpdateMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.MemberRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	member, err := h.mUsecase.UpdateByID(req, ID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return response.ErrNotFound(c, domain.ErrMemberNotFound.Error())
		}
		return response.ErrBadRequest(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, member)
}

func (h *handler) DeleteMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	err = h.mUsecase.DeleteByID(ID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return response.ErrNotFound(c, domain.ErrMemberNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}
