package v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pozedorum/user-balance-service/internal/service"
)

const (
	userIdCtx = "userId"
)

type authMiddleware struct {
	authService service.Auth
}

func (h *authMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())

		if !ok {
			log.Errorf("AuthMiddleware.UserIdentity - bearerToken: %v", ErrInvalidAuthHeader)
			newErrResponce(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			log.Errorf("AuthMiddleware.UserIdentity - h.authService.ParseToken: %v", err)
			newErrResponce(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
		}

		c.Set(userIdCtx, userId)
		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}
	return "", false
}
