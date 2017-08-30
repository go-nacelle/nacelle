package grpc

import (
	"context"

	"github.com/efritz/nacelle/example/api"
)

func (s *Server) PostSecret(ctx context.Context, secret *Secret) (*Id, error) {
	id, err := s.SecretService.Post(secret.Secret)
	if err != nil {
		return nil, err
	}

	return &Id{Name: id}, nil
}

func (s *Server) ReadSecret(ctx context.Context, id *Id) (*Secret, error) {
	secret, err := s.SecretService.Read(id.Name)
	if err != nil {
		if err != api.ErrNoSecret {
			return nil, err
		}

		return nil, nil
	}

	return &Secret{Secret: secret}, nil
}
