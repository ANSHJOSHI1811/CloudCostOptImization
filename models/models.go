package models

type Provider struct {
	ProviderID   uint   `gorm:"primaryKey"`
	ProviderName string `gorm:"unique"`
}

type Region struct {
	RegionID   uint   `gorm:"primaryKey"`
	RegionCode string `gorm:"unique"`
	ProviderID uint   `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
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
