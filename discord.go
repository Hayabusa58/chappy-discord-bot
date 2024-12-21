package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
)

type DiscordBot struct {
	Session          *discordgo.Session
	CompletionParams openai.ChatCompletionNewParams
}

func NewDiscordBot(tk string, initmsg string) *DiscordBot {
	session, err := discordgo.New("Bot " + tk)
	parms := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(initmsg),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT4),
	}
	if err != nil {
		fmt.Println(err)
	}
	return &DiscordBot{
		Session:          session,
		CompletionParams: parms,
	}
}
