package handler

import (
	"log"
	"net/http"
	"net/mail"
	"regexp" // Using regexp for simple email validation as per instructions
	"time"   // Added for JWT expiration

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model"
	"github.com/gocraft/dbr/v2"
	"github.com/golang-jwt/jwt/v5" // Added for JWT generation
	"github.com/labstack/echo/v4"  // Replaced Gin with Echo
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret is used to sign and validate JWT tokens.
// TODO: Move this to a configuration file or environment variable.
var jwtSecret = []byte("temporary-secret-key-please-change")

// RegisterUserRequest defines the structure for the user registration request body.
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// isValidEmail checks if the email format is valid using a simple regex.
// Using net/mail.ParseAddress is generally better for comprehensive validation.
func isValidEmail(email string) bool {
	// A very basic regex for email validation: something@something.something
	// For more robust validation, consider net/mail.ParseAddress or a more complex regex.
	pattern := `.+@.+\..+`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// RegisterUser handles the user registration process.
func RegisterUser(c echo.Context) error { // Changed signature to echo.Context and returns error
	var request RegisterUserRequest

	// Bind JSON request body to the struct
	if err := c.Bind(&request); err != nil { // Changed to c.Bind for Echo
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload: " + err.Error()}) // Changed to echo.Map
	}

	// Validate input
	if request.Username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username cannot be empty"})
	}
	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email cannot be empty"})
	}
	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Password cannot be empty"})
	}

	// Validate email format using net/mail.ParseAddress
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid email format"})
	}
	// Fallback or alternative simple regex validation (as per instruction variations)
	// if !isValidEmail(request.Email) {
	// return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid email format"})
	// }

	// Validate password length
	if len(request.Password) < 8 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Password must be at least 8 characters long"})
	}

	userModel := model.NewUser()

	// Check if username already exists
	_, err := userModel.GetUserByUsername(request.Username)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Username already exists"})
	} else if err != dbr.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error checking username: " + err.Error()})
	}

	// Check if email already exists
	_, err = userModel.GetUserByEmail(request.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Email already exists"})
	} else if err != dbr.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error checking email: " + err.Error()})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password: " + err.Error()})
	}

	// Create the user
	_, err = userModel.Create(request.Username, request.Email, string(hashedPassword))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User created successfully"})
}

// LoginUserRequest defines the structure for the user login request body.
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUser handles user login and JWT token generation.
func LoginUser(c echo.Context) error { // Changed signature to echo.Context and returns error
	var request LoginUserRequest

	// Bind JSON request body to the struct
	if err := c.Bind(&request); err != nil { // Changed to c.Bind for Echo
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload: " + err.Error()}) // Changed to echo.Map
	}

	// Validate input
	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email cannot be empty"})
	}
	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Password cannot be empty"})
	}

	userModel := model.NewUser()

	// Fetch user by email
	dbUserInterface, err := userModel.GetUserByEmail(request.Email)
	if err != nil {
		if err == dbr.ErrNotFound {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error: " + err.Error()})
	}

	// Type assertion to access User struct fields from the interface
	dbUser := dbUserInterface

	log.Printf("LoginUser: Retrieved dbUser.UserID = %d, dbUser = %+v\n", dbUser.UserID, dbUser)

	if dbUser.UserID == 0 {
		log.Printf("LoginUser: Attempt to login with UserID 0 for email %s. Aborting token generation.", request.Email)
		// Return an error similar to other validation errors or invalid credentials.
		// Using StatusUnauthorized to avoid revealing specific account issues.
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials or user account issue."})
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(request.Password))
	if err != nil {
		// If err is bcrypt.ErrMismatchedHashAndPassword, then it's an invalid password.
		// Otherwise, it could be some other error during comparison.
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),                     // Issued at current time
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}
