package security

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/greenbone/opensight-notification-service/pkg/config"
)

func TestSaltManager_UpdateKeys(t *testing.T) {

	mgr := NewEncryptManager()

	// TODO enhance the test cases

	tests := map[string]struct {
		cfgValues config.DatabaseKeyringConfig
		wantErr   bool
	}{
		"Update keys": {
			cfgValues: config.DatabaseKeyringConfig{
				ActiveID: 1,
				Keys: map[int]config.EncryptionKey{
					1: {
						Password:     "password",
						PasswordSalt: "password-salt",
					},
				},
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mgr.UpdateKeys(tc.cfgValues)
			assert.Equal(t, mgr.activeID, tc.cfgValues.ActiveID)

			for k, v := range tc.cfgValues.Keys {
				assert.Equal(t, mgr.keys[k].Password, v.Password)
				assert.Equal(t, mgr.keys[k].PasswordSalt, v.PasswordSalt)
			}
		})
	}
}

func TestSaltManager_Encrypt_Decrypt(t *testing.T) {

	mgr := NewEncryptManager()
	mgr.UpdateKeys(config.DatabaseKeyringConfig{
		ActiveID: 1,
		Keys: map[int]config.EncryptionKey{
			1: {
				Password:     "password",
				PasswordSalt: "password-salt-lorem-ipsum-dolor-salt",
			},
		},
	})

	// TODO enhance the test cases
	// db password salt is too short

	tests := map[string]struct {
		inValue    string
		keyVersion int
		wantErr    bool
	}{
		"Encrypt and Decrypt success": {
			inValue:    "password",
			keyVersion: 1,
			wantErr:    false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encryptedPwd, version, err := mgr.Encrypt(tc.inValue)
			assert.NotZero(t, version)
			assert.NotEmpty(t, encryptedPwd)
			assert.NoError(t, err)

			decryptedPwd, err := mgr.Decrypt(encryptedPwd, tc.keyVersion)
			assert.Equal(t, tc.inValue, decryptedPwd)
			assert.NoError(t, err)
		})
	}
}
