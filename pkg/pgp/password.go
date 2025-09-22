package pgp

import (
	"errors"
	"strings"

	"github.com/gsterjov/go-libsecret"
)

const (
	gpgLabelPrefix = "GnuPG: n/"
)

var (
	ErrNoPGPKeysFound = errors.New("no PGP keys found")
)

func getPGPSecretKeyIDAndPassword(keyID *string) (secretKeyID string, password string, err error) {
	service, err := libsecret.NewService()
	if err != nil {
		return "", "", err
	}

	session, err := service.Open()
	if err != nil {
		return "", "", err
	}

	collections, err := service.Collections()
	if err != nil {
		return "", "", err
	}

	secretKeyCandidates := 0
	for _, collection := range collections {
		if err := service.Unlock(&collection); err != nil {
			return "", "", err
		}

		items, err := collection.Items()
		if err != nil {
			return "", "", err
		}

		for _, item := range items {
			label, err := item.Label()
			if err != nil {
				return "", "", err
			}

			if !strings.HasPrefix(label, gpgLabelPrefix) {
				continue
			}

			keygrip := strings.TrimPrefix(label, gpgLabelPrefix)

			secretKeyID, err := getKeyIDForKeygrip(keygrip)
			if err != nil {
				return "", "", err
			}

			if keyID != nil && *keyID != secretKeyID {
				continue
			}

			secret, err := item.GetSecret(session)
			if err != nil {
				return "", "", err
			}

			return secretKeyID, string(secret.Value), nil
		}
	}

	if secretKeyCandidates <= 0 {
		return "", "", ErrNoPGPKeysFound
	}

	return "", "", ErrPGPKeyNotFound
}
