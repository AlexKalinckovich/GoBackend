package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brota/gobackend/internal/shared/db"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	GetUserByID(ctx context.Context, id int64) (*db.User, error)
	CreateUserWithID(ctx context.Context, params db.CreateUserParams) (int64, error)
	UpdateUser(ctx context.Context, params db.UpdateUserParams) error
	DeleteUser(ctx context.Context, id int64) error
}

type CachedUserRepository struct {
	repo   UserStore
	client *redis.Client
	ttl    time.Duration
}

func NewCachedUserRepository(repo UserStore, client *redis.Client, ttl time.Duration) *CachedUserRepository {
	return &CachedUserRepository{
		repo:   repo,
		client: client,
		ttl:    ttl,
	}
}

func (c *CachedUserRepository) cacheKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (c *CachedUserRepository) GetUserByID(ctx context.Context, id int64) (*db.User, error) {
	key := c.cacheKey(id)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == nil {
		var user db.User
		if json.Unmarshal(data, &user) == nil {
			return &user, nil
		}

		log.Printf("redis unmarshal error for key %s: %v", key, err)
	} else if !errors.Is(redis.Nil, err) {

		log.Printf("redis get error for key %s: %v", key, err)
	}

	user, err := c.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if jsonData, e := json.Marshal(user); e == nil {
			if setErr := c.client.Set(ctx, key, jsonData, c.ttl).Err(); setErr != nil {
				log.Printf("redis set error for key %s: %v", key, setErr)
			}
		}
	}
	return user, nil
}

func (c *CachedUserRepository) CreateUserWithID(ctx context.Context, params db.CreateUserParams) (int64, error) {
	return c.repo.CreateUserWithID(ctx, params)
}

func (c *CachedUserRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) error {
	if err := c.repo.UpdateUser(ctx, params); err != nil {
		return err
	}

	if err := c.client.Del(ctx, c.cacheKey(params.ID)).Err(); err != nil {
		log.Printf("redis del error for key %s: %v", c.cacheKey(params.ID), err)
	}
	return nil
}

func (c *CachedUserRepository) DeleteUser(ctx context.Context, id int64) error {
	if err := c.repo.DeleteUser(ctx, id); err != nil {
		return err
	}
	if err := c.client.Del(ctx, c.cacheKey(id)).Err(); err != nil {
		log.Printf("redis del error for key %s: %v", c.cacheKey(id), err)
	}
	return nil
}
