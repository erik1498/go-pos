package rest

import (
	"context"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	uUsecase domain.UserUsecase
}

func GetAuthMiddleware(uUsecase domain.UserUsecase) *authMiddleware {
	return &authMiddleware{
		uUsecase: uUsecase,
	}
}

func getUserSessionData(r *http.Request) (model.UserSession, error) {
	authString := r.Header.Get("Authorization")
	splitString := strings.Split(authString, " ")
	if len(splitString) != 2 {
		return model.UserSession{}, domain.ErrUnauthorized
	}

	return model.UserSession{
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

		userID, role, err := am.uUsecase.CheckSession(userSession)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Internal: err,
				Message:  domain.ErrSessionExpired.Error(),
			}
		}

		authContext := context.WithValue(c.Request().Context(), model.AuthContextKey, userID)
		c.SetRequest(c.Request().WithContext(authContext))
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
