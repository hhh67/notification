# 通知CLI

### 1. GCPでOAuth2.0の設定

下記を参考にしてOAuth 2.0 のクライアントID、クライアントシークレット、リダイレクトURLを控えておいてね。

https://developers.google.com/admob/api/v1/getting-started?hl=ja


### 2. 環境変数の設定

下記に従って必要な環境変数を設定してね。

|KEY|VALUE|
|-|-|
|ADMOB_PUBLISHER_ID|ca-app-pub-xxxxxxxxxxxxxxxxのxの部分だけ|
|GOOGLE_OAUTH2_CLIENT_ID|1で控えたやつ|
|GOOGLE_OAUTH2_CLIENT_SECRET|1で控えたやつ|
|GOOGLE_OAUTH2_REDIRECT_URL|1で控えたやつ|
|SLACK_API_TOKEN|SlackのOAuth Tokens|
|SLACK_ADMOB_CHANNEL_ID|送信したいチャンネルID|

### 3. びるど

下記のコマンドを実行し、びるどしてから使ってね。

```
# プロジェクトルートで実行してね。
go build -o bin

# コマンドの一覧が表示されるよ。
./bin

# コマンドのオプションの一覧が表示されるよ
./bin [コマンド] --help
```

### 4. OAuth2.0を含むコマンドをはじめて実行する場合

下記のような文字列が出力されるよ。

```
このURLをブラウザで開いて、認証を完了してね！: https://accounts.google.com/o/oauth2/auth.....
```

認証完了後のリダイレクト先にあるGETパラメータ`code`の値をコピーして、URLデコードしてからターミナルに貼り付けてEnterを押下しようね。

下記のような表示がされたら通知がされるはずだよ！

```
OAuth2認証が完了したよ！
```

## 認証に失敗する場合

`config/token.json`の中身を確認してみよう。

`refresh_token`のフィールドが存在しなかったらトークンの交換ができないのでエラーになってしまうよ。

refresh_tokenは初回の認証時にのみ発行されるので、以下から一度アクセス権を削除してからコマンドを実行し直してみよう。

https://myaccount.google.com/u/0/permissions?pli=1

生成される`config/token.json`にrefresh_tokenが含まれるようになったら成功だよ。