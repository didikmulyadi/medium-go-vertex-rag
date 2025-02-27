package main

import (
	"fmt"
	"log"
	"medium-rag/config"
	"medium-rag/modules/chat"
	"medium-rag/modules/transaction"
	vertexai "medium-rag/utils/vertex-ai"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	env := config.GetEnv()
	router := gin.Default()

	if env.Mode != "release" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "sa.json")
	}

	// Google Vertex AI Client
	vertexAIClient, err := vertexai.NewVertexAIClient(env)
	if err != nil {
		log.Fatalf("Failed to create Vertex AI client: %v", err)
	}

	// Chat Service
	chatService := chat.NewChatService(env, vertexAIClient)
	chatHandler := chat.NewChatHandler(chatService, router, env)
	chatHandler.RegisterRoutes()

	// Transaction Service
	transactionService := transaction.NewTransactionService(env)
	transactionHandler := transaction.NewTransactionHandler(transactionService, router, env)
	transactionHandler.RegisterRoutes()

	router.Run(fmt.Sprintf(":%s", env.Port))
}
