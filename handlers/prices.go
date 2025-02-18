package handlers

import (
	"net/http"
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
)
func GetPrices(c *gin.Context) {
	skuID := c.Query("sku_id")
	if skuID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU ID is required"})
		return
	}

	var prices []models.Price
	if err := database.DB.Where("sku_id = ?", skuID).Order("price_id DESC").Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
		return
	}

	c.JSON(http.StatusOK, prices)
}