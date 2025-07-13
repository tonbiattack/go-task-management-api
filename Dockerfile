# ビルドステージ
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Go モジュールファイルをコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 実行ステージ
FROM alpine:latest

# CA証明書とタイムゾーンデータを追加
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# ポート8080を公開
EXPOSE 8080

# アプリケーション実行
CMD ["./main"]