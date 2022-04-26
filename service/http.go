package main

import (
	"os"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	kafka *KafkaBroker
}

func (h *HttpHandler) ReturnError(err error) gin.H {
	return gin.H{
		"error":             err.Error(),
		"error_description": err.Error(),
	}
}

func (h *HttpHandler) Handle() {
	router := gin.Default()
	// the jwt middleware
	jwtPublicKey := os.Getenv("JWT_PUBLIC_KEY")
	authMiddleware := JWTMiddleware{
		SigningAlgorithm: "RS256",
		PubKeyBytes:      []byte(jwtPublicKey),
		Unauthorized: func(c *gin.Context, code int, err error) {
			c.JSON(code, h.ReturnError(err))
			c.Abort()
		},
		SetUserIdentity: func(c *gin.Context, claims MapClaims) {
			if userid, ok := claims["sub"]; ok {
				c.Set("UserID", userid)
			}
			if role, ok := claims["role"]; ok {
				c.Set("UserRole", role)
			}
		},
	}

	auth := router.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/v2/options/grades", h.GetGrades)
		auth.GET("/v2/options/programs", h.GetPrograms)
		auth.GET("/v2/options/subjects", h.GetSubjects)

		auth.GET("/v2/users/me/subaccounts", h.GetSubaccounts)
		auth.GET("/v2/users/me/subaccounts/:userid", h.GetChildProfile)
		auth.PUT("/v2/users/me/subaccounts/:userid", h.UpdateSubaccount)
		auth.DELETE("/v2/users/me/subaccounts/:userid", h.DeleteSubaccount)
		auth.PUT("/v2/profiles/me", h.UpdateProfile)
		auth.GET("/v2/profiles/:userid", h.GetUserProfile)
	}
	router.GET("/healthcheck", h.Healthcheck)
	router.Run(":80")
}
