package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tnolet/tatanka/control"
	"log"
	"net/http"
)

func New(version string, ctrl *control.Controller) (*gin.Engine, error) {

	log.Println("Initializing api...")

	gin.SetMode("release")

	r := gin.New()
	r.Use(gin.Recovery())
	{
		r.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, version)
		})

		r.GET("/state", func(c *gin.Context) {
			c.JSON(http.StatusOK, ctrl.State())
		})
	}
	return r, nil
}
