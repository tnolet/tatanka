package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func New(version string) (*gin.Engine, error) {

	gin.SetMode("release")

	r := gin.New()
	r.Use(gin.Recovery())
	{
		r.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, version)
		})
	}
	return r, nil
}
