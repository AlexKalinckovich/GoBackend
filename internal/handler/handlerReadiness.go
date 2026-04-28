package handler

import (
	"net/http"

	"github.com/brota/gobackend/internal/transport"
)

type ReadinessHandler struct {
	rh *RequestHandler
}

func NewReadinessHandler() *ReadinessHandler {
	registry := transport.NewErrorRegistry()
	rh := NewRequestHandler(registry)
	return &ReadinessHandler{rh: rh}
}

func (h *ReadinessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := func() (transport.HTTPResponse, error) {
		return transport.NewHTTPResponse(http.StatusOK, "OK"), nil
	}
	h.rh.TryAction(w, action)
}
