package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// mockUserModel is a mock implementation of the model.User interface.
type mockUserModel struct {
	CreateFunc              func(username, email, passwordHash string) (model.User, error)
	GetUserByUsernameFunc   func(username string) (model.User, error)
	GetUserByEmailFunc      func(email string) (model.User, error)
	GetUserByIDFunc         func(userID int64) (model.User, error)
	// Embed UserImpl to satisfy the interface if methods are called directly on the struct itself (not via interface methods)
	// However, for handler tests, we typically only care about the interface methods being mocked.
	// model.UserImpl 
}

// Implement the model.User interface for mockUserModel
func (m *mockUserModel) Create(username, email, passwordHash string) (model.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(username, email, passwordHash)
	}
	return nil, errors.New("CreateFunc not implemented in mock")
}

func (m *mockUserModel) GetUserByUsername(username string) (model.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return nil, errors.New("GetUserByUsernameFunc not implemented in mock")
}

func (m *mockUserModel) GetUserByEmail(email string) (model.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, errors.New("GetUserByEmailFunc not implemented in mock")
}

func (m *mockUserModel) GetUserByID(userID int64) (model.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(userID)
	}
	return nil, errors.New("GetUserByIDFunc not implemented in mock")
}


// --- RegisterUser Tests ---

func TestRegisterUser_Success(t *testing.T) {
	e := echo.New()
	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Backup and restore original NewUser function
	originalNewUser := model.NewUser
	defer func() { model.NewUser = originalNewUser }()

	mockUser := &model.UserImpl{User: model.User{UserID: 1, Username: "testuser", Email: "test@example.com"}}

	model.NewUser = func() model.User {
		return &mockUserModel{
			GetUserByUsernameFunc: func(username string) (model.User, error) {
				return nil, dbr.ErrNotFound // Username not found
			},
			GetUserByEmailFunc: func(email string) (model.User, error) {
				return nil, dbr.ErrNotFound // Email not found
			},
			CreateFunc: func(username, email, passwordHash string) (model.User, error) {
				return mockUser, nil // Successful creation
			},
		}
	}

	err := RegisterUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User created successfully", resp["message"])
}

func TestRegisterUser_ValidationErrors(t *testing.T) {
	e := echo.New()
	testCases := []struct {
		name          string
		payload       string
		expectedError string
	}{
		{"empty username", `{"username":"","email":"test@example.com","password":"password123"}`, "Username cannot be empty"},
		{"empty email", `{"username":"testuser","email":"","password":"password123"}`, "Email cannot be empty"},
		{"empty password", `{"username":"testuser","email":"test@example.com","password":""}`, "Password cannot be empty"},
		{"invalid email", `{"username":"testuser","email":"invalid","password":"password123"}`, "Invalid email format"},
		{"short password", `{"username":"testuser","email":"test@example.com","password":"pass"}`, "Password must be at least 8 characters long"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := RegisterUser(c)
			assert.NoError(t, err) // Handler itself doesn't error, it writes to response
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			
			var resp map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedError, resp["error"])
		})
	}
}

func TestRegisterUser_Conflict_UsernameExists(t *testing.T) {
	e := echo.New()
	reqBody := `{"username":"existinguser","email":"new@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalNewUser := model.NewUser
	defer func() { model.NewUser = originalNewUser }()

	model.NewUser = func() model.User {
		return &mockUserModel{
			GetUserByUsernameFunc: func(username string) (model.User, error) {
				// Simulate username already exists
				return &model.UserImpl{User: model.User{UserID: 1, Username: "existinguser"}}, nil
			},
		}
	}

	err := RegisterUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Username already exists", resp["error"])
}


// --- LoginUser Tests ---

func TestLoginUser_Success(t *testing.T) {
	e := echo.New()
	email := "login@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	reqBody := `{"email":"` + email + `","password":"` + password + `"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalNewUser := model.NewUser
	defer func() { model.NewUser = originalNewUser }()

	mockUserInstance := &model.UserImpl{User: model.User{UserID: 1, Email: email, PasswordHash: string(hashedPassword)}}

	model.NewUser = func() model.User {
		return &mockUserModel{
			GetUserByEmailFunc: func(e string) (model.User, error) {
				if e == email {
					return mockUserInstance, nil
				}
				return nil, dbr.ErrNotFound
			},
		}
	}

	// Need to ensure jwtSecret is the same as in handler
	// var jwtSecret = []byte("temporary-secret-key-please-change") // from handler/user.go

	err := LoginUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["token"], "Token should not be empty on successful login")
}

func TestLoginUser_InvalidCredentials_UserNotFound(t *testing.T) {
	e := echo.New()
	reqBody := `{"email":"nonexistent@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalNewUser := model.NewUser
	defer func() { model.NewUser = originalNewUser }()

	model.NewUser = func() model.User {
		return &mockUserModel{
			GetUserByEmailFunc: func(email string) (model.User, error) {
				return nil, dbr.ErrNotFound // Simulate user not found
			},
		}
	}

	err := LoginUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid credentials", resp["error"])
}

func TestLoginUser_InvalidCredentials_WrongPassword(t *testing.T) {
	e := echo.New()
	email := "user@example.com"
	correctPassword := "correctPassword"
	wrongPassword := "wrongPassword"
	hashedCorrectPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	reqBody := `{"email":"`+email+`","password":"`+wrongPassword+`"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalNewUser := model.NewUser
	defer func() { model.NewUser = originalNewUser }()
	
	mockUserInstance := &model.UserImpl{User: model.User{UserID: 1, Email: email, PasswordHash: string(hashedCorrectPassword)}}

	model.NewUser = func() model.User {
		return &mockUserModel{
			GetUserByEmailFunc: func(e_addr string) (model.User, error) {
				if e_addr == email {
					return mockUserInstance, nil
				}
				return nil, dbr.ErrNotFound
			},
		}
	}

	err := LoginUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid credentials", resp["error"])
}

// Further tests could include:
// - RegisterUser: Email conflict, model.Create error, bcrypt hashing error (harder to mock)
// - LoginUser: Validation errors (empty email/password), JWT generation error (harder to mock)
```
**Step 3: Create `backend/app/middleware/auth_test.go` (Partial Implementation)**
This will involve mocking `jwt.Parse`.
