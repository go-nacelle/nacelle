package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	"github.com/gorilla/mux"

	"github.com/efritz/nacelle/examples/multi-service/secret"
)

type EndpointSet struct {
	Logger        nacelle.Logger       `service:"logger"`
	SecretService secret.SecretService `service:"secret-service"`
}

func NewEndpointSet() *EndpointSet {
	return &EndpointSet{}
}

func (es *EndpointSet) Init(config nacelle.Config, server *http.Server) error {
	router := mux.NewRouter()
	server.Handler = router

	// Register routes
	router.HandleFunc("/post", response.Convert(es.post)).Methods("POST")
	router.HandleFunc("/load/{id}", response.Convert(es.load)).Methods("GET")

	return nil
}

func (es *EndpointSet) post(r *http.Request) response.Response {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		es.Logger.Error("failed to read payload (%s)", err.Error())
		return response.Empty(http.StatusInternalServerError)
	}

	id, err := es.SecretService.Post(string(data))
	if err != nil {
		es.Logger.Error("failed to post secret (%s)", err.Error())
		return response.Empty(http.StatusInternalServerError)
	}

	es.Logger.Info("Posted secret %s", id)

	resp := response.Empty(http.StatusOK)
	resp.AddHeader("Location", fmt.Sprintf("/load/%s", id))
	return resp
}

func (es *EndpointSet) load(r *http.Request) response.Response {
	id := mux.Vars(r)["id"]

	data, err := es.SecretService.Read(id)
	if err != nil {
		if err == secret.ErrNoSecret {
			es.Logger.Info("Secret %s requested but not found", id)
			return response.Empty(http.StatusNotFound)
		}

		es.Logger.Error("failed to retrieve secret (%s)", err.Error())
		return response.Empty(http.StatusInternalServerError)
	}

	es.Logger.Info("Fetched secret %s", id)
	return response.Respond([]byte(data))
}
