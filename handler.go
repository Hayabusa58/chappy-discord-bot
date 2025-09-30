package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
)

var sysprompt = `あなたは Discord サーバでテキストチャットを通して、複数のユーザと同時にやり取りするAI チャットBotです。
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

func readyHandler(b *DiscordBot, oai *OpenAiService, cid string) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Info: Bot logged in as %s", r.User.String())

		// 過去のメッセージ履歴をロード
		historyMessages := b.History.GetMessages(cid)
		if historyMessages == nil {
			// 初回起動時
			log.Println("Info: No history found. Starting initalize...")
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, openai.SystemMessage(sysprompt))
			completion, err := oai.Client.Chat.Completions.New(context.TODO(), b.CompletionParams)
			if err != nil {
				log.Fatalf("Error: An error happend while initalize: %w", err)
				msg := fmt.Sprintf("⚠エラー: Botの初期化処理中にエラーが発生しました。\ndetail:\n```\n%s```", err)
				s.ChannelMessageSend(cid, msg)
				return
			} else {
				log.Println("Info: System prompt response: ", completion.Choices[0].Message.Content)

			}
		} else {
			// 過去の履歴を読み込んで起動
			log.Println("Info: history.json found. Starting load chat history...")
			var apiMessages []openai.ChatCompletionMessageParamUnion
			for _, msg := range historyMessages {
				if msg.Role == "user" {
					apiMessages = append(apiMessages, openai.UserMessage(msg.Content))
				} else {
					apiMessages = append(apiMessages, openai.AssistantMessage(msg.Content))
				}
			}
			b.CompletionParams.Messages.Value = apiMessages
		}

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
			// メッセージが空であれば return
			if m.Content == "" {
				log.Println("Warning: User message has no content.")
				return
			}
			usermassage := m.Author.GlobalName + ": " + m.Content

			// メッセージ履歴に追加
			b.History.AddMessage(cid, "user", usermassage)
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
				msg := fmt.Sprintf("⚠エラー: メッセージの応答処理中にエラーが発生しました。\ndetail:\n```\n%s```", err)
				s.ChannelMessageSend(m.ChannelID, msg)
				return
			}
			// メッセージ履歴に追加
			b.History.AddMessage(cid, "assistant", completion.Choices[0].Message.Content)
			s.ChannelMessageSend(m.ChannelID, completion.Choices[0].Message.Content)
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, completion.Choices[0].Message)

		}
	}

}

func forgetCommandHandler(b *DiscordBot) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "forget" {
			// 保持している会話履歴のリセット
			err := b.History.Forget(i.ChannelID)
			b.CompletionParams.Messages.Value = []openai.ChatCompletionMessageParamUnion{}
			// システムプロンプトだけ入れ直す
			b.CompletionParams.Messages.Value = append(b.CompletionParams.Messages.Value, openai.SystemMessage(sysprompt))
			log.Println("Info: Removing bot history...")
			if err != nil {
				msg := fmt.Sprintf("⚠エラー: 記憶消去処理中にエラーが発生しました。\ndetail:\n```\n%s```", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: msg,
					},
				})
			} else {
				log.Println("Info: Removing bot history...")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "✅チャットボットの記憶を消去しました。",
					},
				})
			}
		}
	}
}
