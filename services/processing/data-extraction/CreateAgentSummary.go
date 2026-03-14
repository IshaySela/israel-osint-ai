package dataextraction

import (
	"context"
	"encoding/json"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

const prompt = `You are a proffesional text analayzer. Extract the location data from a text and summarize the event.
Produce output with the following format:
{
"enLocations": ["first location, "second location",....],
"heSummary": "short event summary in hebrew. note only data from the event and nothing else."
}`

type AgentSummary struct {
	EnLocations []string `json:"enLocations"`
	HeSummary   string   `json:"heSummary"`
}

func CreateAgentSummary(event models.RawOsintEvent, ctx context.Context, apiKey string, modelName string) (AgentSummary, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(event.Text)},
		Model:        openai.ChatModel(modelName),
	})

	if err != nil {
		return AgentSummary{}, err
	}

	var agentSummary AgentSummary

	err = json.Unmarshal([]byte(resp.OutputText()), &agentSummary)

	if err != nil {
		return AgentSummary{}, err
	}

	return agentSummary, nil
}
