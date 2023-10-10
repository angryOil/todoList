package handler

import "net/http"

type DecoratorHandleFunc func(w http.ResponseWriter, r *http.Request, h http.Handler)

type DecoHandler struct {
	dhf DecoratorHandleFunc
	h   http.Handler
}

func (dh *DecoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dh.dhf(w, r, dh.h)
}

func NewDecoHandler(h http.Handler, dhf DecoratorHandleFunc) http.Handler {
	return &DecoHandler{
		dhf: dhf,
		h:   h,
	}
}
