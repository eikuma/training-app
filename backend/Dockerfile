# goのバージョン指定
ARG GO_VERSION=1.23.3
FROM golang:${GO_VERSION}-alpine AS base

# 作業ディレクトリ変更
WORKDIR /go/src

# main.go実行
CMD ["go", "run", "main.go"]
