package pgp

import (
	"bytes"
	"errors"
	"os/exec"
)

var (
	ErrPGPKeyNotFound = errors.New("PGP key not found")
)

func getKeyIDForKeygrip(keygrip string) (string, error) {
	output, err := exec.Command("gpg", "--list-keys", "--with-colons", "--with-keygrip").Output()
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(output, []byte("\n"))
	for i, line := range lines {
		if bytes.Contains(line, []byte(keygrip)) {
			for j := i - 1; j >= 0; j-- {
				if bytes.HasPrefix(lines[j], []byte("pub:")) {
					return string(bytes.Split(lines[j], []byte(":"))[4]), nil
				}
			}
		}
	}

	return "", ErrPGPKeyNotFound
}
