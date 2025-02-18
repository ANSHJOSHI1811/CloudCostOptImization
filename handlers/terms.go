package handlers

import (
	"net/http"
	"strconv"
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
)
func GetTerms(c *gin.Context) {
	skuIDStr := c.Query("sku_id")

	// Validate SKU ID
	skuID, err := strconv.Atoi(skuIDStr)
	if err != nil || skuID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SKU ID"})
		return
	}

	// Fetch terms by SKU_ID
	var terms []models.Term
	if err := database.DB.Where("sku_id = ?", skuID).Find(&terms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch terms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"SKU_ID": skuID, "terms": terms})
}