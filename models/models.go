package models

type Region struct {
	RegionID   uint   `gorm:"primaryKey"`
	RegionCode string `gorm:"unique"`
	ProviderID uint   `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
}

type Provider struct {
	ProviderID   uint      `gorm:"primaryKey"`
	ProviderName string    `gorm:"unique"`
	DisableFlag  bool      `gorm:"default:false"`
}

type SavingPlan struct {
	ID                  uint   `gorm:"primaryKey"`
	DiscountedSku       string
	Sku                 string
	LeaseContractLength int
	DiscountedRate      string
	ProviderID 			uint
	RegionCode			string
	DiscountedInstanceType string
	Unit 				string
	RegionID            uint `gorm:"not null;constraint:OnDelete:CASCADE;"`
}

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

type Term struct {
	OfferTermID         int    `gorm:"primaryKey;autoIncrement"`
	SKU_ID              uint   `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
	PriceID             uint   `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
	LeaseContractLength string `gorm:"size:255"`
	PurchaseOption      string `gorm:"size:255"`
	OfferingClass       string `gorm:"size:255"`
}
// Define a struct to merge Term and Price fields
type TermWithPrice struct {
	OfferTermID         int    `json:"offer_term_id"`
	SKU_ID              uint   `json:"sku_id"`
	LeaseContractLength string `json:"lease_contract_length"`
	PurchaseOption      string `json:"purchase_option"`
	OfferingClass       string `json:"offering_class"`
	PricePerUnit        string `json:"price_per_unit"`
	Unit                string `json:"unit"`
	PriceID             uint   `json:"price_id"`
}