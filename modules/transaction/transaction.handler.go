package transaction

import (
	"fmt"
	"medium-rag/config"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService ITransactionService
	router             *gin.Engine
	config             *config.EnvVariable
}

func NewTransactionHandler(transactionService ITransactionService, router *gin.Engine, config *config.EnvVariable) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService, router: router, config: config}
}

func (h *TransactionHandler) getTotalTransactionPerMonth(c *gin.Context) {
	var query GetTotalTransactionReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	data, err := h.transactionService.getTotalTransactionPerMonth(c, query)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"message": "not found", "data": nil})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": data})
}

func (h *TransactionHandler) RegisterRoutes() {
	// In real case, this is authenticated route, and we can get the userId from the token
	h.router.GET("/v1/transactions/total-per-month", h.getTotalTransactionPerMonth)
}
