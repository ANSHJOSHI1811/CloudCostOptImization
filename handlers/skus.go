package handlers

import (
	"net/http"
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetSKUs handles the request to fetch SKUs with region filter
func GetSKUs(c *gin.Context) {
	region := c.DefaultQuery("region", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.SKU{})

	// Region filter with validation
	if region != "" {
		var regionRecord models.Region
		if err := database.DB.Where("region_id = ?", region).First(&regionRecord).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
			return
		}
		query = query.Where("region_id = ?", regionRecord.RegionID)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count SKUs"})
		return
	}
	totalPages := int64((totalCount + int64(limit) - 1) / int64(limit))
	query = query.Limit(limit).Offset(offset)

	var skus []models.SKU
	if err := query.Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currentPage": page,
		"totalPages":  totalPages,
		"totalCount":  totalCount,
		"data":        skus,
	})
}
