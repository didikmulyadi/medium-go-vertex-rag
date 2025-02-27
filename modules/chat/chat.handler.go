package chat

import (
	"fmt"
	"medium-rag/config"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService IChatService
	router      *gin.Engine
	config      *config.EnvVariable
}

func NewChatHandler(chatService IChatService, router *gin.Engine, config *config.EnvVariable) *ChatHandler {
	return &ChatHandler{chatService: chatService, router: router, config: config}
}

func (h *ChatHandler) chatWithAI(c *gin.Context) {
	var body ChatReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	responseMessage, err := h.chatService.Chat(c, body)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message": "we are currently cant process your request, please try again later"})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": responseMessage})
}

func (h *ChatHandler) RegisterRoutes() {
	// In real case, this is authenticated route, and we can get the userId from the token
	h.router.POST("/v1/chats", h.chatWithAI)
}
