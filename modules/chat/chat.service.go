package chat

import (
	"context"
	"fmt"
	"medium-rag/config"
	vertexai "medium-rag/utils/vertex-ai"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"resty.dev/v3"
)

type IChatService interface {
	Chat(ctx context.Context, body ChatReq) (string, error)
}

type ChatService struct {
	Env            *config.EnvVariable
	VertexAIClient vertexai.IVertexAIClient
}

func NewChatService(env *config.EnvVariable, vertexAIClient vertexai.IVertexAIClient) IChatService {

	return &ChatService{Env: env, VertexAIClient: vertexAIClient}
}

func (s *ChatService) Chat(ctx context.Context, body ChatReq) (string, error) {
	// In real implementation, we can get the userId from the token in the context
	userId := "123"

	// In real implementation, we can get the user's timezone from the user's profile
	now := time.Now().Format("2006-01-02 15:04:05")

	functionPrompt := fmt.Sprintf("Variables:\n User current time: %s\n\n\n%s", now, body.Message)
	model, functionCalls, err := s.VertexAIClient.GetFinancialFunctions(ctx, functionPrompt)
	if err != nil {
		return "", err
	}

	var funcResp []*genai.FunctionResponse
	for _, fnCall := range functionCalls {
		fmt.Println("The model suggests to call the function %q with args: %v\n", fnCall.Name, fnCall.Args)
		if fnCall.Name == "getTotalTransactionPerMonth" {
			year := fnCall.Args["year"]
			month := fnCall.Args["month"]
			// Convert year and month values to strings for the API call
			yearStr := fmt.Sprintf("%v", year)
			monthStr := fmt.Sprintf("%v", month)

			client := resty.New()
			defer client.Close()

			res, err := client.R().
				EnableTrace().
				SetQueryParam("user_id", userId).
				SetQueryParam("year", yearStr).
				SetQueryParam("month", monthStr).
				Get("http://localhost:8080/v1/transactions/total-per-month")

			if err != nil {
				fmt.Println("fnCall getTotalTransactionPerMonth error", err)
				continue
			}

			if res.StatusCode() == 200 {
				// Update function response with API result
				funcResp = append(funcResp, &genai.FunctionResponse{
					Name: "getTotalTransactionPerMonth",
					Response: map[string]any{
						"content": res.String(),
					},
				})
			}
		}
	}

	// it can be personalized based on user configuration
	styleChat := "Make the response sound chill, fun, and conversational!"
	body.Message = fmt.Sprintf("%s (just keep it %s)", body.Message, styleChat)

	response, err := s.VertexAIClient.Generation(ctx, model, body.Message, funcResp)
	if err != nil {
		return "", err
	}

	return response, nil
}
