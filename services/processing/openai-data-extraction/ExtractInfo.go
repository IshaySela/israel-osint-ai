package openaidataextraction

import (
	"context"
	"os"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	dotenv "github.com/joho/godotenv"
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

var openaiApiKey string = ""

func ExtractInfo(event models.RawOsintEvent, ctx context.Context) (string, error) {
	if openaiApiKey == "" {
		dotenv.Load()
		openaiApiKey = os.Getenv("OPENAI_API_KEY")
	}

	client := openai.NewClient(
		option.WithAPIKey(openaiApiKey),
	)

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Instructions: openai.String(prompt),
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(event.Text)},
		Model:        openai.ChatModelGPT5Mini,
	})

	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}
