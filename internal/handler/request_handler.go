package handler

import (
	"encoding/json"
	"github.com/brota/gobackend/internal/transport"
	"net/http"
)

type Action func() error
type ResultChan chan transport.HTTPResponse

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
	err := action()
	return h.translator.Translate(err)
}

func (h *RequestHandler) WriteResponse(w http.ResponseWriter, resp transport.HTTPResponse) {
	w.WriteHeader(resp.Status)
	h.setContentType(w)
	h.encodeResponse(w, resp)
}

func (h *RequestHandler) setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (h *RequestHandler) encodeResponse(w http.ResponseWriter, resp transport.HTTPResponse) {
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RequestHandler) TryActionAsync(action Action) ResultChan {
	ch := make(ResultChan, 1)
	go h.runAsyncTask(action, ch)
	return ch
}

func (h *RequestHandler) runAsyncTask(action Action, ch ResultChan) {
	result := h.CatchError(action)
	ch <- result
	close(ch)
}
