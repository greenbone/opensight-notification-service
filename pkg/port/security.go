package port

type EncryptManager interface {
	Encrypt(plaintext string) ([]byte, int, error)
	Decrypt(data []byte, version int) (string, error)
}
