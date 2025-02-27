package chat

type ChatReq struct {
	Message string `form:"message" json:"message" binding:"required"`
}
