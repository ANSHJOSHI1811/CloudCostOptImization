package handlers

import (
	"net/http"
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
)

func GetRegions(c *gin.Context) {
	providerName := c.Query("provider") // Get provider from query param
	var regions []models.Region

	if providerName != "" {
		// Fetch regions based on provider name
		if err := database.DB.Joins("JOIN providers ON providers.provider_id = regions.provider_id").
			Where("providers.provider_name = ?", providerName).
			Find(&regions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		// Fetch all regions
		if err := database.DB.Find(&regions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	c.JSON(http.StatusOK, regions)
}
