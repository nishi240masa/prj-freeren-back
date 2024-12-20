# Goの公式イメージを使用
FROM golang:1.23-alpine

# 作業ディレクトリを作成
WORKDIR /go/src/app/
# パッケージをインストール
RUN apk update && apk add git

# アプリケーションコードをコンテナ内にコピー
COPY ./src/ ./


# モジュールをインストール
RUN go mod download

# Airをインストール (公式スクリプトを使用)
RUN go install github.com/air-verse/air@latest
# デフォルトでポート 4000 を使用
ENV PORT=8080

# コンテナがリッスンするポートを指定
EXPOSE 8080


# Air設定ファイルに従ってアプリを起動
CMD ["air", "-c", ".air.toml"]