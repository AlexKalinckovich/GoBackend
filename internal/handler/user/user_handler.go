package handler

import (
	"encoding/json"
	"github.com/brota/gobackend/internal/domain/user"
	"github.com/brota/gobackend/internal/handler"
	service "github.com/brota/gobackend/internal/service/user"
	"net/http"
)

type UserHandler struct {
	service        service.UserService
	requestHandler handler.RequestHandler
}

func NewUserHandler(service service.UserService, requestHandler handler.RequestHandler) *UserHandler {
	return &UserHandler{service: service, requestHandler: requestHandler}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	action := h.buildCreateAction(r)
	h.requestHandler.TryAction(w, action)
}

func (h *UserHandler) buildCreateAction(r *http.Request) handler.Action {
	return func() error {
		return h.decodeAndProcess(r)
	}
}

func (h *UserHandler) decodeAndProcess(r *http.Request) error {
	var entity user.Entity
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&entity)
	if err != nil {
		return err
	}
	return h.service.CreateUser(entity)
}
