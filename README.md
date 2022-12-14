# go_todo_app

このリポジトリは[詳解 Go 言語 Web アプリケーション開発](https://www.amazon.co.jp/%E8%A9%B3%E8%A7%A3Go%E8%A8%80%E8%AA%9EWeb%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E9%96%8B%E7%99%BA-%E6%B8%85%E6%B0%B4%E9%99%BD%E4%B8%80%E9%83%8E-ebook/dp/B0B62K55SL/ref=sr_1_1?adgrpid=138104784699&gclid=CjwKCAjwkaSaBhA4EiwALBgQaB2hrIsngzgcQHoEWL4dUCFG9y4KJY6V8WSwcPX7P51uQRp3KypYFhoCBu8QAvD_BwE&hvadid=626779852727&hvdev=c&hvlocphy=1009318&hvnetw=g&hvqmt=e&hvrand=18193561217761404218&hvtargid=kwd-1686770390525&hydadcr=1798_13549897&jp-ad-ap=0&keywords=%E8%A9%B3%E8%A7%A3go%E8%A8%80%E8%AA%9Eweb%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3%E9%96%8B%E7%99%BA&qid=1665813548&qu=eyJxc2MiOiIwLjgyIiwicXNhIjoiMC44NCIsInFzcCI6IjAuMzcifQ%3D%3D&sr=8-1)に記載のコードを勉強用に模写しながら、カスタマイズしたもの。  
また、書籍のコードの最新版は[こちら](https://github.com/budougumi0617/go_todo_app)にあるので、それも参考にすること。

## ローカル環境構築

### 必要

- Docker
- Go
- openssl

### pem ファイル配置

ローカルで起動する場合は、pem ファイルを配置する必要があります。  
public.pem と secret.pem をを auth/cert に配置して下さい。  
以下に mac の場合の手順を記載します。

```bash
cd [このリポジトリのrootディレクトリ]
mkdir -p auth/cert
openssl genrsa 4096 > auth/cert/secret.pem
openssl rsa -pubout < auth/cert/secret.pem > auth/cert/public.pem
```

※ Mac デフォルトの OpenSSL ではなく、Homebrew で インストールした OpenSSL(openssl@3) で動作確認済み。

### ツール install

```bash
go install github.com/k0kubun/sqldef/cmd/mysqldef@latest
```

### DB 作成

DB スキーマを作成する必要があります。

```bash
make up
make migrate
```

### テスト環境構築(mock 作成)

```bash
make generate
```

## テスト実行

```bash
make up
make test
```

## ローカル環境の動作確認

```bash
# 起動
make up
# Curlでリクエストする例
## ユーザ作成
curl -i -XPOST localhost:18000/register -d '{"name": "admin_user", "password": "test", "role": "admin"}'
## ログイン
curl -i -XPOST localhost:18000/login -d '{"user_name": "admin_user", "password": "test"}'
## トークンセット（ログインのレスポンスに含まれるaccess_tokenを変数にセット）
## 例）
export TODO_TOKEN=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjU0MjY4NzQsImlhdCI6MTY2NTQyNTA3NCwiaXNzIjoiZ2l0aHViLmNvbS91ZWtpR2l0eXV0by9nb190b2RvX2FwcCIsImp0aSI6IjVjZDNjNGNjLTViNzQtNDhkNC1iZDA0LTNhNmE5MWRiOGY0YyIsInJvbGUiOiJhZG1pbiIsInN1YiI6ImFjY2Vzc190b2tlbiIsInVzZXJfbmFtZSI6ImFkbWluX3VzZXIifQ.WoCBJ5MkPYyrxLgitbCB4-LWVTcrt--93MtfyyuoAvaX4npKGEgHgD0XDF-OX9L-WMKPWrumVfCk1Bbpw_IEY8QkFZMAmeiyNUO71UOgwbJKbCWf06ZSWoLVvQra4YLu8NrzUMLrj-0diMg2pyWG62RJPZ8N0iBGE9XdxYKV5_WAPJDe_l9C2HcGLlLqIss-W09VQ5NOmvDwDK_9wVACYGo5YWjLxJaQfK8Zznvzf47jqzH0hH_NYpK_399FDdXV35hedD-bjL8qO_anvplOz-KnZFQEQ_p8CNdOwTlB1wdH7HPDlaOx1kQ5smHAeFyeyRTgiGlEJY08ydZ0gRzypLz0ic41_ch7nEzFarkY4Ub8ny-iUBsf6vbmV5LAZDXBYFfLn1Wu7umyW0hjCWFtouFpBNjWCMAwn8mrbQPECDxfnspcqWcyaamm3f77d8EntcYBnW5RyCJgD7w4FM2RV-juA6gZ4DdiY7Beo_z9Tbx6GrZtGFbehvTXPt4-WOsMw0KrAPoYyrqmExXAeLZvObiE87AOI0-XYNnhap8eZFDXIefpiVPf_6aJ_r6o1CHOMEjINOcRhkiOg3ebzzxLK1SIurT1dEXNOxc7YsGplTNkVin9nDEOvXmNka40z7UGW5PiYC4q2osFy9xCypZJV_q-NKv3kzbJHkOTFJskvj8
## タスク作成
curl -XPOST -H "Authorization: Bearer $TODO_TOKEN" localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden
## タスク一覧取得
curl -XGET -H "Authorization: Bearer $TODO_TOKEN" localhost:18000/tasks
## ログアウト
curl -XGET -H "Authorization: Bearer $TODO_TOKEN" localhost:18000/logout
```

##　その他
lint
```
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

golangci-lint run
```
