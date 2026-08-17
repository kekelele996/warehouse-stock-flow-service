package dto

// DashboardSummary 仓库总览数据。
type DashboardSummary struct {
	TodayInboundCount  int64               `json:"today_inbound_count"`
	TodayOutboundCount int64               `json:"today_outbound_count"`
	PendingInbound     int64               `json:"pending_inbound"`
	PendingOutbound    int64               `json:"pending_outbound"`
	BinOccupancyRate   float64             `json:"bin_occupancy_rate"`
	TotalBinCount      int64               `json:"total_bin_count"`
	OccupiedBinCount   int64               `json:"occupied_bin_count"`
	OwnerTop10         []OwnerTopItem      `json:"owner_top10"`
	PendingTasks       []PendingTaskItem   `json:"pending_tasks"`
}

// OwnerTopItem 货主库存金额 TOP10。
type OwnerTopItem struct {
	OwnerID    uint    `json:"owner_id"`
	OwnerName  string  `json:"owner_name"`
	TotalValue float64 `json:"total_value"`
}

// PendingTaskItem 待处理任务。
type PendingTaskItem struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	OrderNo string `json:"order_no"`
	Status  string `json:"status"`
}
