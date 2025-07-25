package models

import (
	"time"
)

// PriceProductInfo represents the price_product_info table structure
type PriceProductInfo struct {
	ID                      int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID               string     `json:"product_id" gorm:"column:product_id;type:varchar(50)"`
	CatalogID               string     `json:"catalog_id" gorm:"column:catalog_id;type:varchar(50)"`
	SscatID                 string     `json:"sscat_id" gorm:"column:sscat_id;type:varchar(50)"`
	Category                string     `json:"category" gorm:"column:category;type:varchar(100)"`
	Images                  string     `json:"images" gorm:"column:images;type:text"`
	Weight                  string     `json:"weight" gorm:"column:weight;type:varchar(50)"`
	MallSuperPort           string     `json:"mall_super_port" gorm:"column:mall_super_port;type:varchar(100)"`
	Sscat                   string     `json:"sscat" gorm:"column:sscat;type:varchar(200)"`
	BizFinCategory          string     `json:"biz_fin_category" gorm:"column:biz_fin_category;type:varchar(100)"`
	PidCreated              *time.Time `json:"pid_created" gorm:"column:pid_created;type:datetime"`
	Portfolio               string     `json:"portfolio" gorm:"column:portfolio;type:varchar(200)"`
	CategoryID              string     `json:"category_id" gorm:"column:category_id;type:varchar(50)"`
	ScatID                  string     `json:"scat_id" gorm:"column:scat_id;type:varchar(50)"`
	SuperPortfolio          string     `json:"super_portfolio" gorm:"column:super_portfolio;type:text"`
	PidValid                string     `json:"pid_valid" gorm:"column:pid_valid;type:varchar(10)"`
	Pod                     string     `json:"pod" gorm:"column:pod;type:varchar(100)"`
	Scat                    string     `json:"scat" gorm:"column:scat;type:varchar(100)"`
	Name                    string     `json:"name" gorm:"column:name;type:text"`
	CatalogDt               *time.Time `json:"catalog_dt" gorm:"column:catalog_dt;type:date"`
	CatalogValid            string     `json:"catalog_valid" gorm:"column:catalog_valid;type:varchar(10)"`
	CatalogCreated          *time.Time `json:"catalog_created" gorm:"column:catalog_created;type:datetime"`
	BrandNonBrandTag        string     `json:"brand_non_brand_tag" gorm:"column:brand_non_brand_tag;type:varchar(200)"`
	ParentBrandName         string     `json:"parent_brand_name" gorm:"column:parent_brand_name;type:varchar(200)"`
	MallFlag                string     `json:"mall_flag" gorm:"column:mall_flag;type:varchar(100)"`
	BrandName               string     `json:"brand_name" gorm:"column:brand_name;type:varchar(200)"`
	MallPort                string     `json:"mall_port" gorm:"column:mall_port;type:varchar(100)"`
	Scale                   string     `json:"scale" gorm:"column:scale;type:varchar(100)"`
	BrandTypeTag            string     `json:"brand_type_tag" gorm:"column:brand_type_tag;type:varchar(200)"`
	MelpsFlag               string     `json:"melps_flag" gorm:"column:melps_flag;type:varchar(10)"`
	Cluster                 string     `json:"cluster" gorm:"column:cluster;type:varchar(200)"`
	SupplierID              string     `json:"supplier_id" gorm:"column:supplier_id;type:varchar(50)"`
	SupplierListedPrice     float64    `json:"supplier_listed_price" gorm:"column:supplier_listed_price;type:decimal(10,2)"`
	ShippingRevenue         float64    `json:"shipping_revenue" gorm:"column:shipping_revenue;type:decimal(10,2)"`
	MeeshoPriceWithShipping float64    `json:"meesho_price_with_shipping" gorm:"column:meesho_price_with_shipping;type:decimal(10,2)"`
	CreatedAt               time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for PriceProductInfo
func (PriceProductInfo) TableName() string {
	return "price_product_info"
}
