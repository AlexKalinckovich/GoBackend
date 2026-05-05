package handler

import (
	"github.com/brota/gobackend/internal/shared/transport"
	"net/http"

	transporterror "github.com/brota/gobackend/internal/shared/custom_errors/transport"
)

type TestErrorHandler struct {
	rh *RequestHandler
}

func NewTestErrorHandler() *TestErrorHandler {
	registry := transport.NewErrorRegistry()

	registry.Register(transporterror.InternalErrorCode, func(err error, m map[string]any) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, err.Error())
	})

	rh := NewRequestHandler(registry)
	return &TestErrorHandler{rh: rh}
}

func (th *TestErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		return transport.HTTPResponse{}, transporterror.NewInternalError("TEST_ERROR")
	}
	th.rh.TryAction(w, action)
}
