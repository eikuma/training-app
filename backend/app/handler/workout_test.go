package handler_test // Use _test package to avoid import cycles if handler imports its own test utilities

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	// "github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/form" // Not directly used if form binding is mocked or simple
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/handler"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/response"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkoutService is a mock type for the service.Workout interface
type MockWorkoutService struct {
	mock.Mock
}

// Implement service.Workout interface
func (m *MockWorkoutService) List(id int64, date time.Time) (response.WorkoutSessions, error) {
	args := m.Called(id, date)
	if args.Get(0) == nil { // Handle nil case for response.WorkoutSessions if error occurs
		return nil, args.Error(1)
	}
	return args.Get(0).(response.WorkoutSessions), args.Error(1)
}

func (m *MockWorkoutService) Get(id int64) (*response.GetWorkoutSession, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.GetWorkoutSession), args.Error(1)
}

func (m *MockWorkoutService) CreateWorkoutSession(date time.Time, userId int64) (*response.WorkoutSession, error) {
	args := m.Called(date, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.WorkoutSession), args.Error(1)
}

func (m *MockWorkoutService) CreateExercise(sessionId int64, exerciseName string) (*response.Exercise, error) {
	args := m.Called(sessionId, exerciseName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.Exercise), args.Error(1)
}

func (m *MockWorkoutService) CreateSet(exerciseID int64, setNumber int64, weight float64, reps int64) (*response.Sets, error) {
	args := m.Called(exerciseID, setNumber, weight, reps)
	if args.Get(0) == nil { // Handle nil case for response.Sets
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.Sets), args.Error(1)
}

func TestCreateWorkoutSession_UserIDZero(t *testing.T) {
	e := echo.New()

	// Valid form data for other checks to pass before userID check
	// The actual date value doesn't matter much here as the handler should error out before service call.
	dateStr := time.Now().Format(time.RFC3339) 
	formData := `{"date":"` + dateStr + `"}` // UserID from form is ignored by handler
	
	req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(formData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.Set("userID", int64(0)) // Simulate userID 0 from AuthMiddleware

	mockService := new(MockWorkoutService)
	// We don't expect CreateWorkoutSession on the service to be called.
	// mockService.On("CreateWorkoutSession", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("should not be called")).Maybe()


	// Use the actual handler implementation, injecting the mock service.
	// The handler.WorkoutImpl struct has WorkoutService as a field.
	// We need to ensure our handler instance uses the mock.
	// The NewWorkout function in handler/workout.go creates WorkoutImpl with a real service.
	// For testing, we directly instantiate WorkoutImpl with our mock.
	workoutHandler := handler.WorkoutImpl{WorkoutService: mockService}

	err := workoutHandler.CreateWorkoutSession(c)

	if assert.Error(t, err, "Expected an error from CreateWorkoutSession") {
		httpError, ok := err.(*echo.HTTPError)
		if assert.True(t, ok, "Error should be an *echo.HTTPError") {
			assert.Equal(t, http.StatusBadRequest, httpError.Code, "HTTP status code should be 400 Bad Request")
			
			// Check the message within the HTTPError
			// echo.NewHTTPError stores the message in the Message field.
			// If it's created with a string, it will be a string.
			errorMessage, isString := httpError.Message.(string)
			if assert.True(t, isString, "HTTPError message should be a string") {
				assert.Equal(t, "Invalid user ID (0) for creating workout session. Please ensure the user account is valid and login is correct.", errorMessage)
			}
		}
	}
	// Assert that CreateWorkoutSession was NOT called on the service
	mockService.AssertNotCalled(t, "CreateWorkoutSession", mock.AnythingOfTypeArgument("time.Time"), int64(0))
	mockService.AssertExpectations(t) // General check for any unexpected calls
}


func TestCreateWorkoutSession_Success(t *testing.T) {
	e := echo.New()
	dateStr := time.Now().Format(time.RFC3339) // "2006-01-02T15:04:05Z07:00"
	formData := `{"date":"` + dateStr + `"}` // UserID from form is ignored

	req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(formData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	expectedUserID := int64(123) // Non-zero UserID
	c.Set("userID", expectedUserID)

	mockService := new(MockWorkoutService)

	// Parse the date string to time.Time for the mock expectation,
	// as the handler will parse it before calling the service.
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	assert.NoError(t, err, "Failed to parse date string for test setup")

	// Expected response from the service
	expectedServiceResponse := &response.WorkoutSession{
		ID:     1,
		Date:   parsedDate.Format("2006-01-02"), // Matches format from response.WorkoutSessionFromModel
		UserID: expectedUserID,
	}

	// Expect a call to CreateWorkoutSession with the parsedDate and expectedUserID
	mockService.On("CreateWorkoutSession", parsedDate, expectedUserID).Return(expectedServiceResponse, nil).Once()

	workoutHandler := handler.WorkoutImpl{WorkoutService: mockService}

	err = workoutHandler.CreateWorkoutSession(c)
	if assert.NoError(t, err, "CreateWorkoutSession should not return an error on success") {
		assert.Equal(t, http.StatusOK, rec.Code, "HTTP status code should be 200 OK")
		
		// Construct the expected JSON response string
		// The key for the workout object in the response is "workout"
		expectedJSON := `{"workout":{"id":1,"date":"` + parsedDate.Format("2006-01-02") + `","user_id":123}}`
		assert.JSONEq(t, expectedJSON, rec.Body.String(), "Response body should match expected JSON")
	}

	mockService.AssertExpectations(t) // Verifies that the expected call to CreateWorkoutSession was made
}

// Add a test for when c.Get("userID") fails the type assertion (returns !ok)
func TestCreateWorkoutSession_UserIDContextError(t *testing.T) {
    e := echo.New()
    dateStr := time.Now().Format(time.RFC3339)
    formData := `{"date":"` + dateStr + `"}`
    req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(formData))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()

    c := e.NewContext(req, rec)
    // Do not set "userID", or set it to a wrong type to simulate context error
    c.Set("userID", "not_an_int64") 

    mockService := new(MockWorkoutService)
    workoutHandler := handler.WorkoutImpl{WorkoutService: mockService}

    err := workoutHandler.CreateWorkoutSession(c)

    if assert.Error(t, err) {
        httpError, ok := err.(*echo.HTTPError)
        if assert.True(t, ok) {
            assert.Equal(t, http.StatusInternalServerError, httpError.Code)
            errorMessage, isString := httpError.Message.(string)
            if assert.True(t, isString) {
                assert.Equal(t, "Error retrieving user ID from context. Ensure you are logged in.", errorMessage)
            }
        }
    }
    mockService.AssertNotCalled(t, "CreateWorkoutSession", mock.AnythingOfTypeArgument("time.Time"), mock.AnythingOfTypeArgument("int64"))
    mockService.AssertExpectations(t)
}
```
