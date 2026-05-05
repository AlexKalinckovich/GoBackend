package repository

import (
	"context"
	"database/sql"

	"github.com/brota/gobackend/internal/db"
)

type BaseRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewBaseRepository(queries *db.Queries, database *sql.DB) *BaseRepository {
	return &BaseRepository{
		queries: queries,
		db:      database,
	}
}

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(base *BaseRepository) *UserRepository {
	return &UserRepository{
		BaseRepository: base,
	}
}

func (r *UserRepository) CreateUserWithID(ctx context.Context, params db.CreateUserParams) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
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

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
