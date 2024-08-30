package middleware

import "github.com/labstack/echo"

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (a *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Request().Header.Get("Authorization") == "" {
			return ctx.JSON(401, "Unauthorized")
		} else {
			return next(ctx)
		}
	}
}
