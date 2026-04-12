package main

import (
	"github.com/brota/gobackend/internal/constants/errors/domain"
	"github.com/brota/gobackend/internal/handler"
	handler2 "github.com/brota/gobackend/internal/handler/user"
	service "github.com/brota/gobackend/internal/service/user"
	"github.com/brota/gobackend/internal/transport"
	validationuser "github.com/brota/gobackend/internal/validation/user"
	"net/http"
)

func main() {
	registry := transport.NewErrorRegistry()
	registerErrorHandlers(registry)
	validator := validationuser.NewValidator()
	userService := service.NewUserService(*validator)
	requestHandler := handler.NewRequestHandler(registry)
	userHandler := handler2.NewUserHandler(*userService, *requestHandler)
	http.HandleFunc("/users", userHandler.Create)
	_ = http.ListenAndServe(":8080", nil)
}

func registerErrorHandlers(registry *transport.ErrorRegistry) {
	registry.Register(domain.ValidationAggregateErrorCode, handleValidationAggregate)
	registry.Register(domain.ValidationErrorCode, handleSingleValidation)
}

func handleValidationAggregate(err error, ctx map[string]any) transport.HTTPResponse {
	details := extractValidationDetails(ctx)
	return transport.HTTPResponse{
		Status: 422,
		Code:   string(domain.ValidationAggregateErrorCode),
		Detail: details,
	}
}

func handleSingleValidation(err error, ctx map[string]any) transport.HTTPResponse {
	field := extractStringField(ctx, "field")
	message := extractStringField(ctx, "message")
	return transport.HTTPResponse{
		Status: 422,
		Code:   string(domain.ValidationErrorCode),
		Detail: field + ": " + message,
	}
}

func extractValidationDetails(ctx map[string]any) map[string]string {
	details, ok := ctx["details"].(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return details
}

func extractStringField(ctx map[string]any, key string) string {
	value, ok := ctx[key].(string)
	if !ok {
		return ""
	}
	return value
}
