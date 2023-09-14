package utils

import (
	"ale/core/crypto"
	"crypto/rand"
	"github.com/btcsuite/btcd/btcec"
)

func GenerateKeyPair() (*crypto.ECKeyPair, error) {
	privateKeyBytes := make([]byte, 32)

	// Generate random private key
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Generate public key from private key
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	publicKey := privateKey.PubKey()

	// Serialize public key
	serializedPublicKey := publicKey.SerializeCompressed()

	keyPair := &crypto.ECKeyPair{
		PublicKey:  serializedPublicKey,
		PrivateKey: privateKeyBytes,
	}

	return keyPair, nil
}
