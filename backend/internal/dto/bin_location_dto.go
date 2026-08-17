package dto

import (
	"time"

	"github.com/wmsflow/wmsflow/internal/model"
)

// BinCreateRequest 创建库位请求。
type BinCreateRequest struct {
	Area               string  `json:"area" binding:"required,oneof=A B C D"`
	RackNo             string  `json:"rack_no" binding:"required,min=1,max=16"`
	LayerNo            int     `json:"layer_no" binding:"required,min=1,max=20"`
	ColumnNo           int     `json:"column_no" binding:"required,min=1,max=50"`
	Capacity           float64 `json:"capacity" binding:"required,gt=0"`
	StorageRequirement string  `json:"storage_requirement" binding:"required,oneof=Normal ColdChain Dangerous Fragile"`
}

// BinBatchCreateRequest 批量创建库位请求。
type BinBatchCreateRequest struct {
	Area               string `json:"area" binding:"required,oneof=A B C D"`
	RackNo             string `json:"rack_no" binding:"required,min=1,max=16"`
	LayerStart         int    `json:"layer_start" binding:"required,min=1,max=20"`
	LayerEnd           int    `json:"layer_end" binding:"required,gtefield=LayerStart,max=20"`
	ColumnStart        int    `json:"column_start" binding:"required,min=1,max=50"`
	ColumnEnd          int    `json:"column_end" binding:"required,gtefield=ColumnStart,max=50"`
	Capacity           float64 `json:"capacity" binding:"required,gt=0"`
	StorageRequirement string  `json:"storage_requirement" binding:"required,oneof=Normal ColdChain Dangerous Fragile"`
}

// BinUpdateRequest 更新库位请求。
type BinUpdateRequest struct {
	Capacity float64 `json:"capacity" binding:"required,gt=0"`
	Status   string  `json:"status" binding:"required,oneof=Available Occupied Reserved Maintenance"`
}

// BinQuery 库位查询参数。
type BinQuery struct {
	Area   string `form:"area"`
	Status string `form:"status"`
}

// BinView 库位视图。
type BinView struct {
	ID                uint      `json:"id"`
	Code              string    `json:"code"`
	Area              string    `json:"area"`
	RackNo            string    `json:"rack_no"`
	LayerNo           int       `json:"layer_no"`
	ColumnNo          int       `json:"column_no"`
	Capacity          float64   `json:"capacity"`
	OccupancyRate     float64   `json:"occupancy_rate"`
	StorageRequirement string   `json:"storage_requirement"`
	StorageText       string    `json:"storage_text"`
	Status            string    `json:"status"`
	StatusText        string    `json:"status_text"`
	CreatedAt         time.Time `json:"created_at"`
}

// ToBinView 模型转视图。
func ToBinView(b *model.BinLocation) *BinView {
	return &BinView{
		ID:                 b.ID,
		Code:               binCode(b.Area, b.RackNo, b.LayerNo, b.ColumnNo),
		Area:               b.Area,
		RackNo:             b.RackNo,
		LayerNo:            b.LayerNo,
		ColumnNo:           b.ColumnNo,
		Capacity:           b.Capacity,
		OccupancyRate:      b.OccupancyRate,
		StorageRequirement: b.StorageRequirement,
		StorageText:        storageText(b.StorageRequirement),
		Status:             b.Status,
		StatusText:         binStatusText(b.Status),
		CreatedAt:          b.CreatedAt,
	}
}
// binCode 生成库位编码。
func binCode(area, rack string, layer, column int) string {
	return area + "-" + rack + "-" + itoa(layer) + "-" + itoa(column)
}

// itoa 简易整数转字符串。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// BinContent 库位存放商品信息。
type BinContent struct {
	ProductID  uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	SKU        string `json:"sku"`
	BatchNo    string `json:"batch_no"`
	Quantity   int    `json:"quantity"`
	OwnerName  string `json:"owner_name"`
}

// BinRecommendRequest 推荐库位请求。
type BinRecommendRequest struct {
	ProductID uint `json:"product_id" binding:"required,gt=0"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// BinRecommendResponse 推荐库位响应。
type BinRecommendResponse struct {
	Bin *BinView `json:"bin"`
}
