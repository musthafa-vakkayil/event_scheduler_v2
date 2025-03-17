package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a paseto token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewJWTMaker creates a new JWTMaker
func NewPasetoMaker(symmetrickEy string) (Maker, error) {
	if len(symmetrickEy) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invlaid key size: must be atleast %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetrickEy),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		Subject:    username,
		IssuedAt:   now,
		Expiration: exp,
	}

	jsonToken.Set("id", uuid.NewString())

	return paseto.NewV2().Encrypt(maker.symmetricKey, jsonToken, nil)
}

// VerifyToken checks if a token is valid or not
func (maker PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	var jsonToken paseto.JSONToken

	err := paseto.NewV2().Decrypt(tokenString, maker.symmetricKey, &jsonToken, nil)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	expTime := jsonToken.Expiration
	if time.Now().After(expTime) {
		return nil, ErrTokenExpired
	}

	// Extract claims into the Payload struct
	tokenID, _ := uuid.Parse(jsonToken.Get("id"))

	payload := &Payload{
		ID:       tokenID,
		Username: jsonToken.Subject,
	}

	return payload, nil
}
