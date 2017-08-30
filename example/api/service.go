package api

import (
	"errors"
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/satori/go.uuid"
)

type (
	SecretService interface {
		Post(string) (string, error)
		Read(string) (string, error)
	}

	secretService struct {
		ttl    int
		prefix string
		client deepjoy.Client
	}
)

var ErrNoSecret = errors.New("secret does not exist")

func (s *secretService) Post(secret string) (string, error) {
	var (
		id  = uuid.NewV4().String()
		key = s.key(id)
	)

	_, err := s.client.Transaction(
		deepjoy.NewCommand("set", key, secret),
		deepjoy.NewCommand("expire", key, s.ttl),
	)

	return id, err
}

func (s *secretService) Read(id string) (string, error) {
	key := s.key(id)
	result, err := s.client.Transaction(
		deepjoy.NewCommand("get", key),
		deepjoy.NewCommand("del", key),
	)

	if err != nil {
		return "", err
	}

	if val := result.([]interface{})[0]; val != nil {
		if arr, ok := val.([]byte); ok {
			return string(arr), nil
		}
	}

	return "", ErrNoSecret
}

func (s *secretService) key(id string) string {
	return fmt.Sprintf("%s:%s", s.prefix, id)
}
