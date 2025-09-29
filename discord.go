package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
)

type DiscordBot struct {
	Session          *discordgo.Session
	Model            string
	CompletionParams openai.ChatCompletionNewParams
	StopTyping       chan bool
	History          *HistoryManager
}

func NewDiscordBot(tk string, model string, initmsg string, hm *HistoryManager) *DiscordBot {
	session, err := discordgo.New("Bot " + tk)
	parms := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(initmsg),
		}),
		// Geminiのサポートのため一時的にSeedを無効化
		// Seed:  openai.Int(1),
		Model: openai.F(model),
	}
	ch := make(chan bool)
	if err != nil {
		fmt.Println(err)
	}
	return &DiscordBot{
		Session:          session,
		Model:            model,
		CompletionParams: parms,
		StopTyping:       ch,
		History:          hm,
	}
}
