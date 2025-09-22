package pgp

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// GetPGPSecretKey gets a secret key by it's key ID and unlocks it
// If keyID is nil, then the default/first secret key will be used
func GetPGPSecretKey(ctx context.Context, keyID *string) (*crypto.Key, error) {
	secretKeyID, password, err := getPGPSecretKeyIDAndPassword(keyID)
	if err != nil {
		panic(err)
	}

	cmd := exec.CommandContext(ctx, "gpg", "--export-secret-keys", "--armor", "--pinentry-mode", "loopback", "--passphrase-fd", "0", secretKeyID)
	cmd.Stdin = strings.NewReader(password)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("could not export secret key: %s", output), err)
	}

	lockedKey, err := crypto.NewKeyFromArmored(string(output))
	if err != nil {
		return nil, err
	}

	return lockedKey.Unlock([]byte(password))
}
