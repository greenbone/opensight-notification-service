package security

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/greenbone/opensight-notification-service/pkg/config"
)

func TestSaltManager_UpdateKeys(t *testing.T) {
	t.Parallel()

	cfg := config.DatabaseEncryptionKey{
		Password:     "password",
		PasswordSalt: "password-salt",
	}

	mgr := NewEncryptManager()
	mgr.UpdateKeys(cfg)

	em, ok := mgr.(*encryptManager)
	if !ok {
		t.Fatalf("mgr is not of type *encryptManager")
	}

	assert.Equal(t, em.activeKey.Password, cfg.Password)
	assert.Equal(t, em.activeKey.PasswordSalt, cfg.PasswordSalt)
}

func TestSaltManager_Encrypt_Decrypt(t *testing.T) {

	configKey := config.DatabaseEncryptionKey{
		Password:     "password",
		PasswordSalt: "password-salt-should-no-be-short-fyi",
	}

	t.Run("Encrypt and Decrypt success", func(t *testing.T) {
		t.Parallel()

		mgr := NewEncryptManager()
		mgr.UpdateKeys(configKey)

		plaintext := "password"

		encryptedPwd, err := mgr.Encrypt(plaintext)
		assert.NotEmpty(t, encryptedPwd)
		assert.NoError(t, err)

		decryptedPwd, err := mgr.Decrypt(encryptedPwd)
		assert.Equal(t, plaintext, decryptedPwd)
		assert.NoError(t, err)
	})

	t.Run("Encrypt fails when empty plaintext is passed", func(t *testing.T) {
		t.Parallel()

		mgr := NewEncryptManager()
		mgr.UpdateKeys(configKey)

		encryptedPwd, err := mgr.Encrypt("")
		assert.Empty(t, encryptedPwd)
		assert.Error(t, err)
	})

	t.Run("Encryption fails when salt is not secure", func(t *testing.T) {
		t.Parallel()

		mgr := NewEncryptManager()
		mgr.UpdateKeys(config.DatabaseEncryptionKey{
			Password:     "password",
			PasswordSalt: "weak",
		})

		encryptedPwd, err := mgr.Encrypt("something")
		assert.NotEmpty(t, err)
		assert.Empty(t, encryptedPwd)
	})

	t.Run("Decryption fails when passed value is empty", func(t *testing.T) {
		t.Parallel()

		mgr := NewEncryptManager()
		mgr.UpdateKeys(configKey)

		encryptedPwd, err := mgr.Decrypt([]byte(``))
		assert.NotEmpty(t, err)
		assert.Empty(t, encryptedPwd)
	})
}
