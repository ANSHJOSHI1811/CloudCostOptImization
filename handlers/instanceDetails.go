package handlers

import (
    "net/http"
    "cco_api/database"
    "github.com/gin-gonic/gin"
)
type SKU struct {
	ID              uint   `gorm:"primaryKey"`
	RegionID        uint   `gorm:"not null;constraint:OnDelete:CASCADE;"` 
	ProviderID      uint   `gorm:"not null"`
	RegionCode      string `gorm:"not null"`
	SKUCode         string `gorm:"unique"`
	InstanceSKU     string
	ProductFamily   string
	VCPU            int
	CpuArchitecture string
	InstanceType    string
	Storage         string
	Network         string
	OperatingSystem string
	Type 		    string 
	Memory          string
	PhysicalProcessor    string  
	MaxThroughput        string     
	EnhancedNetworking   string   
	GPU                  string                   
	MaxIOPS              string 

    // ✅ Add this line to establish the relation
    Prices []Price `gorm:"foreignKey:SKU_ID"`  
}
type Price struct {
    PriceID       uint   `gorm:"primaryKey;autoIncrement"`
    SKU_ID        uint   `gorm:"not null;constraint:OnDelete:CASCADE;"`
    EffectiveDate string `gorm:"type:varchar(255)"`
    Unit          string `gorm:"type:varchar(50)"`
    PricePerUnit  string `gorm:"type:varchar(50)"`
}

func GetDetails(c *gin.Context) {
	var sku SKU

	// Get query parameters
	skuID := c.Query("sku_id")
	skuCode := c.Query("skuCode")

	// Fetch SKU based on provided parameter
	if skuID != "" {
		if err := database.DB.Preload("Prices").Where("id = ?", skuID).First(&sku).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
			return
		}
	} else if skuCode != "" {
		if err := database.DB.Preload("Prices").Where("sku_code = ?", skuCode).First(&sku).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide either sku_id or skuCode"})
		return
	}

	// Respond with SKU details
	c.JSON(http.StatusOK, sku)
}