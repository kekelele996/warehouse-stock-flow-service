package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// ProductCreateRequest 创建商品请求。
type ProductCreateRequest struct {
	OwnerID           uint    `json:"owner_id" binding:"required,gt=0"`
	Name              string  `json:"name" binding:"required,min=1,max=128"`
	SKU               string  `json:"sku" binding:"required,min=1,max=64"`
	Barcode           string  `json:"barcode" binding:"max=64"`
	Category          string  `json:"category" binding:"max=64"`
	Spec              string  `json:"spec" binding:"max=128"`
	Unit              string  `json:"unit" binding:"required,max=16"`
	ShelfLifeDays     *int    `json:"shelf_life_days" binding:"omitempty,min=0"`
	StorageRequirement string `json:"storage_requirement" binding:"required,oneof=Normal ColdChain Dangerous Fragile"`
	Volume            float64 `json:"volume" binding:"min=0"`
	Weight            float64 `json:"weight" binding:"min=0"`
	Price             float64 `json:"price" binding:"min=0"`
}

// ProductUpdateRequest 更新商品请求。
type ProductUpdateRequest struct {
	Name              string  `json:"name" binding:"required,min=1,max=128"`
	Barcode           string  `json:"barcode" binding:"max=64"`
	Category          string  `json:"category" binding:"max=64"`
	Spec              string  `json:"spec" binding:"max=128"`
	Unit              string  `json:"unit" binding:"required,max=16"`
	ShelfLifeDays     *int    `json:"shelf_life_days" binding:"omitempty,min=0"`
	StorageRequirement string `json:"storage_requirement" binding:"required,oneof=Normal ColdChain Dangerous Fragile"`
	Volume            float64 `json:"volume" binding:"min=0"`
	Weight            float64 `json:"weight" binding:"min=0"`
	Price             float64 `json:"price" binding:"min=0"`
}

// ProductQuery 商品查询参数。
type ProductQuery struct {
	OwnerID            uint   `form:"owner_id"`
	Keyword            string `form:"keyword"`
	Category           string `form:"category"`
	StorageRequirement string `form:"storage_requirement"`
}

// ProductView 商品视图。
type ProductView struct {
	ID                 uint      `json:"id"`
	OwnerID            uint      `json:"owner_id"`
	OwnerName          string    `json:"owner_name"`
	Name               string    `json:"name"`
	SKU                string    `json:"sku"`
	Barcode            string    `json:"barcode"`
	Category           string    `json:"category"`
	Spec               string    `json:"spec"`
	Unit               string    `json:"unit"`
	ShelfLifeDays      *int      `json:"shelf_life_days"`
	StorageRequirement string    `json:"storage_requirement"`
	StorageText        string    `json:"storage_text"`
	Volume             float64   `json:"volume"`
	Weight             float64   `json:"weight"`
	Price              float64   `json:"price"`
	CreatedAt          time.Time `json:"created_at"`
}

// ToProductView 模型转视图。
func ToProductView(p *model.Product, ownerName string) *ProductView {
	return &ProductView{
		ID:                 p.ID,
		OwnerID:            p.OwnerID,
		OwnerName:          ownerName,
		Name:               p.Name,
		SKU:                p.SKU,
		Barcode:            p.Barcode,
		Category:           p.Category,
		Spec:               p.Spec,
		Unit:               p.Unit,
		ShelfLifeDays:      p.ShelfLifeDays,
		StorageRequirement: p.StorageRequirement,
		StorageText:        storageText(p.StorageRequirement),
		Volume:             p.Volume,
		Weight:             p.Weight,
		Price:              p.Price,
		CreatedAt:          p.CreatedAt,
	}
}
