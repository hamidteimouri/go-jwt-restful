package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hamidteimouri/go-jwt-restful/models"
	utils "github.com/hamidteimouri/go-jwt-restful/utils"
	"net/http"
	"os"
	"strings"
)

var JwtAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		/* list of url that does not need to authorization */

		noNeedAuthUrl := []string{
			"/api/auth/register",
			"/api/auth/login",
		}

		/* getting current url path */
		requestPath := request.URL.Path

		/*
		 * checking current url path need auth or not
		 * if not, we just pass the auth and continuing serving the request
		 */
		for _, value := range noNeedAuthUrl {
			if value == requestPath {
				next.ServeHTTP(writer, request)
				return
			}
		}

		/* getting the auth token */
		tokenHeader := request.Header.Get("Authorization")

		/* check whether token auth exist */
		if tokenHeader == "" {
			utils.ERROR(writer, http.StatusForbidden, errors.New("invalid auth token"))
		}

		/* checking token format */
		if len(strings.Split(tokenHeader, " ")) != 2 {
			utils.ERROR(writer, http.StatusForbidden, errors.New("Invalid/Malformed auth token"))
		}

		/* getting token */
		tokenHeader = strings.Split(tokenHeader, " ")[1]

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			utils.ERROR(writer, http.StatusForbidden, errors.New("malformed auth token"))
		}

		if !token.Valid {
			utils.ERROR(writer, http.StatusForbidden, errors.New("invalid auth token"))
		}

		/* writing the token's user into context of current request */
		fmt.Sprintf("User %s", tk.UserId)
		ctx := context.WithValue(request.Context(), "user", tk.UserId)
		request = request.WithContext(ctx)
		next.ServeHTTP(writer, request)

	})
}
