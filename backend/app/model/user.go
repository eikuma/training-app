package model

import (
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/db"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// User represents the structure of a user in the database.
type User struct {
	UserID       int64  `db:"user_id,auto_increment"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

// UserImpl embeds User and implements the User interface.
type UserImpl struct {
	User
}

// User defines the interface for user operations.
type User interface {
	Create(username, email, passwordHash string) (User, error)
	GetUserByUsername(username string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByID(userID int64) (User, error)
}

// NewUser creates a new User object.
func NewUser() User {
	return &UserImpl{}
}

// Create inserts a new user record into the database.
func (u *UserImpl) Create(username, email, passwordHash string) (User, error) {
	sess := db.GetSession("training_db")
	if sess == nil {
		return nil, errors.New("database session not found")
	}

	user := UserImpl{
		User: User{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	}

	result, err := sess.InsertInto("users").
		Columns("username", "email", "password_hash").
		Record(&user.User).
		Exec()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert user: %s", username)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get last insert ID for user: %s", username)
	}
	user.UserID = userID

	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (u *UserImpl) GetUserByUsername(username string) (User, error) {
	sess := db.GetSession("training_db")
	if sess == nil {
		return nil, errors.New("database session not found")
	}

	var user UserImpl
	err := sess.Select("*").
		From("users").
		Where("username = ?", username).
		LoadOne(&user.User)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, dbr.ErrNotFound
		}
		return nil, errors.Wrapf(err, "failed to get user by username: %s", username)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
func (u *UserImpl) GetUserByEmail(email string) (User, error) {
	sess := db.GetSession("training_db")
	if sess == nil {
		return nil, errors.New("database session not found")
	}

	var user UserImpl
	err := sess.Select("*").
		From("users").
		Where("email = ?", email).
		LoadOne(&user.User)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, dbr.ErrNotFound
		}
		return nil, errors.Wrapf(err, "failed to get user by email: %s", email)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func (u *UserImpl) GetUserByID(userID int64) (User, error) {
	sess := db.GetSession("training_db")
	if sess == nil {
		return nil, errors.New("database session not found")
	}

	var user UserImpl
	err := sess.Select("*").
		From("users").
		Where("user_id = ?", userID).
		LoadOne(&user.User)

	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, dbr.ErrNotFound
		}
		return nil, errors.Wrapf(err, "failed to get user by ID: %d", userID)
	}
	return &user, nil
}
