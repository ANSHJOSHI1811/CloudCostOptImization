package handlers

import (
	"cco_api/database"
	"cco_api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func GetTerms(c *gin.Context) {
	skuIDStr := c.Query("sku_id")

	// Validate SKU ID
	skuID, err := strconv.Atoi(skuIDStr)
	if err != nil || skuID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SKU ID"})
		return
	}

	// Perform JOIN operation to get Terms with Price details, including PriceID
	var terms []models.TermWithPrice
	err = database.DB.Raw(`
		SELECT 
			t.offer_term_id, t.sku_id, t.price_id, 
			t.lease_contract_length, t.purchase_option, t.offering_class, 
			p.price_per_unit, p.unit
		FROM terms t
		JOIN prices p ON t.price_id = p.price_id
		WHERE t.sku_id = ?`, skuID).Scan(&terms).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch terms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"SKU_ID": skuID, "terms": terms})
}
