package port

type EncryptManager interface {
	Encrypt(plaintext string) ([]byte, error)
	Decrypt(data []byte) (string, error)
}
