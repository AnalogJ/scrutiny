package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/version"
	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"version": version.String(),
	})
}
