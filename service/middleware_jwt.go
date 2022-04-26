package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MapClaims map[string]interface{}

type JWTMiddleware struct {
	SigningAlgorithm string
	PubKeyBytes      []byte
	Unauthorized     func(*gin.Context, int, error)
	SetUserIdentity  func(*gin.Context, MapClaims)
}

func (mw *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *JWTMiddleware) middlewareImpl(c *gin.Context) {
	claims, err := mw.GetClaimsFromJWT(c)
	if err != nil {
		mw.Unauthorized(c, http.StatusUnauthorized, err)
		return
	}

	if claims["exp"] == nil {
		mw.Unauthorized(c, http.StatusBadRequest, errors.New("TokenExpiredException"))
		return
	}

	if _, ok := claims["exp"].(float64); !ok {
		mw.Unauthorized(c, http.StatusBadRequest, errors.New("TokenExpiredException"))
		return
	}

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		mw.Unauthorized(c, http.StatusUnauthorized, errors.New("TokenExpiredException"))
		return
	}

	c.Set("JWT_PAYLOAD", claims)

	if mw.SetUserIdentity != nil {
		mw.SetUserIdentity(c, claims)
	}

	c.Next()
}

func (mw *JWTMiddleware) GetClaimsFromJWT(c *gin.Context) (MapClaims, error) {
	token, err := mw.ParseToken(c)

	if err != nil {
		return nil, err
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

func (mw *JWTMiddleware) ParseToken(c *gin.Context) (*jwt.Token, error) {
	token, err := mw.jwtFromHeader(c)

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, errors.New("InvalidTokenException")
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(mw.PubKeyBytes)
		if err != nil {
			return nil, errors.New("InvalidPublicKeyException")
		}
		return key, nil
	})
}

func (mw *JWTMiddleware) jwtFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("InvalidTokenException")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("InvalidTokenException")
	}

	return parts[1], nil
}
