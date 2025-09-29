package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	openaitk := os.Getenv("OPENAI_TOKEN")

	// Bot の作成
	discordtk := os.Getenv("DISCORD_BOT_TOKEN")
	// 初期プロンプトを自由に設定できるようにはしている
	initprompt := ""
	// チャットモデルの設定
	model := os.Getenv("OPENAI_MODEL")
	// 投稿するチャンネルID
	cid := os.Getenv("DISCORD_CHANNEL_ID")
	// API のエンドポイントURL
	baseUrl := os.Getenv("OPENAI_BASEURL")
	// 設定されていなければ openai のURLを使用する
	if baseUrl == "" {
		baseUrl = "https://api.openai.com/v1/"
	}
	log.Printf("Info: API baseurl is %s", baseUrl)
	log.Println("Info: Bot starting...")
	log.Printf("Info: OpenAI model is %s", model)

	// 記憶の保持期間
	days, _ := strconv.Atoi(os.Getenv("HISTORY_DAYS"))

	if openaitk == "" || model == "" {
		log.Fatal("Error: OpenAI token or OpenAI model not set")
		return
	}

	if discordtk == "" {
		log.Fatal("Error: Discord token not set")
		return
	}
	hm := NewHistoryManager("history.json", days)
	openaisv := NewOpenAiService(openaitk, baseUrl)
	bot := NewDiscordBot(discordtk, model, initprompt, hm)

	bot.Session.AddHandler(readyHandler(bot, openaisv, cid))
	bot.Session.AddHandler(messageCreateHandler(bot, cid, openaisv))
	bot.Session.AddHandler(forgetCommandHandler(bot))
	bot.Session.Open()
	// スラッシュコマンドの登録
	_, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "forget",
		Description: "このチャンネルにおけるBotの記憶を削除します。",
	})
	if err != nil {
		log.Printf("Error: Error happend registering slash command: %v", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	hm.Stop()
	err = bot.Session.Close()

	log.Println("Info: Bot stopping...")
	if err != nil {
		log.Printf("Error: Something went wrong when session closing: %f", err)
	}

}
