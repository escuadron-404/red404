package repositories_test

import (
	"context"
	"fmt"
	"log/slog" // Keep slog import, used for logger in NewUserRepository
	"os"       // For TestMain
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v4" // v4 is correct!

	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/escuadron-404/red404/backend/internal/repositories"
)

// TestMain sets up and tears down test environment for all tests in this package.
func TestMain(m *testing.M) {
	// Initialize default logger for tests
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Run all tests
	code := m.Run()

	// Perform any global tear-down here if needed.
	// For pgxmock, individual mockPool.Close() and mockPool.ExpectationsWereMet()
	// are typically handled within each test function using defer.

	os.Exit(code)
}

func TestUserRepository_Create(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close()

	// `slog.Default()` is used here because TestMain sets the default logger.
	repo := repositories.NewUserRepository(mockPool)

	testUser := &models.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
		// Initialize pointers to non-nil zero values if you intend to send empty strings/nil to DB
		FullName:       new(string),
		Bio:            new(string),
		ProfilePicture: new(string),
		// DeletedAt explicitly nil if it's supposed to be NULL in DB
		DeletedAt: nil, // This is the crucial fix for the first error
	}
	*testUser.FullName = "Test User"
	*testUser.Bio = "A short bio."
	*testUser.ProfilePicture = "http://example.com/pic.jpg"

	mockPool.ExpectQuery(`
		INSERT INTO users \(full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\)
		RETURNING id`).
		WithArgs(
			testUser.FullName,
			testUser.Email,
			testUser.Password,
			testUser.Bio,
			testUser.ProfilePicture,
			pgxmock.AnyArg(), // For CreatedAt
			pgxmock.AnyArg(), // For UpdatedAt
			false,            // For Deleted
			nil,              // For DeletedAt (now correctly matches the explicit nil)
		).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(context.Background(), testUser)

	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if testUser.ID != 1 {
		t.Errorf("expected user ID to be 1, but got: %d", testUser.ID)
	}

	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Create_Error(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close()

	repo := repositories.NewUserRepository(mockPool)

	testUser := &models.User{
		Email:          "error@example.com",
		Password:       "password",
		FullName:       new(string),
		Bio:            new(string),
		ProfilePicture: new(string),
		DeletedAt:      nil, // Ensure this is nil for consistency
	}
	*testUser.FullName = "Error User"
	*testUser.Bio = "Error bio."
	*testUser.ProfilePicture = ""

	// Crucial fix: Expect all 9 arguments that the `Create` method sends.
	// The query pattern can be more general if you prefer.
	mockPool.ExpectQuery(`
		INSERT INTO users \(full_name, email, password, bio, profile_picture, created_at, updated_at, deleted, deleted_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\)
		RETURNING id`).
		WithArgs(
			testUser.FullName,
			testUser.Email,
			testUser.Password,
			testUser.Bio,
			testUser.ProfilePicture,
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
			false,            // Deleted
			nil,              // DeletedAt
		).
		WillReturnError(fmt.Errorf("database connection lost"))

	err = repo.Create(context.Background(), testUser)

	if err == nil {
		t.Error("expected an error, but got none")
	}
	// Check for a specific error message if you wrap it in your repo
	if !strings.Contains(err.Error(), "database connection lost") {
		t.Errorf("expected 'database connection lost' error, but got: %v", err)
	}

	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
