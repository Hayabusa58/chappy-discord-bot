package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

func main() {
	var (
		debug = flag.Bool("debug", false, "Debug enable")
	)
	flag.Parse()

	if *debug {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error: (debug) Can't load enviroment variables.")
			return
		}
	}

	openaitk := os.Getenv("OPENAI_TOKEN")

	// Bot の作成
	discordtk := os.Getenv("DISCORD_BOT_TOKEN")
	// 初期プロンプトを自由に設定できるようにはしている
	initprompt := ""
	// チャットモデルの設定
	model := os.Getenv("OPENAI_MODEL")
	// 投稿するチャンネルID
	cid := os.Getenv("DISCORD_CHANNEL_ID")
	log.Println("Info: Bot starting...")
	log.Printf("Info: OpenAI model is %s", model)

	if openaitk == "" || model == "" {
		log.Fatal("Error: OpenAI token or OpenAI model not set")
		return
	}

	if discordtk == "" {
		log.Fatal("Error: Discord token not set")
		return
	}
	openaisv := NewOpenAiService(openaitk)
	bot := NewDiscordBot(discordtk, model, initprompt)

	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Info: Bot logged in as %s", r.User.String())
		s.ChannelMessageSend(cid, "Botがログインしました。")
	})

	bot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Bot 自身のメッセージの場合
		if m.Author.ID == s.State.User.ID {
			// 入力中... 表示を停止する
			bot.StopTyping <- true
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
			bot.CompletionParams.Messages.Value = append(bot.CompletionParams.Messages.Value, openai.UserMessage(m.Content))
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
					case <-bot.StopTyping:
						// fmt.Println("Stopping")
						return
					case <-timeout:
						return
					}
				}
			}()
			// OpenAI APIへ投げ、返ってきた応答を送信する
			completion, err := openaisv.Client.Chat.Completions.New(context.TODO(), bot.CompletionParams)

			if err != nil {
				log.Println("Warning: API error, %w", err)
				s.ChannelMessageSend(m.ChannelID, "Error: Something went wrong. Try contact to administrator. \n 何らかのエラーが発生したようです。管理者にご連絡ください。")
				return
			}
			s.ChannelMessageSend(m.ChannelID, completion.Choices[0].Message.Content)
			bot.CompletionParams.Messages.Value = append(bot.CompletionParams.Messages.Value, completion.Choices[0].Message)

		}
	})

	bot.Session.Open()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err := bot.Session.Close()
	log.Println("Info: Bot stopping...")
	if err != nil {
		log.Printf("Error: Something went wrong when session closing: %f", err)
	}

}
