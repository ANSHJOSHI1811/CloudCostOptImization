package handlers

import (
	"net/http"
	"cco_api/database"
	"github.com/gin-gonic/gin"
)

// Define struct to merge Price and Terms data
type PriceAndTerms struct {
	PriceID             uint   `json:"price_id"`
	SKU_ID              uint   `json:"sku_id"`
	EffectiveDate       string `json:"effective_date"`
	Unit                string `json:"unit"`
	PricePerUnit        string `json:"price_per_unit"`
	OfferTermID         *int   `json:"offer_term_id,omitempty"` // Nullable, since some prices may not have a term
	LeaseContractLength *string `json:"lease_contract_length,omitempty"`
	PurchaseOption      *string `json:"purchase_option,omitempty"`
	OfferingClass       *string `json:"offering_class,omitempty"`
}

func GetPriceAndTerms(c *gin.Context) {
	skuID := c.Query("sku_id")
	if skuID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU ID is required"})
		return
	}

	// Use LEFT JOIN to get all prices and terms (if they exist)
	var results []PriceAndTerms
	err := database.DB.Raw(`
		SELECT 
			p.price_id, p.sku_id, p.effective_date, p.unit, p.price_per_unit, 
			t.offer_term_id, t.lease_contract_length, t.purchase_option, t.offering_class
		FROM prices p
		LEFT JOIN terms t ON p.price_id = t.price_id
		WHERE p.sku_id = ?
		ORDER BY p.price_id DESC
	`, skuID).Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price and terms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"SKU_ID": skuID, "data": results})
}
