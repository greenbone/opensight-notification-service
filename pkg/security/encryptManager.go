package security

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/greenbone/opensight-golang-libraries/pkg/dbcrypt"
	"github.com/greenbone/opensight-notification-service/pkg/config"
	"github.com/rs/zerolog/log"
)

type Key struct {
	Password     string
	PasswordSalt string
}

type EncryptManager struct {
	activeID int
	keys     map[int]Key
	mu       sync.RWMutex
}

func NewEncryptManager() *EncryptManager {
	return &EncryptManager{}
}

func (sm *EncryptManager) UpdateKeys(keyringConfig config.DatabaseKeyringConfig) {
	newKeys := make(map[int]Key)
	for id, key := range keyringConfig.Keys {
		// TODO add validations before updating
		newKeys[id] = Key{
			Password:     key.Password,
			PasswordSalt: key.PasswordSalt,
		}
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.keys = newKeys
	sm.activeID = keyringConfig.ActiveID
	log.Info().Msgf("Keyring successfully refreshed in memory")
}

func (sm *EncryptManager) StartRefresher(ctx context.Context, keyringConfig config.DatabaseKeyringConfig, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sm.UpdateKeys(keyringConfig)
		}
	}
}

func (sm *EncryptManager) Encrypt(plaintext string) ([]byte, int, error) {
	sm.mu.RLock()
	currentVersion := sm.activeID
	activeKey := sm.keys[currentVersion]
	sm.mu.RUnlock()

	cipher, err := dbcrypt.NewDBCipher(dbcrypt.Config{
		Password:     activeKey.Password,
		PasswordSalt: activeKey.PasswordSalt,
	})
	if err != nil {
		return nil, 0, err
	}

	encryptedBytes, err := cipher.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, 0, err
	}

	return encryptedBytes, currentVersion, nil
}

func (sm *EncryptManager) Decrypt(data []byte, version int) (string, error) {
	sm.mu.RLock()
	activeKey, ok := sm.keys[version]
	sm.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("unknown salt version: %d", version)
	}

	cipher, err := dbcrypt.NewDBCipher(dbcrypt.Config{
		Password:     activeKey.Password,
		PasswordSalt: activeKey.PasswordSalt,
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
