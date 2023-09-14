package utils

import (
	"ale/core/types"
	pb "ale/protobuf/generated"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	secp256 "github.com/haltingstate/secp256k1-go"
	"log"
	"reflect"
	"unsafe"
)

// Base58StringToAddress address to  bytes.
func Base58StringToAddress(addr string) (*pb.Address, error) {
	var address = new(pb.Address)
	addressBytes, err := DecodeCheck(addr)
	if err != nil {
		return nil, errors.New("Base58String To Address error")
	}
	address.Value = addressBytes
	return address, nil
}

func AddressToBase58String(address *pb.Address) string {
	return EncodeCheck(address.Value)
}

// BytesToString  Bytes To String.
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

// GetBytesSha256 Get Bytes Sha256.
func GetBytesSha256(s string) []byte {
	sha := sha256.New()
	sha.Write([]byte(s))
	return sha.Sum(nil)
}

// BytesToInt Bytes To Int.
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var data int32
	binary.Read(bytesBuffer, binary.LittleEndian, &data)
	return int(data)
}

// EncodeCheck prepends appends a four byte checksum.
func EncodeCheck(input []byte) string {
	b := make([]byte, 0, 1+len(input)+4)
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return base58.Encode(b)
}

// checksum: first four bytes of sha256^2.
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

// DecodeCheck decodes a string that was encoded with CheckEncode and verifies the checksum.
func DecodeCheck(input string) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, errors.New("invalid format: version and/or checksum bytes missing")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, errors.New("checksum error")
	}
	payload := decoded[0 : len(decoded)-4]
	result = append(result, payload...)
	return
}

// GetAddressByBytes  sha256 twice Bytes to get address.
func GetAddressByBytes(b []byte) string {
	firstBytes := sha256.Sum256(b)
	secondBytes := sha256.Sum256(firstBytes[:])
	address := EncodeCheck(secondBytes[:])
	return address
}

// Base64DecodeBytes Base64 Decode Bytes
func Base64DecodeBytes(param string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		fmt.Println("json Marshal error")
	}
	return b, nil
}

// SignWithPrivateKey Get Signature With PrivateKey.
func SignWithPrivateKey(privateKey string, txData []byte) (string, error) {
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	txDataBytes := sha256.Sum256(txData)
	signatureBytes := secp256.Sign(txDataBytes[:], privateKeyBytes)
	return hex.EncodeToString(signatureBytes), nil
}

// GetAddressFromPrivateKey Get the account address through the private key.
func GetAddressFromPrivateKey(privateKey string) string {
	bytes, _ := hex.DecodeString(privateKey)
	pubkeyBytes := secp256.UncompressedPubkeyFromSeckey(bytes)
	return GetAddressByBytes(pubkeyBytes)
}

// GenerateKeyPairInfo Generate KeyPair Info.
func GenerateKeyPairInfo() *types.KeyPair {
	publicKeyBytes, privateKeyBytes := secp256.GenerateKeyPair()
	publicKey := hex.EncodeToString(publicKeyBytes)
	privateKey := hex.EncodeToString(privateKeyBytes)
	privateKeyAddress := GetAddressFromPrivateKey(privateKey)
	var keyPair = &types.KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    privateKeyAddress,
	}
	return keyPair
}

func HexStringToByteArray(hexString string) []byte {
	if len(hexString) >= 2 && hexString[0] == '0' && (hexString[1] == 'x' || hexString[1] == 'X') {
		hexString = hexString[2:]
	}
	length := len(hexString)
	byteArray := make([]byte, length/2)
	for startIndex := 0; startIndex < length; startIndex += 2 {
		byteValue, err := hex.DecodeString(hexString[startIndex : startIndex+2])
		if err != nil {
			log.Fatal(err)
		}
		byteArray[startIndex/2] = byteValue[0]
	}
	return byteArray
}
