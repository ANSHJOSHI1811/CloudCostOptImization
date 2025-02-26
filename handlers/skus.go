package handlers
 
import (
    "net/http"
    "cco_api/database"
    "cco_api/models"
    "github.com/gin-gonic/gin"
    "strconv"
)
 
// GetSKUs handles the request to fetch SKUs with filters
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
    query := database.DB.Model(&models.SKU{}).Debug() // Add Debug() to log SQL
 
    // **Fix: Correct Field Names**
    if region != "" {
        var regionRecord models.Region
        if err := database.DB.Where("region_code = ?", region).First(&regionRecord).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
            return
        }
        query = query.Where("region_id = ?", regionRecord.RegionID)
    }
 
    // **Fix: Correct field name "VCPU"**
    if vcpuStr != "" {
        vcpu, err := strconv.Atoi(vcpuStr)
        if err != nil || vcpu < 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vcpu parameter"})
            return
        }
        query = query.Where("v_cpu = ?", vcpu) // Correct field name as per DB
    }
 
    // **Fix: Ensure correct column names**
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
 
    // **Fetch filtered results**
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