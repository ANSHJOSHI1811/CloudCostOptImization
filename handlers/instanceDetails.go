package handlers

import (
    "net/http"
    "cco_api/database"
    "cco_api/models"
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
	PhysicalProcessor    string `gorm:"column:physical_processor"`   
	MaxThroughput        string `gorm:"column:max_throughput"`        
	EnhancedNetworking   string `gorm:"column:enhanced_networking"`   
	GPU                  string `gorm:"column:gpu"`                  
	MaxIOPS              string `gorm:"column:max_iops"`       

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
	var sku models.SKU

	// Get query parameters
	skuID := c.Query("sku_id")
	skuCode := c.Query("skuCode")
	vcpu := c.Query("vcpu")
	operatingSystem := c.Query("operatingSystem")
	instanceType := c.Query("instanceType")
	storage := c.Query("storage")
	network := c.Query("network")
	instanceSKU := c.Query("instanceSKU")
	memory := c.Query("memory")
	regionCode := c.Query("regionCode")
	regionID := c.Query("regionID")
	providerID := c.Query("providerID")
	physical_processor:= c.Query("physicalProcessor")
	max_throughput:=c.Query("maxThroughput")
	enhanced_networking:=c.Query("enhancedNetworking")

	// Initialize query with preloads
	query := database.DB.Preload("Prices").
		Preload("Region").
		Preload("Provider")

	// Build dynamic filter conditions
	if skuID != "" {
		query = query.Where("id = ?", skuID)
	}
	if skuCode != "" {
		query = query.Where("sku_code = ?", skuCode)
	}
	if vcpu != "" {
		query = query.Where("vcpu = ?", vcpu)
	}
	if operatingSystem != "" {
		query = query.Where("operating_system = ?", operatingSystem)
	}
	if instanceType != "" {
		query = query.Where("instance_type = ?", instanceType)
	}
	if storage != "" {
		query = query.Where("storage = ?", storage)
	}
	if network != "" {
		query = query.Where("network = ?", network)
	}
	if instanceSKU != "" {
		query = query.Where("instance_sku = ?", instanceSKU)
	}
	if memory != "" {
		query = query.Where("memory = ?", memory)
	}
	if regionCode != "" {
		query = query.Where("region_code = ?", regionCode)
	}
	if regionID != "" {
		query = query.Where("region_id = ?", regionID)
	}
	if providerID != "" {
		query = query.Where("provider_id = ?", providerID)
	}
	if physical_processor != "" {
		query = query.Where("physical_processor = ?", physical_processor)
	}
	if max_throughput != "" {
		query = query.Where("max_throughput = ?", max_throughput)
	}
	if enhanced_networking != "" {
		query = query.Where("enhanced_networking = ?", enhanced_networking)
	}
	// Execute query and fetch SKU
	if err := query.First(&sku).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}

	// Respond with SKU details
	c.JSON(http.StatusOK, sku)
}