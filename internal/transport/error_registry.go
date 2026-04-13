package transport

import (
	"github.com/brota/gobackend/internal/custom_errors"
	errconsts "github.com/brota/gobackend/internal/custom_errors/abstract_error_code"
	"github.com/brota/gobackend/internal/custom_errors/transport"
)

type ContextualErrorHandler = func(error, map[string]any) HTTPResponse

type ErrorCodeProvider interface {
	Code() errconsts.ErrorCode
}

type ErrorTranslator interface {
	Translate(err error) HTTPResponse
}

type ErrorRegistry struct {
	handlers map[errconsts.ErrorCode]ContextualErrorHandler
	fallback ContextualErrorHandler
}

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		handlers: make(map[errconsts.ErrorCode]ContextualErrorHandler),
		fallback: DefaultFallbackHandler,
	}
}

func DefaultFallbackHandler(error, map[string]any) HTTPResponse {
	return HTTPResponse{
		Status: 500,
		Code:   string(transport.InternalErrorCode),
		Detail: "unexpected error",
	}
}

func (r *ErrorRegistry) Register(code errconsts.ErrorCode, handler ContextualErrorHandler) {
	r.handlers[code] = handler
}

func (r *ErrorRegistry) Translate(err error) HTTPResponse {
	code := ExtractErrorCode(err)
	context := extractContext(err)
	handler := r.ResolveHandler(code)
	return handler(err, context)
}

func extractContext(err error) map[string]any {
	if carrier, ok := err.(custom_errors.ContextCarrier); ok {
		return carrier.ContextData()
	}
	return map[string]any{}
}

func ExtractErrorCode(err error) errconsts.ErrorCode {
	if provider, ok := err.(ErrorCodeProvider); ok {
		return provider.Code()
	}
	return transport.UnknownErrorCode
}

func (r *ErrorRegistry) ResolveHandler(code errconsts.ErrorCode) ContextualErrorHandler {
	handler, exists := r.handlers[code]
	if exists {
		return handler
	}
	return r.fallback
}
