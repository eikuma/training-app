# node.jsのバージョン指定
ARG NODE_VERSION=22.14.0

FROM node:${NODE_VERSION}-alpine
# 作業ディレクトリ変更
WORKDIR /src/app
# ホストのpackage.jsonをコンテナにコピー
COPY ./package.json ./
# コンテナに必要なパッケージをインストール
RUN npm install
