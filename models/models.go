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
	Processor       string
	UsageType       string
	RegionID        uint `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
}