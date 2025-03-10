package crypto

import (
	"crypto/rc4"
)

// RC4Manager handles packet encryption/decryption
type RC4Manager struct {
	inbound  *rc4.Cipher
	outbound *rc4.Cipher
	inKey    []byte
	outKey   []byte
}

func NewRC4Manager(inKey, outKey []byte) (*RC4Manager, error) {
	manager := &RC4Manager{
		inKey:  inKey,
		outKey: outKey,
	}
	if err := manager.Reset(); err != nil {
		return nil, err
	}
	return manager, nil
}

// Reset reinitializes both RC4 ciphers with their original keys
func (m *RC4Manager) Reset() error {
	var err error
	m.inbound, err = rc4.NewCipher(m.inKey)
	if err != nil {
		return err
	}

	m.outbound, err = rc4.NewCipher(m.outKey)
	if err != nil {
		return err
	}
	return nil
}

// Decrypt decrypts the packet data
func (m *RC4Manager) Decrypt(data []byte) {
	// Decrypt the entire payload
	m.inbound.XORKeyStream(data, data)
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
