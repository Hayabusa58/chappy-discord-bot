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

### サポートされているモデル
.env の `OPENAI_MODEL` に設定可能な値は以下の通り
| 値                                  |
|-------------------------------------------------|
| o1                                              |
| o1-2024-12-17                                   |
| o1-preview                                      |
| o1-preview-2024-09-12                           |
| o1-mini                                         |
| o1-mini-2024-09-12                              |
| gpt-4o                                          |
| gpt-4o-2024-11-20                               |
| gpt-4o-2024-08-06                               |
| gpt-4o-2024-05-13                               |
| gpt-4o-audio-preview                            |
| gpt-4o-audio-preview-2024-10-01                 |
| gpt-4o-audio-preview-2024-12-17                 |
| gpt-4o-mini-audio-preview                       |
| gpt-4o-mini-audio-preview-2024-12-17           |
| chatgpt-4o-latest                               |
| gpt-4o-mini                                     |
| gpt-4o-mini-2024-07-18                          |
| gpt-4-turbo                                     |
| gpt-4-turbo-2024-04-09                          |
| gpt-4-0125-preview                              |
| gpt-4-turbo-preview                             |
| gpt-4-1106-preview                              |
| gpt-4-vision-preview                            |
| gpt-4                                           |
| gpt-4-0314                                      |
| gpt-4-0613                                      |
| gpt-4-32k                                       |
| gpt-4-32k-0314                                  |
| gpt-4-32k-0613                                  |
| gpt-3.5-turbo                                   |
| gpt-3.5-turbo-16k                               |
| gpt-3.5-turbo-0301                              |
| gpt-3.5-turbo-0613                              |
| gpt-3.5-turbo-1106                              |
| gpt-3.5-turbo-0125                              |
| gpt-3.5-turbo-16k-0613                          |



### Docker
```
$ docker build -t chappy-discord-bot .
# コンテナのビルド

$ mv compose.yaml.sample compose.yaml
$ vi compose.yaml
# OPENAI_TOKEN, OPENAI_MODEL, DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID を入力する

$ docker run -d --rm --name chappy-discord-bot chappy-discord-bot
# コンテナの起動
```
デバッグ時には `go run . -debug`と実行することで、直接 .env ファイルを読み出すことができます。
