package repository_test

import (
	"context"
	"database/sql"
	"github.com/brota/gobackend/internal/user/repository/repository"
	"log"
	"testing"
	"time"

	"github.com/brota/gobackend/internal/shared/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const MySqlImage string = "mysql:8.0"

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		MySqlImage,
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	)
	require.NoError(t, err)

	dsn, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	database, err := sql.Open("mysql", dsn)
	require.NoError(t, err)

	_, err = database.Exec(`
		CREATE TABLE users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			surname VARCHAR(255) NOT NULL,
			age INT,
			country_code VARCHAR(10),
			account_balance VARCHAR(50),
			role VARCHAR(50) NOT NULL,
			is_premium BOOLEAN,
			subscription_tier VARCHAR(50) NOT NULL,
			timezone VARCHAR(50)
		)
	`)
	require.NoError(t, err)

	cleanup := func() {
		err := database.Close()
		if err != nil {
			log.Fatal("Error closing database:" + err.Error())
			return
		}
		err = mysqlContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error terminating container:" + err.Error())
			return
		}
	}

	return database, cleanup
}

func TestUserRepository_CreateUserWithID(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	queries := db.New(database)
	repo := repository.NewUserRepositoryWithQueriesAndConn(queries, database)

	validParams := db.CreateUserParams{
		Name:             "John",
		Surname:          "Doe",
		Role:             db.UsersRoleUser,
		SubscriptionTier: db.UsersSubscriptionTierFree,
	}

	t.Run("Positive - Success", func(t *testing.T) {
		ctx := context.Background()
		id, err := repo.CreateUserWithID(ctx, validParams)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("Negative - Canceled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repo.CreateUserWithID(ctx, validParams)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Negative - Timeout Context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Millisecond)
		_, err := repo.CreateUserWithID(ctx, validParams)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("Negative - Closed DB Connection", func(t *testing.T) {
		err := database.Close()
		if err != nil {
			log.Fatal("Error closing database:" + err.Error())
			return
		}
		ctx := context.Background()
		_, err = repo.CreateUserWithID(ctx, validParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	queries := db.New(database)
	repo := repository.NewUserRepositoryWithQueriesAndConn(queries, database)
	ctx := context.Background()

	id, err := repo.CreateUserWithID(ctx, db.CreateUserParams{
		Name:             "Alice",
		Surname:          "Smith",
		Role:             db.UsersRoleAdmin,
		SubscriptionTier: db.UsersSubscriptionTierPro,
	})
	require.NoError(t, err)

	t.Run("Positive - Success", func(t *testing.T) {
		user, err := repo.GetUserByID(context.Background(), id)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("Negative - Not Found", func(t *testing.T) {
		user, err := repo.GetUserByID(context.Background(), 999999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, user)
	})

	t.Run("Negative - Canceled Context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repo.GetUserByID(canceledCtx, id)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Negative - Closed DB Connection", func(t *testing.T) {
		err := database.Close()
		if err != nil {
			log.Fatal("Error closing database:" + err.Error())
			return
		}
		_, err = repo.GetUserByID(context.Background(), id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	queries := db.New(database)
	repo := repository.NewUserRepositoryWithQueriesAndConn(queries, database)
	ctx := context.Background()

	id, err := repo.CreateUserWithID(ctx, db.CreateUserParams{
		Name:             "Bob",
		Surname:          "Marley",
		Role:             db.UsersRoleUser,
		SubscriptionTier: db.UsersSubscriptionTierBasic,
	})
	require.NoError(t, err)

	updateParams := db.UpdateUserParams{
		ID:      id,
		Name:    "Robert",
		Surname: "Marley",
	}

	t.Run("Positive - Success", func(t *testing.T) {
		err := repo.UpdateUser(context.Background(), updateParams)
		assert.NoError(t, err)

		updatedUser, err := repo.GetUserByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, "Robert", updatedUser.Name)
	})

	t.Run("Negative - Canceled Context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		err := repo.UpdateUser(canceledCtx, updateParams)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Negative - Timeout Context", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Millisecond)
		err := repo.UpdateUser(timeoutCtx, updateParams)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("Negative - Closed DB Connection", func(t *testing.T) {
		err := database.Close()
		if err != nil {
			log.Fatal("Error closing database:" + err.Error())
			return
		}
		err = repo.UpdateUser(context.Background(), updateParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	queries := db.New(database)
	repo := repository.NewUserRepositoryWithQueriesAndConn(queries, database)
	ctx := context.Background()

	id, err := repo.CreateUserWithID(ctx, db.CreateUserParams{
		Name:             "Charlie",
		Surname:          "Chaplin",
		Role:             db.UsersRoleUser,
		SubscriptionTier: db.UsersSubscriptionTierFree,
	})
	require.NoError(t, err)

	t.Run("Positive - Success", func(t *testing.T) {
		err := repo.DeleteUser(context.Background(), id)
		assert.NoError(t, err)

		_, err = repo.GetUserByID(context.Background(), id)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Negative - Canceled Context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		err := repo.DeleteUser(canceledCtx, id)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Negative - Timeout Context", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Millisecond)
		err := repo.DeleteUser(timeoutCtx, id)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("Negative - Closed DB Connection", func(t *testing.T) {
		err := database.Close()
		if err != nil {
			log.Fatal("Error closing database:" + err.Error())
			return
		}
		err = repo.DeleteUser(context.Background(), id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})

}
