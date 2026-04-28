package repository

import (
	"context"
	"database/sql"

	"github.com/brota/gobackend/internal/db"
)

type UserRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewUserRepository(queries *db.Queries, db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: queries,
		db:      db,
	}
}

func (r *UserRepository) CreateUserWithID(ctx context.Context, params db.CreateUserParams) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
		}
	}(tx)

	qtx := r.queries.WithTx(tx)

	err = qtx.CreateUser(ctx, params)
	if err != nil {
		return 0, err
	}

	var id int64
	err = tx.QueryRowContext(ctx, "SELECT LAST_INSERT_ID()").Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*db.User, error) {
	user, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) error {
	err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
