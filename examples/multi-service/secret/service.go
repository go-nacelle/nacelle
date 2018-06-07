package secret

import (
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/efritz/overcurrent"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
)

type (
	SecretService interface {
		Post(string) (string, error)
		Read(string) (string, error)
	}

	secretService struct {
		ttl             int
		prefix          string
		client          deepjoy.Client
		breakerRegistry overcurrent.Registry
	}
)

var ErrNoSecret = fmt.Errorf("secret does not exist")

func (s *secretService) Post(secret string) (string, error) {
	var (
		id  = uuid.Must(uuid.NewV4()).String()
		key = s.key(id)
	)

	_, err := s.client.Transaction(
		deepjoy.NewCommand("set", key, secret),
		deepjoy.NewCommand("expire", key, s.ttl),
	)

	return id, err
}

func (s *secretService) Read(id string) (string, error) {
	var (
		key    = s.key(id)
		result interface{}
		err    error
	)

	result, err = s.client.Transaction(
		deepjoy.NewCommand("get", key),
		deepjoy.NewCommand("del", key),
	)

	if err != nil {
		return "", err
	}

	if result := result.([]interface{})[0]; result != nil {
		secret, err := redis.String(result, nil)
		if err != nil {
			return "", err
		}

		return secret, nil
	}

	return "", ErrNoSecret
}

func (s *secretService) key(id string) string {
	return fmt.Sprintf("%s:%s", s.prefix, id)
}
