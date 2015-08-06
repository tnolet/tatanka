package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVersion(c *gin.Context, version string) {
	c.JSON(http.StatusOK, version)
}
