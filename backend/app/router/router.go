package router

import (
	"log"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init(e *echo.Echo) {
	err := godotenv.Load("app/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // フロントエンドのオリジン
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	// ワークアウトのハンドラを取得
	workoutHandler := handler.NewWorkout()

	// ワークアウトのルーティングを設定
	e.GET("/workouts", workoutHandler.List)
	e.GET("/workouts/:id", workoutHandler.Get)
	e.POST("/workouts", workoutHandler.CreateWorkoutSession)
	e.POST("/workouts/:id/exercises", workoutHandler.CreateExercise)
	e.POST("/workouts/:id/exercises/:exercise_id/sets", workoutHandler.CreateSet)

	recommendationHandler := handler.NewRecommendation()
	e.POST("/recommendations", recommendationHandler.ProposeTrainingMenu)
}
