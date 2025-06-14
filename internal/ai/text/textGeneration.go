package text

import (
	"context"

	"google.golang.org/genai"
)

// Text Generation Gemini Model
func Client(ctx context.Context, geminiKey string) (*genai.Client, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GenTextContent(ctx context.Context, client *genai.Client, promt string) (string, error) {

	generatedText, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(promt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return generatedText.Text(), nil

}
