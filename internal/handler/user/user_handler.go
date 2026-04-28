package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/brota/gobackend/internal/handler"
	"io"
	"net/http"
	"strconv"

	"github.com/brota/gobackend/internal/custom_errors/validation"
	"github.com/brota/gobackend/internal/db"
	"github.com/brota/gobackend/internal/repository"
	"github.com/brota/gobackend/internal/transport"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	rh        *handler.RequestHandler
	repo      *repository.UserRepository
	validator *Validator
}

func NewUserHandler(repo *repository.UserRepository) *Handler {
	registry := transport.NewErrorRegistry()
	registry.Register(validation.AggregateErrorCode, func(err error, ctx map[string]any) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusBadRequest, err)
	})

	registry.Register(validation.ErrorCode, func(err error, ctx map[string]any) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusBadRequest, map[string]any{
			"error": ctx["message"],
			"field": ctx["field"],
		})
	})
	rh := handler.NewRequestHandler(registry)
	return &Handler{
		rh:        rh,
		repo:      repo,
		validator: NewUserValidator(),
	}
}

type Response struct {
	ID               int64                    `json:"id"`
	Name             string                   `json:"name"`
	Surname          string                   `json:"surname"`
	Role             db.UsersRole             `json:"role"`
	SubscriptionTier db.UsersSubscriptionTier `json:"subscription_tier"`
	Age              int                      `json:"age,omitempty"`
	CountryCode      string                   `json:"country_code,omitempty"`
	Timezone         string                   `json:"timezone,omitempty"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Age         int    `json:"age,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

type PatchUserRequest struct {
	Name        *string `json:"name,omitempty"`
	Surname     *string `json:"surname,omitempty"`
	Age         *int    `json:"age,omitempty"`
	CountryCode *string `json:"country_code,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("body", "invalid request body")
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		errCh := make(chan *validation.AggregateError)
		go h.validator.ValidateCreateOrUpdate(req, errCh)
		if err := <-errCh; err != nil {
			return transport.HTTPResponse{}, err
		}

		params := db.CreateUserParams{
			Name:             req.Name,
			Surname:          req.Surname,
			Role:             db.UsersRoleUser,
			IsPremium:        sql.NullBool{Bool: false, Valid: true},
			SubscriptionTier: db.UsersSubscriptionTierFree,
			AccountBalance:   sql.NullString{String: "0", Valid: true},
		}
		if req.Age > 0 {
			params.Age = sql.NullInt32{Int32: int32(req.Age), Valid: true}
		}
		if req.CountryCode != "" {
			params.CountryCode = sql.NullString{String: req.CountryCode, Valid: true}
		}
		if req.Timezone != "" {
			params.Timezone = sql.NullString{String: req.Timezone, Valid: true}
		}

		userID, err := h.repo.CreateUserWithID(r.Context(), params)
		if err != nil {
			return transport.HTTPResponse{}, err
		}

		resp := Response{
			ID:               userID,
			Name:             req.Name,
			Surname:          req.Surname,
			Role:             db.UsersRoleUser,
			SubscriptionTier: db.UsersSubscriptionTierFree,
			Age:              req.Age,
			CountryCode:      req.CountryCode,
			Timezone:         req.Timezone,
		}

		return transport.NewHTTPResponse(http.StatusCreated, resp), nil
	}

	h.rh.TryAction(w, action)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("id", "invalid user id")
		}

		userEntity, err := h.repo.GetUserByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return transport.HTTPResponse{}, validation.NewValidationError("id", "user not found")
			}
			return transport.HTTPResponse{}, err
		}

		resp := Response{
			ID:               userEntity.ID,
			Name:             userEntity.Name,
			Surname:          userEntity.Surname,
			Role:             userEntity.Role,
			SubscriptionTier: userEntity.SubscriptionTier,
			Age:              int(userEntity.Age.Int32),
			CountryCode:      userEntity.CountryCode.String,
			Timezone:         userEntity.Timezone.String,
		}

		return transport.NewHTTPResponse(http.StatusOK, resp), nil
	}

	h.rh.TryAction(w, action)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("id", "invalid user id")
		}

		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("body", "invalid request body")
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		errCh := make(chan *validation.AggregateError)
		go h.validator.ValidateCreateOrUpdate(req, errCh)
		if err := <-errCh; err != nil {
			return transport.HTTPResponse{}, err
		}

		params := db.UpdateUserParams{
			ID:      id,
			Name:    req.Name,
			Surname: req.Surname,
		}
		if req.Age > 0 {
			params.Age = sql.NullInt32{Int32: int32(req.Age), Valid: true}
		}
		if req.CountryCode != "" {
			params.CountryCode = sql.NullString{String: req.CountryCode, Valid: true}
		}
		if req.Timezone != "" {
			params.Timezone = sql.NullString{String: req.Timezone, Valid: true}
		}

		if err := h.repo.UpdateUser(r.Context(), params); err != nil {
			return transport.HTTPResponse{}, err
		}

		userEntity, err := h.repo.GetUserByID(r.Context(), id)
		if err != nil {
			return transport.HTTPResponse{}, err
		}

		resp := Response{
			ID:               id,
			Name:             req.Name,
			Surname:          req.Surname,
			Role:             userEntity.Role,
			SubscriptionTier: userEntity.SubscriptionTier,
			Age:              req.Age,
			CountryCode:      req.CountryCode,
			Timezone:         req.Timezone,
		}

		return transport.NewHTTPResponse(http.StatusOK, resp), nil
	}

	h.rh.TryAction(w, action)
}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("id", "invalid user id")
		}

		var req PatchUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return transport.HTTPResponse{}, validation.NewValidationError("body", "invalid request body")
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		errCh := make(chan *validation.AggregateError)
		go h.validator.ValidatePatch(req, errCh)
		if err := <-errCh; err != nil {
			return transport.HTTPResponse{}, err
		}

		userEntity, err := h.repo.GetUserByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return transport.HTTPResponse{}, validation.NewValidationError("id", "user not found")
			}
			return transport.HTTPResponse{}, err
		}

		params := db.UpdateUserParams{
			ID:          id,
			Name:        userEntity.Name,
			Surname:     userEntity.Surname,
			Age:         userEntity.Age,
			CountryCode: userEntity.CountryCode,
			Timezone:    userEntity.Timezone,
		}

		if req.Name != nil {
			params.Name = *req.Name
		}
		if req.Surname != nil {
			params.Surname = *req.Surname
		}
		if req.Age != nil {
			params.Age = sql.NullInt32{Int32: int32(*req.Age), Valid: true}
		}
		if req.CountryCode != nil {
			params.CountryCode = sql.NullString{String: *req.CountryCode, Valid: true}
		}
		if req.Timezone != nil {
			params.Timezone = sql.NullString{String: *req.Timezone, Valid: true}
		}

		if err := h.repo.UpdateUser(r.Context(), params); err != nil {
			return transport.HTTPResponse{}, err
		}

		resp := Response{
			ID:               id,
			Name:             params.Name,
			Surname:          params.Surname,
			Role:             userEntity.Role,
			SubscriptionTier: userEntity.SubscriptionTier,
			Age:              int(params.Age.Int32),
			CountryCode:      params.CountryCode.String,
			Timezone:         params.Timezone.String,
		}

		return transport.NewHTTPResponse(http.StatusOK, resp), nil
	}

	h.rh.TryAction(w, action)
}
