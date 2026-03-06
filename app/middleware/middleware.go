package middleware

import (
	"Mini-Project-Game-Vault-API/util/response"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func forbiddenResponse(c echo.Context) error {
	return c.JSON(http.StatusForbidden, response.NewErrorResponse(http.StatusText(http.StatusForbidden)))
}

func JWTMiddleware(jwtSign string) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if strings.Contains(ctx.Request().URL.Path, "/login") {
				return hf(ctx)
			}

			signature := strings.Split(ctx.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return forbiddenResponse(ctx)
			}
			if signature[0] != "Bearer" {
				return forbiddenResponse(ctx)
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(jwtSign), nil
			})

			if err != nil {
				return forbiddenResponse(ctx)
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return forbiddenResponse(ctx)
			}

			expAt, err := claim.GetExpirationTime()
			if err != nil {
				return forbiddenResponse(ctx)
			}

			if time.Now().After(expAt.Time) {
				return forbiddenResponse(ctx)
			}

			userID, _ := claim["id"].(string)
			role, _ := claim["role"].(string)
			ctx.Set("id", userID)
			ctx.Set("role", role)

			return hf(ctx)
		}
	}
}

func ACLMiddleware(rolesMap map[string]bool) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			role, _ := ctx.Get("role").(string)
			if role == "superadmin" {
				return hf(ctx)
			}

			if rolesMap[role] {
				return hf(ctx)
			}

			return forbiddenResponse(ctx)
		}
	}
}

func JWTEchoMiddleware(jwtSign string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSign),
	})
}
