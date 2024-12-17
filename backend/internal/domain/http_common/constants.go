package http_common

import (
	"os"
	"time"
)

type errorResponseCode struct {
	InvalidRequest      string
	InternalServerError string
	RecordNotFound      string
	MissingIdParameter  string
	InvalidDataType     string
	Unauthorized        string
}

var ErrorResponseCode = errorResponseCode{
	InvalidRequest:      "INVALID_REQUEST",
	InternalServerError: "INTERNAL_SERVER_ERROR",
	RecordNotFound:      "RECORD_NOT_FOUND",
	MissingIdParameter:  "MISSING_ID_PARAMETER",
	InvalidDataType:     "INVALID_DATA_TYPE",
	Unauthorized:        "UNAUTHORIZED",
}

type customValidationErrCode map[string]string

var CustomValidationErrCode = customValidationErrCode{}

type errorMessage struct {
	ErrUserAlreadyExists string
	ErrUserDoesNotExist  string
	InvalidDataType      string
	InvalidRequest       string
	BadCredentials       string
	SilentRefreshFailed  string
	TokenExpired         string
}

var ErrorMessage = errorMessage{
	ErrUserAlreadyExists: "user already exists",
	ErrUserDoesNotExist:  "user does not exist",
	InvalidDataType:      "invalid data type",
	InvalidRequest:       "invalid request",
	BadCredentials:       "bad credentials",
	SilentRefreshFailed:  "silent refresh failed",
	TokenExpired:         "token has invalid claims: token is expired",
}

type jwtConstants struct {
	AccessSecretKey      string
	RefreshSecretKey     string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	CookieMaxAge         int
}

var JwtConstants = jwtConstants{
	AccessSecretKey:      os.Getenv("JWT_ACCESS_SECRET"),
	RefreshSecretKey:     os.Getenv("JWT_REFRESH_SECRET"),
	AccessTokenDuration:  15 * time.Minute,
	RefreshTokenDuration: 24 * time.Hour,
	CookieMaxAge:         int(7 * 24 * time.Hour.Seconds()),
}

type dbConstants struct {
	Timeout time.Duration
}

var DbConstants = dbConstants{
	Timeout: 5 * time.Second,
}

type contextKey string

type contextKeyConstants struct {
	UserId       contextKey
	RefreshToken contextKey
}

var ContextKeyConstants = contextKeyConstants{
	UserId:       contextKey("user_id"),
	RefreshToken: contextKey("refresh_token"),
}
