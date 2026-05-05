package user

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brota/gobackend/internal/shared/db"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUserWithID(ctx context.Context, params db.CreateUserParams) (int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*db.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*db.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupChiContext(req *http.Request, key string, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestHandler_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	t.Run("Success", func(t *testing.T) {
		reqBody := CreateUserRequest{Name: "John", Surname: "Doe", Age: 25}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		mockRepo.On("CreateUserWithID", mock.Anything, mock.AnythingOfType("db.CreateUserParams")).
			Return(int64(1), nil).Once()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Failed", func(t *testing.T) {
		reqBody := CreateUserRequest{Name: "", Surname: "Doe", Age: 25}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandler_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req = setupChiContext(req, "id", "1")
		rr := httptest.NewRecorder()

		mockUser := &db.User{ID: 1, Name: "John", Surname: "Doe"}
		mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return(mockUser, nil).Once()

		handler.GetUser(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		req = setupChiContext(req, "id", "999")
		rr := httptest.NewRecorder()

		mockRepo.On("GetUserByID", mock.Anything, int64(999)).Return((*db.User)(nil), sql.ErrNoRows).Once()

		handler.GetUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		req = setupChiContext(req, "id", "abc")
		rr := httptest.NewRecorder()

		handler.GetUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandler_PatchUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	t.Run("Success", func(t *testing.T) {
		reqBody := PatchUserRequest{Name: ptr("Johnny")}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(bodyBytes))
		req = setupChiContext(req, "id", "1")
		rr := httptest.NewRecorder()

		mockUser := &db.User{ID: 1, Name: "John", Surname: "Doe"}
		mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return(mockUser, nil).Once()
		mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("db.UpdateUserParams")).Return(nil).Once()

		handler.PatchUser(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		reqBody := PatchUserRequest{Name: ptr("Johnny")}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(bodyBytes))
		req = setupChiContext(req, "id", "1")
		rr := httptest.NewRecorder()

		mockUser := &db.User{ID: 1, Name: "John", Surname: "Doe"}
		mockRepo.On("GetUserByID", mock.Anything, int64(1)).Return(mockUser, nil).Once()
		mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("db.UpdateUserParams")).Return(errors.New("db error")).Once()

		handler.PatchUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestHandler_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	handler := NewUserHandler(mockRepo)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		req = setupChiContext(req, "id", "1")
		rr := httptest.NewRecorder()

		mockRepo.On("DeleteUser", mock.Anything, int64(1)).Return(nil).Once()

		handler.DeleteUser(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
		req = setupChiContext(req, "id", "abc")
		rr := httptest.NewRecorder()

		handler.DeleteUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Repository Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
		req = setupChiContext(req, "id", "2")
		rr := httptest.NewRecorder()

		mockRepo.On("DeleteUser", mock.Anything, int64(2)).Return(errors.New("db error")).Once()

		handler.DeleteUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
	})
}
