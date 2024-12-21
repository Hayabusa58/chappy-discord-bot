package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAiService struct {
	Client *openai.Client
}

func NewOpenAiService(tk string) *OpenAiService {
	cl := openai.NewClient(option.WithAPIKey(tk))

	return &OpenAiService{
		Client: cl,
	}
}
