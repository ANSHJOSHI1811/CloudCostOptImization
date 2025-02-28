package handlers

import (
    "net/http"
    "cco_api/database"
    "cco_api/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm" // ✅ Add this import
    "strconv"
)

// GetSKU handles the request to fetch SKUs with filters and prices
func GetSKU(c *gin.Context) {
    region := c.DefaultQuery("region", "")
    vcpuStr := c.DefaultQuery("vcpu", "")
    operatingSystem := c.DefaultQuery("operatingSystem", "")
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "200")

    // Convert page and limit to integers
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
    var skus []models.SKU

    // Start the query with SKU model
    query := database.DB.Model(&models.SKU{}).
        Preload("Prices", func(db *gorm.DB) *gorm.DB {
            return db.Select("price_id, sku_id, effective_date, unit, price_per_unit")
        }).
        Debug()

    // **Filtering Logic**
    if region != "" {
        var regionRecord models.Region
        if err := database.DB.Where("region_code = ?", region).First(&regionRecord).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
            return
        }
        query = query.Where("region_id = ?", regionRecord.RegionID)
    }

    if vcpuStr != "" {
        vcpu, err := strconv.Atoi(vcpuStr)
        if err != nil || vcpu < 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vcpu parameter"})
            return
        }
        query = query.Where("v_cpu = ?", vcpu)
    }

    if operatingSystem != "" {
        query = query.Where("operating_system = ?", operatingSystem)
    }

    // **Count total results for pagination**
    var totalCount int64
    if err := query.Count(&totalCount).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count SKUs"})
        return
    }

    totalPages := int64((totalCount + int64(limit) - 1) / int64(limit))
    query = query.Limit(limit).Offset(offset)

    // **Fetch filtered results with Prices**
    if err := query.Find(&skus).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
        return
    }

    // **Response**
    c.JSON(http.StatusOK, gin.H{
        "currentPage": page,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
        "data":        skus,
    })
}
