package main

import (
	"errors"
	"github.com/brota/gobackend/internal/constants/errors/domain"
	"github.com/brota/gobackend/internal/domainerrors"
	"github.com/brota/gobackend/internal/handler"
	"github.com/brota/gobackend/internal/transport"
	"net/http"
)

func main() {
	registry := transport.NewErrorRegistry()
	registerDefaultHandlers(registry)
	requestHandler := handler.NewRequestHandler(registry)
	http.HandleFunc("/test", createEndpoint(requestHandler))
	_ = http.ListenAndServe(":8080", nil)
}

func registerDefaultHandlers(registry *transport.ErrorRegistry) {
	registry.Register(domain.ValidationErrorCode, handleValidationError)
	registry.Register(domain.BusinessLogicErrorCode, handleBusinessLogicError)
	registry.Register(domain.ConstraintErrorCode, handleConstraintError)
}

func handleValidationError(err error) transport.HTTPResponse {
	var ve *domainerrors.ValidationError
	errors.As(err, &ve)
	return transport.HTTPResponse{Status: 422, Code: string(ve.Code()), Detail: ve.Field() + ": " + ve.Error()}
}

func handleBusinessLogicError(err error) transport.HTTPResponse {
	var businessLogicError *domainerrors.BusinessLogicError
	errors.As(err, &businessLogicError)
	return transport.HTTPResponse{
		Status: 409,
		Code:   string(businessLogicError.Code()),
		Detail: businessLogicError.Error(),
	}
}

func handleConstraintError(err error) transport.HTTPResponse {
	var ce *domainerrors.ConstraintError
	errors.As(err, &ce)
	return transport.HTTPResponse{Status: 403, Code: string(ce.Code()), Detail: ce.Error()}
}

func createEndpoint(h *handler.RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.TryAction(w, simulateBusinessAction)
	}
}

func simulateBusinessAction() error {
	return domainerrors.NewBusinessLogicError("insufficient funds")
}
