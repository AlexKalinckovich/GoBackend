package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brota/gobackend/internal/transport"
)

type Action func() (transport.HTTPResponse, error)

type RequestHandler struct {
	translator transport.ErrorTranslator
}

func NewRequestHandler(t transport.ErrorTranslator) *RequestHandler {
	return &RequestHandler{translator: t}
}

func (h *RequestHandler) TryAction(w http.ResponseWriter, action Action) {
	result := h.CatchError(action)
	h.WriteResponse(w, result)
}

func (h *RequestHandler) CatchError(action Action) transport.HTTPResponse {
	data, err := action()
	if err == nil {
		return data
	}
	return h.translator.Translate(err)
}

func (h *RequestHandler) WriteResponse(w http.ResponseWriter, resp transport.HTTPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status())
	_ = json.NewEncoder(w).Encode(resp.Message)
}
