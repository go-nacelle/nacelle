package grpc

import (
	"context"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/examples/multi-service/secret"
	"google.golang.org/grpc"
)

type EndpointSet struct {
	Logger        nacelle.Logger       `service:"logger"`
	SecretService secret.SecretService `service:"secret-service"`
}

func NewEndpointSet() *EndpointSet {
	return &EndpointSet{}
}

func (es *EndpointSet) Init(config nacelle.Config, s *grpc.Server) error {
	RegisterSecretServiceServer(s, es)
	return nil
}

func (es *EndpointSet) PostSecret(ctx context.Context, secret *Secret) (*Id, error) {
	id, err := es.SecretService.Post(secret.Secret)
	if err != nil {
		return nil, err
	}

	return &Id{Name: id}, nil
}

func (es *EndpointSet) ReadSecret(ctx context.Context, id *Id) (*Secret, error) {
	data, err := es.SecretService.Read(id.Name)
	if err != nil {
		if err != secret.ErrNoSecret {
			return nil, err
		}

		return nil, nil
	}

	return &Secret{Secret: data}, nil
}
