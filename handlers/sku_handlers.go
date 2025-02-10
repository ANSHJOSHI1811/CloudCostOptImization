package handlers

import (
	"fmt"
	"cco_api/database"
	"cco_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct for request body
type SKURequest struct {
	Provider       string  `json:"provider"`
	Location       string  `json:"location"`
	CPULimit       int     `json:"cpu_limit"`
	PriceLimit     float64 `json:"price_limit"`
	RAMLimit       int     `json:"ram_limit"`
	BandwidthLimit int     `json:"bandwidth_limit"`
}

// Fetch SKUs based on request parameters
func FetchSKUs(c *gin.Context) {
	var req SKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the provider by name
	var provider models.Provider
	if err := database.DB.Where("provider_name = ?", req.Provider).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	// Find the region by name and provider
	var region models.Region
	if err := database.DB.Where("region_code = ? AND provider_id = ?", req.Location, provider.ProviderID).First(&region).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
		return
	}

	// Query SKUs with filters based on request parameters
	var skus []models.SKU
	query := database.DB.Where("region_id = ?", region.RegionID)

	// Apply CPU Limit filter
	if req.CPULimit > 0 {
		query = query.Where("vcpu <= ?", req.CPULimit)
	}

	// Apply RAM Limit filter
	if req.RAMLimit > 0 {
		query = query.Where("storage LIKE ?", fmt.Sprintf("%%RAM %d%%", req.RAMLimit))
	}

	// Apply Bandwidth Limit filter
	if req.BandwidthLimit > 0 {
		query = query.Where("network LIKE ?", fmt.Sprintf("%%Bandwidth %d%%", req.BandwidthLimit))
	}

	// Execute the query
	if err := query.Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}

	// Return the filtered SKUs
	c.JSON(http.StatusOK, skus)
}
