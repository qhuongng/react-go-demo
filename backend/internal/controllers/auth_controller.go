package controllers

import (
	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/domain/models"
	"chi-mysql-boilerplate/internal/services"
	"chi-mysql-boilerplate/internal/utils/helpers"
	"chi-mysql-boilerplate/internal/utils/jwt"
	"chi-mysql-boilerplate/internal/utils/validators"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validators.Validator
}

func NewAuthHandler(db *sql.DB, validator *validators.Validator) *AuthHandler {
	return &AuthHandler{authService: services.NewAuthService(db), validator: validator}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := handler.validator.BindJSONAndValidate(w, r, &req); err != nil {
		// error is already handled in the validator
		return
	}

	if err := handler.authService.Register(req); err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	message := "User registered successfully"
	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&message))
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := handler.validator.BindJSONAndValidate(w, r, &req); err != nil {
		// error is already handled in the validator
		return
	}

	res, err := handler.authService.Login(req)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.ErrUserDoesNotExist {
			// user hasn't registered yet
			helpers.MessageLogs.ErrorLog.Println(err)
			helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				}))
			return
		} else {
			// handle other errors
			helpers.MessageLogs.ErrorLog.Println(err)
			helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.InternalServerError,
				}))
			return
		}

	}

	// generate access token and include it in the response
	accessToken := GenerateToken(w, map[string]interface{}{"id": res.ID}, false)
	if accessToken == "" {
		return
	}
	res.AccessToken = accessToken

	// generate refresh token and set it as a cookie
	refreshToken := GenerateToken(w, map[string]interface{}{"id": res.ID}, true)
	if refreshToken == "" {
		return
	}

	// set refresh token in the database
	if err = handler.authService.UpdateRefreshToken(res.ID, refreshToken); err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: fmt.Sprintf("Failed to save refresh token to database: %s", err.Error()),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&res))
}

func (handler *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	userId := handler.HandleRefreshToken(w, r)

	if userId == 0 {
		// error, already handled in the helper function
		return
	}
	// generate a new access token and include it in the response
	accessToken := GenerateToken(w, map[string]interface{}{"id": userId}, false)
	if accessToken == "" {
		return
	}
	res := map[string]interface{}{"id": userId, "accessToken": accessToken}
	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&res))
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := GetRefreshTokenFromContext(w, r)

	if refreshToken == "" {
		// error, already handled in the helper function
		return
	}

	// delete the refresh token from the database
	if err := handler.authService.RemoveRefreshToken(refreshToken); err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: fmt.Sprintf("Failed to delete refresh token from database: %s", err.Error()),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	// remove the cookie that contains the refresh token
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	message := "User logged out successfully"
	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&message))
}

// helper function to generate tokens and set them as cookies
func GenerateToken(w http.ResponseWriter, payload map[string]interface{}, isRefreshToken bool) string {
	var tokenDuration time.Duration
	if isRefreshToken {
		tokenDuration = httpcommon.JwtConstants.RefreshTokenDuration
	} else {
		tokenDuration = httpcommon.JwtConstants.AccessTokenDuration
	}

	token, err := jwt.GenerateToken(tokenDuration, payload, isRefreshToken)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: fmt.Sprintf("Failed to generate token: %s", err.Error()),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return ""
	}

	// if the token is a refresh token, set it as a cookie
	if isRefreshToken {
		cookie := http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			MaxAge:   httpcommon.JwtConstants.CookieMaxAge,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}

	return token
}

func GetRefreshTokenFromContext(w http.ResponseWriter, r *http.Request) string {
	refreshToken, ok := r.Context().Value(httpcommon.ContextKeyConstants.RefreshToken).(string)
	if !ok {
		helpers.MessageLogs.ErrorLog.Println("Failed to retrieve refresh token from context")
		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return ""
	}

	return refreshToken
}

// helper function to handle every error that the refresh token might possibly have
// and then spit out the claims' payload (here it's just the user ID)
// gosh I hate this thing
func (handler *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) uint64 {
	refreshToken := GetRefreshTokenFromContext(w, r)
	if refreshToken == "" {
		return 0
	}

	// decode the token to get the claims
	refreshTokenClaims, err := jwt.VerifyToken(refreshToken, true)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		var httpErr httpcommon.Error
		if err.Error() == httpcommon.ErrorMessage.TokenExpired {
			// expired token
			httpErr = httpcommon.Error{
				Message: httpcommon.ErrorMessage.TokenExpired,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}
		} else {
			httpErr = httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}
		}

		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(httpErr))
		return 0
	}

	// get the payload from the claims (here it's just the user ID)
	tokenPayload, ok := refreshTokenClaims.Payload.(map[string]interface{})
	if !ok {
		helpers.MessageLogs.ErrorLog.Println("Failed to get payload from refresh token")
		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return 0
	}

	userId := uint64(tokenPayload["id"].(float64))

	// check if the token retrieved from the database is the same as the one in the request
	isValid, err := handler.authService.ValidateRefreshToken(userId, refreshToken)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println("Failed to get payload from refresh token")
		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return 0
	}
	if !isValid {
		helpers.MessageLogs.ErrorLog.Println("Invalid refresh token")
		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return 0
	}

	return userId
}
