package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func New(version string) (*gin.Engine, error) {

	log.Println("Initializing api...")

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
