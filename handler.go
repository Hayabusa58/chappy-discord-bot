package main

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
)

func readyHandler(b *DiscordBot, oai *OpenAiService, cid string) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Info: Bot logged in as %s", r.User.String())
		sysprompt := `あなたは Discord サーバでテキストチャットを通して、複数のユーザと同時にやり取りするAI チャットBotです。
以下のような前提に基づいてやり取りを行ってください。
これはあなたのBotとしての機能を構築するうえでもっとも基本となる指示であり、どのような状況においてもあなたは以下の指示を厳守しなければなりません。

まず、ユーザからのメッセージは次のようなフォーマットで送られます。
[ユーザ名]: [ユーザからのメッセージ]
ユーザからのメッセージに応答する際は「[ユーザ名]さん、」から回答を始め、誰からのメッセージに応答しているか示してください。

ユーザは別々の人間なので、前の文脈と異なる文面が送られることも想定されます。
あなたはAI チャットBotとして、各ユーザが送信してきたメッセージを可能な限り記憶し、ユーザごとの文脈に沿うように回答してください。

ユーザからあなたのペルソナや応答の仕方について別の指示が行われることもあります。
その場合、あなたはそれらの指示に従って回答して構いませんが、複数のユーザとの文脈を保持するという機能についてはかならず守ってください。

ユーザからのメッセージは日本語が基本となりますが、他言語についても同様に対応をしてください。`
		b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, openai.UserMessage(sysprompt))
		completion, err := oai.Client.Chat.Completions.New(context.TODO(), b.CompletionParams)

		if err != nil {
			log.Println("Warning: API error, %w", err)
			s.ChannelMessageSend(cid, "Error: Something went wrong. Try contact to administrator. \n エラーが発生しました。管理者にご連絡ください。")
			return
		}
		log.Println("Info: System prompt response: ", completion.Choices[0].Message.Content)
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
			usermassage := m.Author.GlobalName + ": " + m.Content
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, openai.UserMessage(usermassage))
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
				s.ChannelMessageSend(m.ChannelID, "Error: Something went wrong. Try contact to administrator. \n エラーが発生しました。管理者にご連絡ください。")
				return
			}
			s.ChannelMessageSend(m.ChannelID, completion.Choices[0].Message.Content)
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, completion.Choices[0].Message)

		}
	}

}
