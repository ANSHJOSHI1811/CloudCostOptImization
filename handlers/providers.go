package handlers
import (
	"cco_api/database"
	"cco_api/models"
	"net/http"
	"github.com/gin-gonic/gin"
)
type Provider struct {
	ProviderID   uint   `gorm:"primaryKey"`
	ProviderName string `gorm:"unique"`
	CreatedDate  string `gorm:"default:current_timestamp"`
	ModifiedDate string `gorm:"default:current_timestamp"`
	DisableFlag  bool   `gorm:"default:false"`
}

func GetProviders(c *gin.Context) {
	var providers []models.Provider

	// Select only provider_name field
	if err := database.DB.Select("provider_name").Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	// Extract only provider names
	var providerNames []string
	for _, provider := range providers {
		providerNames = append(providerNames, provider.ProviderName)
	}

	c.JSON(http.StatusOK, gin.H{"providers": providerNames})
}