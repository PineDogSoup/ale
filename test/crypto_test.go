package test

import (
	"ale/utils"
	"fmt"
	"log"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := utils.GenerateKeyPair()
	if err != nil {
		log.Fatal("Error generating key pair:", err)
	}

	fmt.Println("Private Key:", keyPair.PrivateKey)
	fmt.Println("Public Key:", keyPair.PublicKey)
}
