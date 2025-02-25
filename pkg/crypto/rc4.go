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

// Decrypt decrypts the packet data, skipping the first 5 bytes (4 bytes length + 1 byte ID)
func (m *RC4Manager) Decrypt(data []byte) {
	// Skip encryption for packets smaller than 5 bytes
	if len(data) <= 5 {
		return
	}

	// Only decrypt the payload (skip the 4-byte length and 1-byte ID)
	m.inbound.XORKeyStream(data[5:], data[5:])
}

// Encrypt encrypts the packet data, skipping the first 5 bytes (4 bytes length + 1 byte ID)
func (m *RC4Manager) Encrypt(data []byte) {
	// Skip encryption for packets smaller than 5 bytes
	if len(data) <= 5 {
		return
	}

	// Only encrypt the payload (skip the 4-byte length and 1-byte ID)
	m.outbound.XORKeyStream(data[5:], data[5:])
}
