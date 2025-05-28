package main

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
)

func readyHandler(cid string) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Info: Bot logged in as %s", r.User.String())
		s.ChannelMessageSend(cid, "Botがログインしました。")
	}
}

func messageCreateHandler(b *DiscordBot, cid string, oai *OpenAiService) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Bot 自身のメッセージの場合
		if m.Author.ID == s.State.User.ID {
			// 入力中... 表示を停止する
			b.StopTyping <- true
			return
		}
		// メッセージが送られたチャンネルが指定されたものか判定する
		if m.ChannelID == cid {
			// fmt.Println(bot.CompletionParams.Messages.Value)
			// メッセージが空であれば return
			if m.Content == "" {
				log.Println("Warning: User message has no content.")
				return
			}
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, openai.UserMessage(m.Content))
			// 入力中... 表示を開始するゴルーチン
			go func() {
				s.ChannelTyping(cid)
				t := time.NewTicker(10 * time.Second)
				defer t.Stop()
				timeout := time.After(1 * time.Minute)
				for {
					select {
					case <-t.C:
						// fmt.Println("typing called")
						s.ChannelTyping(cid)
					case <-b.StopTyping:
						// fmt.Println("Stopping")
						return
					case <-timeout:
						return
					}
				}
			}()
			// OpenAI APIへ投げ、返ってきた応答を送信する
			completion, err := oai.Client.Chat.Completions.New(context.TODO(), b.CompletionParams)

			if err != nil {
				log.Println("Warning: API error, %w", err)
				s.ChannelMessageSend(m.ChannelID, "Error: Something went wrong. Try contact to administrator. \n 何らかのエラーが発生したようです。管理者にご連絡ください。")
				return
			}
			s.ChannelMessageSend(m.ChannelID, completion.Choices[0].Message.Content)
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, completion.Choices[0].Message)

		}
	}

}
