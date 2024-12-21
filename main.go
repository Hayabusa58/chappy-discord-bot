package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	openaitk := os.Getenv("OPENAI_TOKEN")
	openaisv := NewOpenAiService(openaitk)

	// Bot の作成
	discordtk := os.Getenv("DISCORD_BOT_TOKEN")
	// 初期プロンプトを自由に設定できるようにはしている
	initprompt := ""
	bot := NewDiscordBot(discordtk, initprompt)

	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Logged in as %s", r.User.String())
	})

	bot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		cid := os.Getenv("DISCORD_CHANNEL_ID")

		// Bot 自身のメッセージの場合無視する
		if m.Author.ID == s.State.User.ID {
			return
		}
		// メッセージが送られたチャンネルが指定されたものか判定する
		if m.ChannelID == cid {
			// fmt.Println(bot.CompletionParams.Messages.Value)
			bot.CompletionParams.Messages.Value = append(bot.CompletionParams.Messages.Value, openai.UserMessage(m.Content))

			// OpenAI APIへ投げ、返ってきた応答を送信する
			completion, err := openaisv.Client.Chat.Completions.New(context.TODO(), bot.CompletionParams)
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Println(completion.Choices[0].Message.Content)
			s.ChannelMessageSend(m.ChannelID, completion.Choices[0].Message.Content)

			bot.CompletionParams.Messages.Value = append(bot.CompletionParams.Messages.Value, completion.Choices[0].Message)

		}
	})

	if err != nil {
		fmt.Println(err)
	}

	bot.Session.Open()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = bot.Session.Close()
	if err != nil {
		fmt.Println(err)
	}

}