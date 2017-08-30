package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/gorilla/mux"

	"github.com/efritz/nacelle/example/api"
)

type handlerSet struct {
	logger        nacelle.Logger
	secretService api.SecretService
}

func (h *handlerSet) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(nil, "%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := h.secretService.Post(string(data))
	if err != nil {
		h.logger.Error(nil, "%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info(nil, "Posted secret %s", id)
	w.Header().Set("Location", fmt.Sprintf("/load/%s", id))
	w.WriteHeader(http.StatusOK)
}

func (h *handlerSet) load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	secret, err := h.secretService.Read(id)
	if err != nil {
		if err == api.ErrNoSecret {
			h.logger.Info(nil, "Secret %s requested but not found", id)
			w.WriteHeader(http.StatusNotFound)
		} else {
			h.logger.Error(nil, "%s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	h.logger.Info(nil, "Fetched secret %s", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(secret))
}
