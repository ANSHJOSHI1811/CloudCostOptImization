package handlers

import (
	"cco_api/database"
	"cco_api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetSKUs handles fetching SKUs with filters
func GetSKUs(c *gin.Context) {
	region := c.DefaultQuery("region", "")
	vcpuStr := c.DefaultQuery("vcpu", "")
	memoryStr := c.DefaultQuery("memory", "")
	instanceType := c.DefaultQuery("instanceType", "")
	cpuArchitecture := c.DefaultQuery("cpuArchitecture", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "100")

	var vcpu, memory int
	var err error

	if vcpuStr != "" {
		vcpu, err = strconv.Atoi(vcpuStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vCPU parameter"})
			return
		}
	}
	if memoryStr != "" {
		memory, err = strconv.Atoi(memoryStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memory parameter"})
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

	query := database.DB.Model(&models.SKU{})

	// Region filter
	var regionRecord models.Region
	if region != "" {
		if err = database.DB.Where("region_code = ?", region).First(&regionRecord).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
			return
		}
		query = query.Where("region_id = ?", regionRecord.RegionID)
	}

	// Additional Filters
	if vcpuStr != "" {
		query = query.Where("vcpu = ?", vcpu)
	}
	if memoryStr != "" {
		query = query.Where("CAST(memory AS INTEGER) = ?", memory)
	}
	if instanceType != "" {
		query = query.Where("instance_type = ?", instanceType)
	}
	if cpuArchitecture != "" {
		query = query.Where("cpu_architecture = ?", cpuArchitecture)
	}

	// Count total records
	var totalCount int64
	if err = query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count SKUs"})
		return
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	offset := (page - 1) * limit

	// Fetch filtered SKUs
	var skus []models.SKU
	if err = query.Limit(limit).Offset(offset).Find(&skus).Error; err != nil {
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
