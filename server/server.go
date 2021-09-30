package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ramseskamanda/vat-service/config"
	"github.com/ramseskamanda/vat-service/controllers"
)

func Start() error {
	config.Init()
	router := NewRouter()
	return router.Run(":" + config.GetString("PORT"))
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", controllers.HealthCheck)

	v1 := router.Group("v1")
	{
		vat := new(controllers.VATController)
		v1.POST("/vat", vat.CheckVATNumber)
	}
	return router
}
