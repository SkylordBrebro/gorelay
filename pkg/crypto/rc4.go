package crypto

import (
	"crypto/rc4"
)

// RC4Manager handles packet encryption/decryption
type RC4Manager struct {
	inbound  *rc4.Cipher
	outbound *rc4.Cipher
}

func NewRC4Manager(inKey, outKey []byte) (*RC4Manager, error) {
	inCipher, err := rc4.NewCipher(inKey)
	if err != nil {
		return nil, err
	}

	outCipher, err := rc4.NewCipher(outKey)
	if err != nil {
		return nil, err
	}

	return &RC4Manager{
		inbound:  inCipher,
		outbound: outCipher,
	}, nil
}

func (m *RC4Manager) Decrypt(data []byte) {
	m.inbound.XORKeyStream(data, data)
}

func (m *RC4Manager) Encrypt(data []byte) {
	m.outbound.XORKeyStream(data, data)
}
