package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/form"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/response"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/service"
	"github.com/labstack/echo/v4"
)

type (
	// Workout ワークアウトのハンドラを表す
	Workout interface {
		List(c echo.Context) error
		Get(c echo.Context) error
		CreateWorkoutSession(c echo.Context) error
		CreateExercise(c echo.Context) error
		CreateSet(c echo.Context) error
	}

	// WorkoutImpl ワークアウトのハンドラを表す
	WorkoutImpl struct {
		WorkoutService service.Workout
	}
)

func NewWorkout() Workout {
	return &WorkoutImpl{
		WorkoutService: service.NewWorkout(),
	}
}

func (h *WorkoutImpl) List(c echo.Context) error {
	f := form.NewListWorkout()
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(400, "invalid form"+err.Error())
	}

	var workoutSessions response.WorkoutSessions
	var err error

	var parsedDate time.Time
	if f.Date != "" {
		var err error
		parsedDate, err = time.Parse(time.RFC3339, f.Date)
		if err != nil {
			return echo.NewHTTPError(400, "invalid date format: "+err.Error())
		}
	}

	// Retrieve requestingUserID from context
	requestingUserID, ok := c.Get("userID").(int64)
	if !ok || requestingUserID == 0 {
		// This should ideally not happen if AuthMiddleware is correctly applied and working.
		// It indicates an issue with how userID is set or retrieved, or middleware not being run.
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context, invalid, or is zero.")
	}

	// Pass requestingUserID, f.ID (as sessionID), and parsedDate to the service layer
	workoutSessions, err = h.WorkoutService.List(requestingUserID, f.ID, parsedDate)
	if err != nil {
		return err
	}

	// It's valid for a user to have no workouts on a given date or for a specific session ID.
	// So, returning an empty list is appropriate here.
	if len(workoutSessions) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{"workouts": []interface{}{}}) // Return 200 OK with empty list
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"workouts": workoutSessions}) // Return 200 OK
}

func (h *WorkoutImpl) Get(c echo.Context) error {
	// Retrieve currentUserID from context
	currentUserID, ok := c.Get("userID").(int64)
	if !ok || currentUserID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found or invalid.")
	}

	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID format.") // Use http.StatusBadRequest
	}

	workoutSessionResponse, err := h.WorkoutService.Get(sessionID)
	if err != nil {
		// The service layer might return specific errors (e.g., not found)
		// which should be propagated or handled appropriately.
		return err // Propagate error from service layer
	}

	// Authorization check: Ensure the workout session belongs to the current user
	if workoutSessionResponse.UserID != currentUserID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"workout": workoutSessionResponse}) // Use http.StatusOK
}

func (h *WorkoutImpl) CreateWorkoutSession(c echo.Context) error {
	f := form.NewCreateWorkoutSession()
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(400, "invalid form"+err.Error())
	}
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return echo.NewHTTPError(400, "validation error "+err.Error())
	}

	parsedDate, err := time.Parse(time.RFC3339, f.Date)

	if err != nil {
		return echo.NewHTTPError(400, "invalid date format: "+err.Error())
	}

	// Retrieve userID from context
	userID, ok := c.Get("userID").(int64)
	log.Printf("CreateWorkoutSession: context_userID=%v, ok=%v\n", userID, ok) // Keep existing log

	if !ok || userID == 0 { // Check if userID is 0
		// If userID is not found, or is 0, treat as unauthorized
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found or invalid, authentication required.")
	}

	// Pass the retrieved userID to the service layer
	workoutSession, err := h.WorkoutService.CreateWorkoutSession(parsedDate, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{"workout": workoutSession})
}

func (h *WorkoutImpl) CreateExercise(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(400, "invalid id")
	}

	f := form.NewCreateExercise()
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(400, "invalid form"+err.Error())
	}
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return echo.NewHTTPError(400, "validation error "+err.Error())
	}

	exercise, err := h.WorkoutService.CreateExercise(id, f.ExerciseName)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{"exercise": exercise})
}

func (h *WorkoutImpl) CreateSet(c echo.Context) error {
	// id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	// if err != nil {
	// 	return echo.NewHTTPError(400, "invalid id")
	// }

	exercise_id, err := strconv.ParseInt(c.Param("exercise_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(400, "invalid exercise_id")
	}

	f := form.NewCreateSet()
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(400, "invalid form"+err.Error())
	}
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return echo.NewHTTPError(400, "validation error "+err.Error())
	}

	sets, err := h.WorkoutService.CreateSet(exercise_id, f.SetNumber, f.Weight, f.Reps)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{"sets": sets})
}
