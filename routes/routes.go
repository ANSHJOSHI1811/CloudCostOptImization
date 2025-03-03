package routes

import (
	"cco_api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/providers", handlers.GetProviders)
	r.GET("/regions", handlers.GetRegions)
	r.GET("/savingplans", handlers.GetSavingPlans)
	r.GET("/skus", handlers.GetSKU)
	r.GET("/prices", handlers.GetPrices)
	r.GET("/terms",handlers.GetTerms)
	r.GET("/priceAndTerms",handlers.GetPriceAndTerms)
	r.GET("/sku",handlers.GetSKU)
}