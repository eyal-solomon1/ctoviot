package openai

import (
	"context"
	"fmt"
	openAI "github.com/sashabaranov/go-openai"
)

type OpenAI interface {
	GetAudioTranscription(filePath string, promt string) (string, error)
}

type OpenAIService struct {
	token string
}

// Initialize returns a new OpenAI client configured with specified token
func Initialize(token string) OpenAI {
	return &OpenAIService{
		token: token,
	}
}

// GetAudioTranscription returns the transcription text for the provided audio file and the prompt
func (openai *OpenAIService) GetAudioTranscription(filePath string, promt string) (string, error) {
	var finalPromt string = `Improve the accuracy and timing alignment of the transcriptions,
	splitting the transcriptions into smaller segments as much as possible,
	use this as a description for the audio file` + promt

	client := openAI.NewClient(openai.token)
	resp, err := client.CreateTranscription(
		context.Background(),
		openAI.AudioRequest{
			Model:    openAI.Whisper1,
			FilePath: filePath,
			Prompt:   finalPromt,
			Format:   openAI.AudioResponseFormatSRT,
		},
	)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return "", err
	}
	return resp.Text, nil
}
