package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Store the original jwt.Parse function
var originalJwtParse func(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error)

func mockJwtParse(mockToken *jwt.Token, mockErr error) {
	originalJwtParse = jwt.Parse // Save the original
	jwt.Parse = func(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
		// You could add assertions here on tokenString or keyFunc if needed
		return mockToken, mockErr
	}
}

func restoreJwtParse() {
	if originalJwtParse != nil {
		jwt.Parse = originalJwtParse
	}
}

// Ensure jwtSecret here matches the one in middleware/auth.go
var testJwtSecret = []byte("temporary-secret-key-please-change")


func generateTestToken(userID int64, secret []byte, expiresAt time.Time) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}


func TestAuthMiddleware_Success(t *testing.T) {
	e := echo.New()
	
	// This token is generated for the test, but AuthMiddleware will use the mocked jwt.Parse
	// So the actual content/signature of this token string doesn't strictly matter for THIS test's success
	// as long as jwt.Parse is mocked to return a valid token.
	// However, for a more integrated test, you might generate a real token and NOT mock jwt.Parse
	// but that would require the jwtSecret in the middleware to be the same as used here.
	validTokenString := generateTestToken(123, testJwtSecret, time.Now().Add(time.Hour*1))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+validTokenString)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock jwt.Parse to return a valid token
	mockedClaims := jwt.MapClaims{
		"user_id": float64(123), // JWT decodes numbers as float64
		"exp":     float64(time.Now().Add(time.Hour).Unix()),
	}
	mockedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, mockedClaims)
	// We need to ensure the keyFunc in AuthMiddleware can "validate" this token
	// by providing the correct secret. The mock for jwt.Parse bypasses actual validation.
	
	mockJwtParse(mockedToken, nil) // Mock successful parsing
	defer restoreJwtParse()       // Cleanup

	nextHandlerCalled := false
	handler := AuthMiddleware(func(next_c echo.Context) error {
		nextHandlerCalled = true
		assert.Equal(t, int64(123), next_c.Get("userID"))
		return next_c.String(http.StatusOK, "next handler called")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.True(t, nextHandlerCalled, "Next handler should have been called")
	assert.Equal(t, http.StatusOK, rec.Code) // Assuming next handler returns 200
	assert.Equal(t, int64(123), c.Get("userID"))
}


func TestAuthMiddleware_MissingHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil) // No Authorization header
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := AuthMiddleware(func(next_c echo.Context) error {
		t.Fatal("Next handler should not be called") // Fail test if next is called
		return nil
	})

	err := handler(c)
	assert.NoError(t, err) // Middleware writes response, doesn't return error itself typically
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Authorization header required", resp["error"])
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "InvalidFormat token123") // Not "Bearer "
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := AuthMiddleware(func(next_c echo.Context) error {
		t.Fatal("Next handler should not be called")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Authorization header format must be Bearer {token}", resp["error"])
}

func TestAuthMiddleware_JwtParseError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer sometoken")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockJwtParse(nil, errors.New("mocked jwt parse error")) // Mock jwt.Parse to return an error
	defer restoreJwtParse()

	handler := AuthMiddleware(func(next_c echo.Context) error {
		t.Fatal("Next handler should not be called")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid or expired token", resp["error"]) // Generic error message
}

func TestAuthMiddleware_MissingUserIDClaim(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer validlookingtoken")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock claims without user_id
	mockedClaims := jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())}
	mockedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, mockedClaims)
	
	mockJwtParse(mockedToken, nil)
	defer restoreJwtParse()

	handler := AuthMiddleware(func(next_c echo.Context) error {
		t.Fatal("Next handler should not be called")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid token claims: user_id missing or not a number", resp["error"])
}

func TestAuthMiddleware_UserIDClaimWrongType(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer validlookingtoken")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock claims with user_id as string instead of float64
	mockedClaims := jwt.MapClaims{
		"user_id": "not_a_number", 
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}
	mockedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, mockedClaims)
	
	mockJwtParse(mockedToken, nil)
	defer restoreJwtParse()

	handler := AuthMiddleware(func(next_c echo.Context) error {
		t.Fatal("Next handler should not be called")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid token claims: user_id missing or not a number", resp["error"])
}

// Further tests:
// - Token expired: jwt.Parse would return jwt.ErrTokenExpired. AuthMiddleware currently gives a generic "Invalid or expired token".
//   To test this specifically, the mock for jwt.Parse could return jwt.ErrTokenExpired.
```
The three test files (`user_test.go` for model, `user_test.go` for handler, and `auth_test.go` for middleware) have been created with initial sets of tests.

**Summary of implementation:**

1.  **`backend/app/model/user_test.go`**:
    *   Created `clearUserTable()` helper to ensure test isolation, crucial as these tests run against a real (test) database.
    *   `TestUser_CreateAndGet`: Tests successful user creation (`Create`) and subsequent retrieval using `GetUserByUsername`, `GetUserByEmail`, and `GetUserByID`. Asserts correct data and positive UserID.
    *   `TestUser_Get_NotFound`: Tests that `GetUserByUsername`, `GetUserByEmail`, and `GetUserByID` correctly return `dbr.ErrNotFound` when a user does not exist.
    *   Notes on database dependency: These tests rely on the `training_db` session being available and correctly configured (e.g., via `ENV=test` and `MYSQL_DSN_TEST`). Tests will skip if the session is nil.

2.  **`backend/app/handler/user_test.go`**:
    *   Created `mockUserModel` implementing `model.User` interface to mock model layer interactions.
    *   **`RegisterUser` Tests**:
        *   `TestRegisterUser_Success`: Mocks successful user creation path, checks for 201 Created and success message.
        *   `TestRegisterUser_ValidationErrors`: Table-driven tests for various input validation failures (empty fields, invalid email, short password), checks for 400 Bad Request and specific error messages.
        *   `TestRegisterUser_Conflict_UsernameExists`: Mocks existing username, checks for 409 Conflict.
    *   **`LoginUser` Tests**:
        *   `TestLoginUser_Success`: Mocks successful user lookup and password verification (implicitly, as bcrypt comparison isn't directly mocked here but relies on correct hashed password), checks for 200 OK and presence of a token.
        *   `TestLoginUser_InvalidCredentials_UserNotFound`: Mocks user not found by email, checks for 401 Unauthorized.
        *   `TestLoginUser_InvalidCredentials_WrongPassword`: Mocks user found but with a password that (implicitly) won't match the stored hash, checks for 401 Unauthorized.
    *   Uses `httptest` and `echo.NewContext` for handler testing.
    *   Patches `model.NewUser` to inject the mock model.

3.  **`backend/app/middleware/auth_test.go`**:
    *   Helper functions `mockJwtParse` and `restoreJwtParse` to mock `jwt.Parse`.
    *   `generateTestToken` helper to create valid JWTs for testing (though `jwt.Parse` is mocked for most error cases).
    *   `TestAuthMiddleware_Success`: Mocks successful token parsing, checks that `userID` is set in context and the next handler is called.
    *   `TestAuthMiddleware_MissingHeader`: Checks for 401 if Authorization header is missing.
    *   `TestAuthMiddleware_InvalidHeaderFormat`: Checks for 401 if header format is not "Bearer {token}".
    *   `TestAuthMiddleware_JwtParseError`: Mocks `jwt.Parse` returning an error, checks for 401.
    *   `TestAuthMiddleware_MissingUserIDClaim`: Mocks token claims without `user_id`, checks for 401.
    *   `TestAuthMiddleware_UserIDClaimWrongType`: Mocks token claims with `user_id` of an incorrect type, checks for 401.

This provides a foundational set of unit tests covering happy paths and key error conditions for the authentication features.
