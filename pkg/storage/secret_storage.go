package storage

import "github.com/zalando/go-keyring"

// TODO fix that
func SaveSecret(key, secret string) error {
	return keyring.Set("ConfigWatcher", key, secret)
}

func GetSecret(key string) (string, error) {
	secret, err := keyring.Get("ConfigWatcher", key)
	if err != nil {
		return "", err
	}
	return secret, nil
}
