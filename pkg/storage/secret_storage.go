package storage

import (
	"github.com/zalando/go-keyring"
)

type SecretStorage struct {
	keyring.Keyring
	service_name string
}

const defaultServiceName = "ConfigWatcher"

func NewSecretStorage(serviceName string) *SecretStorage {
	if serviceName == "" {
		serviceName = defaultServiceName
	}
	return &SecretStorage{
		service_name: serviceName,
	}
}

func (sec *SecretStorage) SaveSecret(key, secret string) error {
	err := sec.Set(sec.service_name, key, secret)
	if err != nil {
		return err
	}
	return nil
}

func (sec *SecretStorage) GetSecret(key string) (string, error) {
	secret, err := sec.Get(sec.service_name, key)
	if err != nil {
		return "", err
	}
	return secret, nil
}
