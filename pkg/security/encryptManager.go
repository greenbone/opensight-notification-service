package security

import (
	"errors"
	"fmt"
	"strings"

	"github.com/greenbone/opensight-golang-libraries/pkg/dbcrypt"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/rs/zerolog/log"
)

type Key struct {
	Password     string
	PasswordSalt string
}

type EncryptManager interface {
	UpdateKeys(keyringConfig config.DatabaseEncryptionKey)
	Encrypt(plaintext string) ([]byte, error)
	Decrypt(data []byte) (string, error)
}

type encryptManager struct {
	activeKey Key
}

func NewEncryptManager() EncryptManager {
	return &encryptManager{}
}

func (sm *encryptManager) UpdateKeys(keyringConfig config.DatabaseEncryptionKey) {
	if strings.TrimSpace(keyringConfig.Password) == "" {
		log.Error().Msg("Empty password for keyring")
	}

	if strings.TrimSpace(keyringConfig.PasswordSalt) == "" {
		log.Error().Msg("Empty password_salt for keyring")
	}

	sm.activeKey = Key{
		Password:     keyringConfig.Password,
		PasswordSalt: keyringConfig.PasswordSalt,
	}

	log.Info().Msgf("Keyring successfully refreshed in memory")
}

func (sm *encryptManager) Encrypt(plaintext string) ([]byte, error) {
	if len(strings.TrimSpace(plaintext)) == 0 {
		return nil, errors.New("plaintext must be a value")
	}

	cipher, err := dbcrypt.NewDBCipher(dbcrypt.Config{
		Password:     sm.activeKey.Password,
		PasswordSalt: sm.activeKey.PasswordSalt,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create a new cipher instance: %w", err)
	}

	encryptedBytes, err := cipher.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func (sm *encryptManager) Decrypt(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("data must be a value")
	}

	cipher, err := dbcrypt.NewDBCipher(dbcrypt.Config{
		Password:     sm.activeKey.Password,
		PasswordSalt: sm.activeKey.PasswordSalt,
	})
	if err != nil {
		return "", err
	}

	decryptedBytes, err := cipher.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
