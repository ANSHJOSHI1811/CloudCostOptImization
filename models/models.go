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
    SKUCode         string `gorm:"unique"`
    ProductFamily   string
    VCPU            int
    OperatingSystem string
    InstanceType    string
    Storage         string
    Network         string
    InstanceSKU     string
    Memory          string
	RegionCode		string
    RegionID        uint `gorm:"not null;constraint:OnDelete:CASCADE;"`
	ProviderID 		uint 

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