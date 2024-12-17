package models

type User struct {
	ID           uint64 `db:"id"`
	Username     string `db:"username"`
	Password     string `db:"password"`
	RefreshToken string `db:"refresh_token"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"pwdminlen"`
}

type AuthResponse struct {
	ID          uint64 `json:"id"`
	AccessToken string `json:"accessToken"`
}
