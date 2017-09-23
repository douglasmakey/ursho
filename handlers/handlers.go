// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/douglasmakey/ursho/storages"
)

func New(prefix string, storage storages.IFStorage) http.Handler {
	mux := http.NewServeMux()
	h := handler{prefix, storage}
	mux.HandleFunc("/encode/", h.encode)
	mux.HandleFunc("/", h.redirect)
	mux.HandleFunc("/info/", h.decode)
	return mux
}

type bodyRequest struct {
	URL string
}

type handler struct {
	prefix  string
	storage storages.IFStorage
}

func (h handler) encode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var b bodyRequest
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := Response{Data: "Unable to decode JSON request body: " + err.Error(), Success: false}
		createResponse(w, e)
		return
	}

	b.URL = strings.TrimSpace(b.URL)

	if b.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		e := Response{Data: "URL is Empty", Success: false}
		createResponse(w, e)
		return
	}

	c, err := h.storage.Save(b.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := Response{Data: err.Error(), Success: false}
		createResponse(w, e)
		return
	}

	response := Response{Data: h.prefix + c, Success: true}
	createResponse(w, response)
}

func (h handler) decode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	code := r.URL.Path[len("/info/"):]

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		e := Response{Data: "URL Not Found", Success: false}
		createResponse(w, e)
		return
	}

	response := Response{Data: model, Success: true}
	createResponse(w, response)
}

func (h handler) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Path[len("/"):]

	model, err := h.storage.Load(code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
		return
	}

	http.Redirect(w, r, string(model.Url), 301)
}
