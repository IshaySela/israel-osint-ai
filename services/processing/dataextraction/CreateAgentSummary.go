package dataextraction

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

const prompt = `You are a professional text analyzer.
Summarize the event described in the user text and extract the locations data from the text.
}`

type AgentSummary struct {
	Locations []string `json:"enLocations" jsonschema_description:"List of city/towns/areas extracted from the text."`
	HeSummary string   `json:"heSummary" jsonschema_description:"Short event summary in hebrew if the event is in hebrew, otherwise in english."`
}

/*Create a JSON schema for the agent summary response*/
/*The schema is passed to the OpenAI API in order to ensure the response is structured correctly*/
func generateAgentSummaryResponseSchema() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v AgentSummary
	schema := reflector.Reflect(v)

	b, _ := json.Marshal(schema)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m

}

func CreateAgentSummary(rawText string, ctx context.Context, apiKey string, modelName string) (AgentSummary, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(rawText)},
		Model:        openai.ChatModel(modelName),
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Schema:      generateAgentSummaryResponseSchema(),
					Name:        "summarize_and_extract_locations",
					Description: openai.String("Summarize the event and extract the location data from the text."),
					Strict:      openai.Bool(true),
				},
			},
		}})

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
