package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"loverly/lib/log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	TokenType = "Bearer"

	//AccessTypeOffline for server to server, require client-id & client-secret validation, no refresh_token expiry, not linked to access_token
	AccessTypeOffline string = "offline"
	//AccessTypeOnline  for client to server, doesn't require client-id & client-secret validation, refresh_token has expiry, and linked to access_token
	AccessTypeOnline string = "online"
)

type (
	AccessToken struct {
		jwt.RegisteredClaims
		Scopes string               `json:"scopes"`
		Roles  string               `json:"roles"`
		Data   AccessTokenClaimData `json:"dat"`
	}

	AccessTokenClaimData struct {
		UserId int64 `json:"user_id"`
	}

	RefreshToken struct {
		jwt.RegisteredClaims
		Data RefreshTokenClaimData `json:"dat"`
	}

	RefreshTokenClaimData struct {
		AccessTokenId string `json:"token_id,omitempty"`
		AccessType    string `json:"access_type"`
	}
)

type TokenProvider struct {
	cfg       *Configuration
	log       log.Interface
	signKey   rsa.PrivateKey
	verifyKey rsa.PublicKey
}

type Configuration struct {
	AccessTokenValidity  time.Duration
	RefreshTokenValidity time.Duration
	IatLeeway            time.Duration
	TokenIssuer          string
	KeyId                string
	VerifyKey            string
	SignKey              string
}

func Init(ctx context.Context, cfg *Configuration, log log.Interface) *TokenProvider {
	return &TokenProvider{
		cfg:       cfg,
		log:       log,
		signKey:   *initAccessTokenPrivateKey(ctx, cfg, log),
		verifyKey: *initAccessTokenPublicKey(ctx, cfg, log),
	}
}

/*
Create new accessToken for given user and identity
*/
func (t TokenProvider) NewAccessToken(ctx context.Context, userId int64, audiences []string, accessType string) (*oauth2.Token, error) {
	if accessType != AccessTypeOffline && accessType != AccessTypeOnline {
		return nil, fmt.Errorf("invalid_access_type")
	}

	accessToken, err := t.newAccessToken(ctx, userId, audiences)
	if err != nil {
		return nil, err
	}

	strAccessToken, err := t.encodeAccessToken(ctx, *accessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.newRefreshToken(ctx, *accessToken, accessType)
	if err != nil {
		return nil, err
	}

	strRefreshToken, err := t.encodeRefreshToken(ctx, *refreshToken)
	if err != nil {
		return nil, err
	}

	token := oauth2.Token{
		AccessToken:  strAccessToken,
		RefreshToken: strRefreshToken,
		TokenType:    TokenType,
		Expiry:       accessToken.ExpiresAt.Time,
	}

	return &token, nil
}

func (t TokenProvider) newAccessToken(ctx context.Context, userId int64, audiences []string) (*AccessToken, error) {
	accessToken := AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.cfg.AccessTokenValidity)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(t.cfg.IatLeeway * -1)),
			Issuer:    t.cfg.TokenIssuer,
			Audience:  audiences,
			Subject:   fmt.Sprintf("%d", userId),
		},
		Scopes: "*",
		Roles:  "user",
		Data: AccessTokenClaimData{
			UserId: userId,
		},
	}
	return &accessToken, nil
}

/*
Create new refreshToken for given accessToken
*/
func (t TokenProvider) newRefreshToken(ctx context.Context, accessToken AccessToken, accessType string) (*RefreshToken, error) {
	refreshToken := RefreshToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    accessToken.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(t.cfg.IatLeeway * -1)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.cfg.RefreshTokenValidity)),
			Audience:  accessToken.Audience,
			Subject:   accessToken.Subject,
		},
		Data: RefreshTokenClaimData{
			AccessTokenId: accessToken.ID,
			AccessType:    AccessTypeOnline,
		},
	}
	if accessType == AccessTypeOffline {
		refreshToken.Data.AccessTokenId = ""
		refreshToken.Data.AccessType = AccessTypeOffline
		refreshToken.RegisteredClaims.ExpiresAt = nil
	}
	return &refreshToken, nil
}

/*
Encode given accessToken as string
*/
func (t TokenProvider) encodeAccessToken(ctx context.Context, accessToken AccessToken) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessToken)
	jwtToken.Header["kid"] = t.cfg.KeyId
	jwtTokenString, err := jwtToken.SignedString(&t.signKey)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("encodeAccessToken err: %v", err))
		return "", err
	}

	return jwtTokenString, nil
}

/*
Encode given refreshToken as string
*/
func (t TokenProvider) encodeRefreshToken(ctx context.Context, refreshToken RefreshToken) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshToken)
	tokenString, err := jwtToken.SignedString(&t.signKey)
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("encodeRefreshToken err: %v", err))
		return "", err
	}

	return tokenString, err
}

/*
Decode and validate accessToken, returned decoded AccessToken object on success
*/
func (t TokenProvider) DecodeAccessToken(ctx context.Context, accessToken string) (*AccessToken, error) {
	var validated AccessToken
	_, err := jwt.ParseWithClaims(accessToken, &validated, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		return &t.verifyKey, nil
	})
	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("parse err: %v", err))
		return nil, err
	}

	return &validated, nil
}

/*
Decode and validate refreshToken, returned decoded RefreshToken object on success
*/
func (t TokenProvider) DecodeRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error) {
	var validated RefreshToken
	_, err := jwt.ParseWithClaims(refreshToken, &validated, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		return &t.verifyKey, nil
	})

	if err != nil {
		t.log.Error(ctx, fmt.Sprintf("parse err: %v", err))
		return nil, err
	}

	return &validated, err
}

func initAccessTokenPublicKey(ctx context.Context, cfg *Configuration, log log.Interface) *rsa.PublicKey {
	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.VerifyKey))
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("init acc token public key err: %v", err))
	}
	return rsaPublicKey
}

func initAccessTokenPrivateKey(ctx context.Context, cfg *Configuration, log log.Interface) *rsa.PrivateKey {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.SignKey))
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("init acc token private key err: %v", err))
	}
	return rsaPrivateKey
}
