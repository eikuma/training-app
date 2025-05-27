package router

import (
	"log"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/handler"
	app_middleware "github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/middleware" // Added middleware import
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	// Apply AuthMiddleware to routes requiring authentication
	e.GET("/workouts", workoutHandler.List, app_middleware.AuthMiddleware)
	e.GET("/workouts/:id", workoutHandler.Get) // Assuming :id route does not need auth for now, or it will be specified
	e.POST("/workouts", workoutHandler.CreateWorkoutSession, app_middleware.AuthMiddleware)
	e.POST("/workouts/:id/exercises", workoutHandler.CreateExercise)              // Assuming these nested routes also need auth if parent does
	e.POST("/workouts/:id/exercises/:exercise_id/sets", workoutHandler.CreateSet) // Or apply middleware individually as needed

	recommendationHandler := handler.NewRecommendation()
	e.POST("/recommendations", recommendationHandler.ProposeTrainingMenu)

	// User authentication routes
	// These handlers are expected to be Echo handlers, not Gin.
	// Ensure handler.RegisterUser and handler.LoginUser are compatible.
	e.POST("/auth/register", handler.RegisterUser)
	e.POST("/auth/login", handler.LoginUser)
}
