package middleware

import (
	"context"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ContextKey string

const AuditMetaKey ContextKey = "audit_meta"

type AuditMeta struct {
	UserID    string
	Role      string
	IPAddress string
	UserAgent string
}

type authMiddleware struct {
	uUsecase domain.UserUsecase
}

func GetAuthMiddleware(uUsecase domain.UserUsecase) *authMiddleware {
	return &authMiddleware{
		uUsecase: uUsecase,
	}
}

func getUserSessionData(r *http.Request) (domain.UserSession, error) {
	authString := r.Header.Get("Authorization")
	splitString := strings.Split(authString, " ")
	if len(splitString) != 2 {
		return domain.UserSession{}, domain.ErrUnauthorized
	}

	return domain.UserSession{
		AccessToken: splitString[1],
	}, nil
}

func (am *authMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userSession, err := getUserSessionData(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Internal: err,
				Message:  domain.ErrSessionExpired.Error(),
			}
		}

		ctx := c.Request().Context()

		userID, role, err := am.uUsecase.CheckSession(ctx, userSession)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Internal: err,
				Message:  domain.ErrSessionExpired.Error(),
			}
		}
		ipAddr := c.RealIP()
		userAgent := c.Request().UserAgent()

		auditMeta := AuditMeta{
			UserID:    userID,
			Role:      role,
			IPAddress: ipAddr,
			UserAgent: userAgent,
		}

		authContext := context.WithValue(ctx, model.AuthContextKey, userID)
		auditContext := context.WithValue(authContext, AuditMetaKey, auditMeta)
		c.SetRequest(c.Request().WithContext(auditContext))

		c.Set("userID", userID)
		c.Set("userRole", role)

		if err := next(c); err != nil {
			return err
		}

		return nil
	}
}

func (am *authMiddleware) RequiredRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("userRole").(string)
			if !ok || userRole == "" {
				return &echo.HTTPError{
					Code:    http.StatusForbidden,
					Message: domain.ErrForbidden,
				}
			}

			isAllowed := false
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				return &echo.HTTPError{
					Code:    http.StatusForbidden,
					Message: domain.ErrForbidden,
				}
			}

			return next(c)
		}
	}
}
