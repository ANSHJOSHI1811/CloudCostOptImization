package handlers

import (
    "net/http"
    "cco_api/database"
    "cco_api/models"
    "github.com/gin-gonic/gin"
)


func GetSKU(c *gin.Context) {
	var sku models.SKU

	// Get query parameters
	skuID := c.Query("sku_id")
	skuCode := c.Query("skuCode")

	// Fetch SKU based on provided parameter
	if skuID != "" {
		if err := database.DB.Preload("Prices").Where("id = ?", skuID).First(&sku).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
			return
		}
	} else if skuCode != "" {
		if err := database.DB.Preload("Prices").Where("sku_code = ?", skuCode).First(&sku).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide either sku_id or skuCode"})
		return
	}

	// Respond with SKU details
	c.JSON(http.StatusOK, sku)
}