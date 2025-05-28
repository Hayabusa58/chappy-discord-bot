package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAiService struct {
	Client  *openai.Client
	BaseUrl string
}

func NewOpenAiService(tk string, url string) *OpenAiService {
	cl := openai.NewClient(option.WithAPIKey(tk), option.WithBaseURL(url))

	return &OpenAiService{
		Client:  cl,
		BaseUrl: url,
	}
}
