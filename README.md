# 通知CLI

1. GCPでOAuth2.0の設定

下記を参考にしてOAuth 2.0 のクライアント ID とクライアント シークレットをJSON形式でダウンロードしてね。

https://developers.google.com/admob/api/v1/getting-started?hl=ja

ダウンロードしたJSONは`client_secrets.json`という名前にして`config/`に置いてね。

2. 環境変数の設定

下記に従って必要な環境変数を設定してね。

|KEY|VALUE|
|-|-|
|ADMOB_PUBLISHER_ID|ca-app-pub-xxxxxxxxxxxxxxxxのxの部分だけ|

3. びるど

下記のコマンドを実行し、びるどしてから使ってね。

```
# プロジェクトルートで実行してね。
go build -o bin

# コマンドの一覧が表示されるよ。
./bin
```
