package secret

import (
	"fmt"

	"github.com/efritz/deepjoy"
	"github.com/efritz/overcurrent"
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
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
		id  = uuid.New().String()
		key = s.key(id)
	)

	pipeline := s.client.Pipeline()
	pipeline.Add("set", key, secret)
	pipeline.Add("expire", key, s.ttl)
	_, err := pipeline.Run()

	return id, err
}

func (s *secretService) Read(id string) (string, error) {
	var (
		key    = s.key(id)
		result interface{}
		err    error
	)

	pipeline := s.client.Pipeline()
	pipeline.Add("get", key)
	pipeline.Add("del", key)

	result, err = pipeline.Run()
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
