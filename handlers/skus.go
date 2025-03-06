package handlers

import (
    "net/http"
    "cco_api/database"
    "cco_api/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm" // ✅ Add this import
    "strconv"
)



func GetSKUS(c *gin.Context) {
    region := c.DefaultQuery("region", "")
    minVcpuStr := c.DefaultQuery("minVcpu", "")
    maxVcpuStr := c.DefaultQuery("maxVcpu", "")
    operatingSystem := c.DefaultQuery("operatingSystem", "")
    minPriceStr := c.DefaultQuery("minPrice", "")
    maxPriceStr := c.DefaultQuery("maxPrice", "")
    minMemoryStr := c.DefaultQuery("minMemory", "")
    maxMemoryStr := c.DefaultQuery("maxMemory", "")
    minNetworkStr := c.DefaultQuery("minNetwork", "")
    maxNetworkStr := c.DefaultQuery("maxNetwork", "")
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "200")

    // Convert pagination params
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

    // Start Query
    query := database.DB.Model(&models.SKU{}).
        Preload("Prices", func(db *gorm.DB) *gorm.DB {
            return db.Select("price_id, sku_id, effective_date, unit, price_per_unit")
        }).
        Debug()

    // **Filter by Region**
    if region != "" {
        var regionRecord models.Region
        if err := database.DB.Where("region_code = ?", region).First(&regionRecord).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
            return
        }
        query = query.Where("region_id = ?", regionRecord.RegionID)
    }

    // **Filter by vCPU**
    if minVcpuStr != "" || maxVcpuStr != "" {
        var minVcpu, maxVcpu int
        if minVcpuStr != "" {
            minVcpu, err = strconv.Atoi(minVcpuStr)
            if err != nil || minVcpu < 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minVcpu parameter"})
                return
            }
            query = query.Where("v_cpu >= ?", minVcpu)
        }
        if maxVcpuStr != "" {
            maxVcpu, err = strconv.Atoi(maxVcpuStr)
            if err != nil || maxVcpu < minVcpu {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxVcpu parameter"})
                return
            }
            query = query.Where("v_cpu <= ?", maxVcpu)
        }
    }

    // **Filter by Operating System**
    if operatingSystem != "" {
        query = query.Where("operating_system = ?", operatingSystem)
    }

    // **Filter by Network Bandwidth (Min & Max)**
    if minNetworkStr != "" || maxNetworkStr != "" {
        var minNetwork, maxNetwork int
        if minNetworkStr != "" {
            minNetwork, err = strconv.Atoi(minNetworkStr)
            if err != nil || minNetwork < 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minNetwork parameter"})
                return
            }
            query = query.Where("network_bandwidth >= ?", minNetwork)
        }
        if maxNetworkStr != "" {
            maxNetwork, err = strconv.Atoi(maxNetworkStr)
            if err != nil || maxNetwork < minNetwork {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxNetwork parameter"})
                return
            }
            query = query.Where("network_bandwidth <= ?", maxNetwork)
        }
    }

    // **Filter by Memory (Min & Max)**
    if minMemoryStr != "" || maxMemoryStr != "" {
        var minMemory, maxMemory int
        if minMemoryStr != "" {
            minMemory, err = strconv.Atoi(minMemoryStr)
            if err != nil || minMemory < 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minMemory parameter"})
                return
            }
            query = query.Where("memory >= ?", minMemory)
        }
        if maxMemoryStr != "" {
            maxMemory, err = strconv.Atoi(maxMemoryStr)
            if err != nil || maxMemory < minMemory {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxMemory parameter"})
                return
            }
            query = query.Where("memory <= ?", maxMemory)
        }
    }

    // **Filter by Price**
    if minPriceStr != "" || maxPriceStr != "" {
        var minPrice, maxPrice float64
        if minPriceStr != "" {
            minPrice, err = strconv.ParseFloat(minPriceStr, 64)
            if err != nil || minPrice < 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minPrice parameter"})
                return
            }
        }
        if maxPriceStr != "" {
            maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
            if err != nil || maxPrice < minPrice {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxPrice parameter"})
                return
            }
        }

        query = query.Joins("JOIN prices ON skus.id = prices.sku_id")

        if minPriceStr != "" {
            query = query.Where("CAST(prices.price_per_unit AS FLOAT) >= ?", minPrice)
        }
        if maxPriceStr != "" {
            query = query.Where("CAST(prices.price_per_unit AS FLOAT) <= ?", maxPrice)
        }
    }

    // **Count total results for pagination**
    var totalCount int64
    if err := query.Count(&totalCount).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count SKUs"})
        return
    }
    totalPages := int64((totalCount + int64(limit) - 1) / int64(limit))
    query = query.Limit(limit).Offset(offset)

    // **Fetch results**
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




