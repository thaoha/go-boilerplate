package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
