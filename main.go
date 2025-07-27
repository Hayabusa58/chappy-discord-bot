package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
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
	// API のエンドポイントURL
	baseUrl := os.Getenv("OPENAI_BASEURL")
	if baseUrl == "" {
		baseUrl = "https://api.openai.com/v1/"
	}
	log.Printf("Info: API baseurl is %s", baseUrl)
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
	openaisv := NewOpenAiService(openaitk, baseUrl)
	bot := NewDiscordBot(discordtk, model, initprompt)

	bot.Session.AddHandler(readyHandler(bot, openaisv, cid))
	bot.Session.AddHandler(messageCreateHandler(bot, cid, openaisv))

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
