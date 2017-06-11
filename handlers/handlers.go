// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/douglasmakey/ursho/storages"
)

type bodyRequest struct {
	URL string
}

func EncodeHandler(storage storages.IFStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			w.Header().Set("Content-Type", "application/json")

			var b bodyRequest
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				e := Error{Data: "Unable to decode JSON request body: " + err.Error(), Success: false}
				d := createError(e)
				w.Write(d)
				return
			}

			b.URL = strings.TrimSpace(b.URL)

			if b.URL == "" {
				w.WriteHeader(http.StatusBadRequest)
				e := Error{Data: "URL is Empty", Success: false}
				d := createError(e)
				w.Write(d)
				return
			}

			c, err := storage.Save(b.URL)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				e := Error{Data: err.Error(), Success: false}
				d := createError(e)
				w.Write(d)
				return
			}

			response := Response{Data: prefix + c, Success: true}
			d := createResponse(response)
			w.Write(d)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	return http.HandlerFunc(handleFunc)
}

func DecodeHandler(storage storages.IFStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			code := r.URL.Path[len("/info/"):]

			model, err := storage.LoadInfo(code)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				e := Error{Data: "URL Not Found", Success: false}
				d := createError(e)
				w.Write(d)
				return
			}

			response := Response{Data: model, Success: true}
			d := createResponse(response)
			w.Write(d)

		} else {

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	return http.HandlerFunc(handleFunc)
}

func RedirectHandler(storage storages.IFStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			code := r.URL.Path[len("/"):]

			model, err := storage.Load(code)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("URL Not Found"))
				return
			}

			http.Redirect(w, r, string(model.Url), 301)

		} else {

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	return http.HandlerFunc(handleFunc)
}
