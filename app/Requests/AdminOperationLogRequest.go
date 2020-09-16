package Requests

//日志列表
type AdminOperationLogIndexRequest struct {
	Localtime uint64                             `json:"localtime" binding:"required"`
	Data      *AdminOperationLogIndexRequestData `json:"data" binding:"required"`
}
type AdminOperationLogIndexRequestData struct {
	Page      uint64 `json:"page" binding:"required,gte=1"`
	Limit     uint8  `json:"limit" binding:"required,gte=1"`
	Keywords  string `json:"keywords" binding:"min:2"`
	StartTime string `json:"start_time" binding:"required" time_format:"2020-09-01 10:22:22"`
	EndTime   string `json:"end_time" binding:"required" time_format:"2020-09-01 10:22:22"`
}
