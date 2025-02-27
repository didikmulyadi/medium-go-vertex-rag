package utils

import (
	"context"
	"fmt"
	"medium-rag/config"

	"cloud.google.com/go/vertexai/genai"
)

type IVertexAIClient interface {
	Generation(ctx context.Context, model *genai.GenerativeModel, instruction string, funcResponse []*genai.FunctionResponse) (string, error)
	GetFinancialFunctions(ctx context.Context, instruction string) (*genai.GenerativeModel, []genai.FunctionCall, error)
}

type VertexAIClient struct {
	config *config.EnvVariable
	model  *genai.GenerativeModel
}

func NewVertexAIClient(config *config.EnvVariable) (IVertexAIClient, error) {
	client, err := genai.NewClient(
		context.Background(),
		config.GoogleProjectID, config.GoogleLocation,
	)
	if err != nil {
		return &VertexAIClient{}, err
	}

	model := client.GenerativeModel(config.VertexModel)
	model.SetTemperature(0.000000)
	model.SetTopP(0.000000)

	return &VertexAIClient{
		config: config,
		model:  model,
	}, nil
}

func GetResponse(resp *genai.GenerateContentResponse) (string, error) {
	var result string

	if resp != nil && resp.Candidates != nil && len(resp.Candidates) > 0 {
		candidate := resp.Candidates[0]
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			if text, ok := candidate.Content.Parts[0].(genai.Text); ok {
				result = fmt.Sprintf("%s", text)
			}
		}
	}

	if result == "" {
		return result, fmt.Errorf("content response not found")
	}

	return result, nil
}

func (c *VertexAIClient) Generation(ctx context.Context, model *genai.GenerativeModel, instruction string, funcResponse []*genai.FunctionResponse) (string, error) {
	var prompts []genai.Part
	prompts = append(prompts, genai.Text(instruction))

	if len(funcResponse) > 0 {
		for _, fn := range funcResponse {
			prompts = append(prompts, fn)
		}
	}

	res, err := model.GenerateContent(ctx, prompts...)
	if err != nil {
		return "", fmt.Errorf("generation - failed to generate content: %v", err)
	}

	stringResult, err := GetResponse(res)
	if err != nil {
		for _, candidate := range res.Candidates {
			fmt.Printf("Candidate: %v\n", *candidate.Content)
		}
		return "", fmt.Errorf("generation - failed to get content: %v, candidates %v", err, res.Candidates)
	}

	return stringResult, nil
}

func (c *VertexAIClient) GetFinancialFunctions(ctx context.Context, instruction string) (*genai.GenerativeModel, []genai.FunctionCall, error) {
	getTotalTransactionFuncDecl := &genai.FunctionDeclaration{
		Name:        "getTotalTransactionPerMonth",
		Description: "Fetch the total financial transactions for a specific year and month. You can call this function multiple times to get data for different months if a comparison is needed.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"month": {
					Type:        genai.TypeInteger,
					Description: "The month for which to retrieve transactions (1 = January, 12 = December).",
				},
				"year": {
					Type:        genai.TypeInteger,
					Description: "The year for which to retrieve transactions, e.g., 2025.",
				},
			},
			Required: []string{"month", "year"},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"month": {
					Type:        genai.TypeInteger,
					Description: "The requested month.",
				},
				"year": {
					Type:        genai.TypeInteger,
					Description: "The requested year.",
				},
				"total_transaction": {
					Type:        genai.TypeNumber,
					Description: "The total transaction amount for the given year and month (e.g., in rupiah).",
				},
			},
		},
	}

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{getTotalTransactionFuncDecl},
		},
	}

	model := c.model
	model.Tools = tools
	prompt := genai.Text(instruction)

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, []genai.FunctionCall{}, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, []genai.FunctionCall{}, fmt.Errorf("GetFinancialFunctions - no candidates from model, resp: %v", resp)
	} else if len(resp.Candidates[0].FunctionCalls()) == 0 {
		return nil, []genai.FunctionCall{}, fmt.Errorf("GetFinancialFunctions - no function call suggestions from model, resp: %v", resp.Candidates)
	}

	return model, resp.Candidates[0].FunctionCalls(), nil
}
