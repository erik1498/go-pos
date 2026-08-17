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
