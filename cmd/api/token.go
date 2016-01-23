package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrInvalidToken is the error returned when a token is invalid
	ErrInvalidToken = errors.New("Invalid token")
)

//Token defines parsed token
type Token struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Azp string `json:"azp"`
	Aud string `json:"aud"`
	Iat string `json:"iat"`
	Exp string `json:"exp"`
}

// TokenValidator defines a token validator
type TokenValidator struct {
	clientID string
}

// NewTokenValidator creates a new token validator
func NewTokenValidator(config *Config) *TokenValidator {
	return &TokenValidator{
		clientID: config.Token.ClientID,
	}
}

// ParseToken parses a token string into a token
func (tv *TokenValidator) ParseToken(token string) (*Token, error) {
	var tkn Token

	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrInvalidToken
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&tkn)
	if err != nil {
		return nil, err
	}

	return &tkn, nil
}

// ValidateToken validates token
func (tv *TokenValidator) ValidateToken(token string) (bool, string, error) {
	tkn, err := tv.ParseToken(token)
	if err != nil {
		return false, "", err
	}

	return tkn.Aud == tv.clientID, tkn.Sub, nil
}
