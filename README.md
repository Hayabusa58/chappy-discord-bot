# chappy-discord-bot

## これは何
- OpenAI API を利用したシンプルな ChatGPT のDiscord Botです
- Bot起動中に入力されたメッセージと応答の履歴を保持し、文脈を踏まえて応答します
- **注意: この Bot は複数のサーバに同時にログインする使用方法をサポートしていません**

## 使い方

### 準備
- 事前に適切な権限を与えた Discord Botを作成し、動作させたいサーバに招待しておく
- 以下の情報を取得する
    - Discord bot token 
    - OpenAI API token
    - Bot を動作させたいチャンネルの Chennel ID

### Docker
```
$ mv .env.sample .env

$ vi .env
# OPENAI_TOKEN, DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID を入力する

$ docker build -t chappy-discord-bot .
# コンテナのビルド

$ docker run -d --rm --name chappy-discord-bot chappy-discord-bot
# コンテナの起動
```
動作は保証しませんが、コンテナのリビルドと再起動を行うスクリプトが run.sh です。

