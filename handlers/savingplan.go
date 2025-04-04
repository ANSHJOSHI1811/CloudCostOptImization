package handlers

import (
	"net/http"
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetSavingPlans handles the request to fetch SavingPlans with filters
// GetSavingPlans handles the request to fetch SavingPlans with filters
func GetSavingPlans(c *gin.Context) {
	region := c.DefaultQuery("region", "")
	contractLengthStr := c.DefaultQuery("contractLength", "")
	minPriceStr := c.DefaultQuery("minimumPrice", "")
	maxPriceStr := c.DefaultQuery("maximumPrice", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "200")

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minimumPrice parameter"})
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maximumPrice parameter"})
			return
		}
	}

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
	var savingPlans []models.SavingPlan

	query := database.DB.Model(&models.SavingPlan{})

	if region != "" {
		var regionRecord models.Region
		if err := database.DB.Where("region_code = ?", region).First(&regionRecord).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
			return
		}
		query = query.Where("region_id = ?", regionRecord.RegionID)
	}

	if contractLengthStr != "" {
		contractLength, err := strconv.Atoi(contractLengthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contractLength parameter"})
			return
		}
		query = query.Where("lease_contract_length = ?", contractLength)
	}

	if minPriceStr != "" {
		query = query.Where("CAST(discounted_rate AS DECIMAL) >= ?", minPrice)
	}

	if maxPriceStr != "" {
		query = query.Where("CAST(discounted_rate AS DECIMAL) <= ?", maxPrice)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count saving plans"})
		return
	}

	totalPages := int64((totalCount + int64(limit) - 1) / int64(limit))
	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&savingPlans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch saving plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currentPage": page,
		"totalPages":  totalPages,
		"totalCount":  totalCount,
		"data":        savingPlans,
	})
}
