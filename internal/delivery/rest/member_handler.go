package rest

import (
	"fmt"
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MemberRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Phone string `json:"phone" validate:"required,max=30"`
	Email string `json:"email" validate:"required"`
}

type MemberResponse struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toMemberResponse(m domain.Member) MemberResponse {
	return MemberResponse{
		Name:      m.Name,
		Phone:     m.Phone,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toMemberResponseList(members []domain.Member) []MemberResponse {
	var res []MemberResponse
	for _, m := range members {
		res = append(res, toMemberResponse(m))
	}
	return res
}

func (h *handler) GetAllMember(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	ctx := c.Request().Context()

	members, totalItems, err := h.mUsecase.GetAll(ctx, opts)

	if err != nil {
		return err
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, toMemberResponseList(members), meta)
}

func (h *handler) CreateMember(c echo.Context) error {
	var req MemberRequest

	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][CreateMember] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][CreateMember] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.CreateMemberParam{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}

	member, err := h.mUsecase.Create(ctx, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, domain.SuccessCreateData, toMemberResponse(member))
}

func (h *handler) GetMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][GetMemberByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	ctx := c.Request().Context()

	member, err := h.mUsecase.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, toMemberResponse(member))
}

func (h *handler) UpdateMemberByID(c echo.Context) error {
	idParam := c.Param("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][UpdateMemberByID] invalid UUID format: %w", domain.ErrIDInvalid)
	}

	var req MemberRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][UpdateMemberByID] invalid body: %w", domain.ErrBadRequest)
	}

	if err := c.Validate(&req); err != nil {
		return fmt.Errorf("[delivery][rest][member_handler][UpdateMemberByID] validation error: %w", domain.ErrBadRequest)
	}

	ctx := c.Request().Context()

	param := domain.UpdateMemberParam{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}

	member, err := h.mUsecase.UpdateByID(ctx, ID, param)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, toMemberResponse(member))
}
