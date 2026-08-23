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
				Message:  err.Error(),
			}
		}

		userID, err := am.uUsecase.CheckSession(userSession)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Internal: err,
				Message:  err.Error(),
			}
		}

		authContext := context.WithValue(c.Request().Context(), model.AuthContextKey, userID)
		c.SetRequest(c.Request().WithContext(authContext))
		if err := next(c); err != nil {
			return err
		}

		return nil
	}
}
