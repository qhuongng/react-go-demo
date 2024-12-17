package middleware

import (
	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/utils/helpers"
	"chi-mysql-boilerplate/internal/utils/jwt"
	"context"
	"net/http"
)

func getAccessToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		// remove the "Bearer " prefix if it exists
		const bearerPrefix = "Bearer "
		if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
			token = token[len(bearerPrefix):]
		}
	}

	return token
}

func VerifyAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessToken(r)
		accessTokenClaims, err := jwt.VerifyToken(accessToken, false)
		if err != nil {
			if err.Error() == httpcommon.ErrorMessage.TokenExpired {
				// token expired
				helpers.MessageLogs.ErrorLog.Println(err)
				helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
					httpcommon.Error{
						Message: httpcommon.ErrorMessage.TokenExpired,
						Code:    httpcommon.ErrorResponseCode.Unauthorized,
					}))
				return
			} else {
				// handle other errors (e.g. undefined/empty token)
				helpers.MessageLogs.ErrorLog.Println(err)
				helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
					httpcommon.Error{
						Message: httpcommon.ErrorMessage.BadCredentials,
						Code:    httpcommon.ErrorResponseCode.Unauthorized,
					}))
				return
			}
		}

		// get the user ID from the access token
		// you might think this makes the ID retrieval in AuthHandler.HandleRefreshToken redundant
		// but if the access token expires this ID wouldn't exist in the context at all
		payload, ok := accessTokenClaims.Payload.(map[string]interface{})
		if !ok {
			helpers.MessageLogs.ErrorLog.Println(err)
			helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: httpcommon.ErrorMessage.BadCredentials,
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				}))
			return
		}
		userId := uint64(payload["id"].(float64))

		ctx := context.WithValue(r.Context(), httpcommon.ContextKeyConstants.UserId, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the refresh token from the cookie
		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			helpers.MessageLogs.ErrorLog.Println("Failed to retrieve refresh token from cookie: ", err)
			helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: httpcommon.ErrorMessage.BadCredentials,
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				}))
			return
		}

		// store the ID and refresh token in an HTTP context for handlers to retrieve
		ctx := context.WithValue(r.Context(), httpcommon.ContextKeyConstants.RefreshToken, refreshToken.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
