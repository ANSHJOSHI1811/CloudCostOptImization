package routes

import (
	"cco_api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/skus", handlers.FetchSKUs)
}
