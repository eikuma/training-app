package main

import (
	"log"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/router"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echoのインスタンスを作成
	e := echo.New()

	// ルーティングの初期化
	router.Init(e)

	// サーバー起動
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
