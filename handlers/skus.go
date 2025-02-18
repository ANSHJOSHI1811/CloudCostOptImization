package handlers

import (
	"net/http"
	"strconv"
	"cco_api/database"
	"cco_api/models"

	"github.com/gin-gonic/gin"
)

// GetSKUs handles fetching SKUs with filters
func GetSKUs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "100")
	vcpuStr := c.Query("vcpu")                   // Optional: VCPU filter
	osFilter := c.Query("operating_system")      // Optional: Operating System filter
	regionIDStr := c.Query("region_id")          // Optional: Region filter

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

	// Build the query dynamically
	query := database.DB.Model(&models.SKU{})

	// Apply filters based on query parameters
	if vcpuStr != "" {
		vcpu, err := strconv.Atoi(vcpuStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid VCPU value"})
			return
		}
		query = query.Where("vcpu <= ?", vcpu)
	}

	if osFilter != "" {
		query = query.Where("operating_system = ?", osFilter)
	}

	if regionIDStr != "" {
		regionID, err := strconv.Atoi(regionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region_id value"})
			return
		}
		query = query.Where("region_id = ?", regionID)
	}

	// Execute query with pagination
	var skus []models.SKU
	if err := query.Limit(limit).Offset(offset).Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "skus": skus})
}
