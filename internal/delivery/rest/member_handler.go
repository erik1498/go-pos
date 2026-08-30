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

	ctx := c.Request().Context()

	members, totalItems, err := h.mUsecase.GetAll(ctx, opts)

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

	ctx := c.Request().Context()

	member, err := h.mUsecase.Create(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrIdempotencyKeyDuplicate) {
			return response.ErrConflictRequest(c, domain.ErrIdempotencyKeyDuplicate.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, member)
}

func (h *handler) GetMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	ctx := c.Request().Context()

	member, err := h.mUsecase.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return response.ErrNotFound(c, domain.ErrMemberNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
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

	ctx := c.Request().Context()

	member, err := h.mUsecase.UpdateByID(ctx, ID, req)
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

	ctx := c.Request().Context()

	err = h.mUsecase.DeleteByID(ctx, ID)
	if err != nil {
		if errors.Is(err, domain.ErrMemberNotFound) {
			return response.ErrNotFound(c, domain.ErrMemberNotFound.Error())
		}
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}
