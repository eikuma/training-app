package model

import (
	"testing"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/db"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestMain can be used for setup/teardown, if needed, similar to other tests.
// For now, assume db.Connect() is handled or a test DB is available.

func clearUserTable() {
	sess := db.GetSession("training_db")
	if sess == nil {
		// This might happen if db.Connect() wasn't called or failed.
		// For robust testing, a test setup (e.g., TestMain) should ensure this.
		panic("database session not found for clearing user table")
	}
	// Use a raw query to delete all users to avoid dependencies on model methods
	// and to ensure a clean state.
	_, err := sess.Exec("DELETE FROM users")
	if err != nil {
		panic("failed to clear users table: " + err.Error())
	}
	// Reset auto_increment if necessary (MySQL specific)
	_, err = sess.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		// This might fail if the table doesn't exist or other issues.
		// It's good for keeping IDs predictable in tests but can be optional.
		// log.Printf("Could not reset auto_increment for users table: %v", err)
	}

}

func TestUser_Create(t *testing.T) {
	// Ensure a clean state for this test
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	username := "testuser_create"
	email := "test_create@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	createdUser, err := userModel.Create(username, email, string(hashedPassword))

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	// Type assert to access fields
	createdUserImpl := createdUser
	assert.Equal(t, username, createdUserImpl.Username)
	assert.Equal(t, email, createdUserImpl.Email)
	assert.Equal(t, string(hashedPassword), createdUserImpl.PasswordHash)
	assert.True(t, createdUserImpl.UserID > 0, "UserID should be greater than 0 after creation")

	// Verify by fetching
	fetchedUser, err := userModel.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedUser)
	fetchedUserImpl := fetchedUser
	assert.Equal(t, username, fetchedUserImpl.Username)
}

func TestUser_GetUserByUsername_NotFound(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	_, err := userModel.GetUserByUsername("nonexistentuser")

	assert.Error(t, err)
	assert.Equal(t, dbr.ErrNotFound, err)
}

func TestUser_GetUserByEmail_Existing(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	username := "testgetemail"
	email := "getemail@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err := userModel.Create(username, email, string(hashedPassword))
	assert.NoError(t, err)

	foundUser, err := userModel.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	foundUserImpl := foundUser
	assert.Equal(t, username, foundUserImpl.Username)
	assert.Equal(t, email, foundUserImpl.Email)
}

func TestUser_GetUserByID_Existing(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	username := "testgetid"
	email := "getid@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	createdUser, err := userModel.Create(username, email, string(hashedPassword))
	assert.NoError(t, err)
	createdUserImpl := createdUser
	assert.Equal(t, username, createdUserImpl.Username)
	foundUser, err := userModel.GetUserByID(createdUserImpl.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	foundUserImpl := foundUser
	assert.Equal(t, username, foundUserImpl.Username)
	assert.Equal(t, createdUserImpl.UserID, foundUserImpl.UserID)
}

// Example of how dbr.ErrNotFound should be handled for other Get methods
func TestUser_GetUserByUsername_Existing(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	username := "testgetusername"
	email := "getusername@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err := userModel.Create(username, email, string(hashedPassword))
	assert.NoError(t, err)

	foundUser, err := userModel.GetUserByUsername(username)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	foundUserImpl := foundUser
	assert.Equal(t, username, foundUserImpl.Username)
}

func TestUser_GetUserByEmail_NotFound(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	_, err := userModel.GetUserByEmail("nonexistentemail@example.com")
	assert.Error(t, err)
	assert.Equal(t, dbr.ErrNotFound, err)
}

func TestUser_GetUserByID_NotFound(t *testing.T) {
	clearUserTable()
	sess := db.GetSession("training_db")
	assert.NotNil(t, sess, "Database session should not be nil")

	userModel := NewUser()
	_, err := userModel.GetUserByID(999999) // Assuming this ID won't exist
	assert.Error(t, err)
	assert.Equal(t, dbr.ErrNotFound, err)
}

// It's good practice to also test for duplicate username/email on Create if the DB constraints are active
// However, that requires more complex error checking for duplicate entry errors from dbr/mysql.
// For now, focusing on successful creation and not found cases.
// To test DB errors for Create, one would typically mock the dbr session.
// e.g., mock sess.InsertInto(...).Exec() to return an error.
// This is more involved and might require a different testing approach than current model tests.

// Note: The `db` package and its `Connect` function are assumed to be initialized
// correctly for the "training_db" session, similar to how `exercise_test.go`
// and `workout_session_test.go` operate. If these tests fail due to nil session,
// it indicates that the test DB setup needs to be explicitly managed, possibly
// in a TestMain or by ensuring `db.Connect()` is called with a test DSN.
// The clearUserTable function is a basic way to ensure test isolation.
// For a more robust setup, transactions or a dedicated test database that's reset
// between test runs would be better.
// The current `init()` in `backend/db/db.go` seems to connect to a real DB.
// This means these model tests WILL run against that DB.
// Ensure `MYSQL_DSN_TEST` is set in your environment for these tests to use a test database.
// Or, the `db.Connect()` needs to be adapted for a test environment.
// The `clearUserTable` function is critical if the tests are running against a persistent test DB.

// Placeholder for TestMain if more sophisticated setup/teardown is needed
// func TestMain(m *testing.M) {
// 	// Setup: connect to test DB, run migrations if necessary
// 	// db.Connect("your_test_dsn_here") // Or ensure test DSN is used by db.Init()
// 	// clearUserTable() // Initial clear
//
// 	code := m.Run()
//
// 	// Teardown: clear tables, close connection
// 	// clearUserTable()
// 	// db.Close() // If your db package has a Close function
//
// 	os.Exit(code)
// }

// Mocking dbr.Session for more controlled tests (Example concept - not fully implemented here)
// type mockSessionRunner struct {
// 	dbr.SessionRunner
// 	ExecError error
// 	LoadOneError error
// 	LastInsertId int64
// 	LastInsertIdError error
// }
// func (m *mockSessionRunner) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
// 	if m.ExecError != nil {
// 		return nil, m.ExecError
// 	}
// 	// Return a mock result if needed for LastInsertId
// 	return &mockSQLResult{lastInsertId: m.LastInsertId, lastInsertIdError: m.LastInsertIdError}, nil
// }
// func (m *mockSessionRunner) Select(cols ...string) *dbr.SelectBuilder { /* ... */ }
// // ... other dbr methods to mock ...

// type mockSQLResult struct {
// 	sql.Result
// 	lastInsertId int64
// 	lastInsertIdError error
// }
// func (m *mockSQLResult) LastInsertId() (int64, error) {
// 	return m.lastInsertId, m.lastInsertIdError
// }
// func (m *mockSQLResult) RowsAffected() (int64, error) { return 0, nil }

// And then in test:
// originalGetSession := db.GetSession
// db.GetSession = func(name string) *dbr.Session {
// 	return &dbr.Session{SessionRunner: &mockSessionRunner{...}}
// }
// defer func() { db.GetSession = originalGetSession }()

// This shows how complex dbr mocking can be. For now, the tests rely on the actual test DB connection.

/*

 */
