package models

import (
	"time"
)

// @Description RenewAccessTokenRequest object used for input
type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Description RenewAccessTokenResponse object used for input
type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expiry"`
}
