package transaction

type GetTotalTransactionReq struct {
	UserID string `form:"user_id"  binding:"required"`
	Month  int    `form:"month"  binding:"required"`
	Year   int    `form:"year"  binding:"required"`
}

type GetTotalTransactionResp struct {
	TotalTransaction float64 `json:"total_transaction"`
	Month            int     `json:"month"`
	Year             int     `json:"year"`
}
