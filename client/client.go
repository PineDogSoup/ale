package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/hex"
	secp256 "github.com/haltingstate/secp256k1-go"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// AElfClient AElf Client.
type AElfClient struct {
	Host         string
	Version      string
	PrivateKey   string
	Address      string
	HttpClient   *utils.HttpClient
	endpoints    []url.URL
	rand         *rand.Rand
	ContractInfo sync.Map
	sync.RWMutex
}

// const const.
const (
	CHAINSTATUS             = "/api/blockChain/chainStatus"
	BLOCKHEIGHT             = "/api/blockChain/blockHeight"
	BLOCKBYHASH             = "/api/blockChain/block"
	BLOCKBYHEIGHT           = "/api/blockChain/blockByHeight"
	TRANSACTIONPOOLSTATUS   = "/api/blockChain/transactionPoolStatus"
	RAWTRANSACTION          = "/api/blockChain/rawTransaction"
	SENDTRANSACTION         = "/api/blockChain/sendTransaction"
	SENDRAWTRANSACTION      = "/api/blockChain/sendRawTransaction"
	TASKQUEUESTATUS         = "/api/blockChain/taskQueueStatus"
	TRANSACTIONRESULT       = "/api/blockChain/transactionResult"
	TRANSACTIONRESULTS      = "/api/blockChain/transactionResults"
	MBYTRANSACTIONID        = "/api/blockChain/merklePathByTransactionId"
	ADDPEER                 = "/api/net/peer"
	REMOVEPEER              = "/api/net/peer"
	PEERS                   = "/api/net/peers"
	NETWORKINFO             = "/api/net/networkInfo"
	SENDTRANSACTIONS        = "/api/blockChain/sendTransactions"
	EXECUTETRANSACTION      = "/api/blockChain/executeTransaction"
	EXECUTERAWTRANSACTION   = "/api/blockChain/executeRawTransaction"
	FILEDESCRIPTOR          = "/api/blockChain/contractFileDescriptorSet"
	CALCULATETRANSACTIONFEE = "/api/blockChain/calculateTransactionFee"

	privateKeyForView = "680afd630d82ae5c97942c4141d60b8a9fedfa5b2864fca84072c17ee1f72d9d"
	addressForView    = "SD6BXDrKT2syNd1WehtPyRo3dPBiXqfGUj8UJym7YP9W9RynM"
)

type Config struct {
	Endpoints  []string
	Version    string
	PrivateKey string
	Timeout    time.Duration
}

func New(cfg *Config) (*AElfClient, error) {
	c := &AElfClient{
		Host:       cfg.Endpoints[0],
		PrivateKey: cfg.PrivateKey,
		HttpClient: utils.NewHttpClient(cfg.Version, cfg.Timeout),
	}

	SetAElfClientAddress(cfg, c)
	SetAElfClientEndpoint(cfg, c)
	return c, nil
}

func SetAElfClientAddress(cfg *Config, client *AElfClient) {
	if cfg.PrivateKey != "" {
		client.Address = utils.GetAddressFromPrivateKey(cfg.PrivateKey)
	}
}
func SetAElfClientEndpoint(cfg *Config, client *AElfClient) {
}

// IsConnected Verify whether this sdk successfully connects the chain.
func (c *AElfClient) IsConnected() bool {
	data, err := c.GetChainStatus()
	if err != nil || data == nil {
		return false
	}
	return true
}

// GenerateKeyPairInfo Generate KeyPair Info.
func (c *AElfClient) GenerateKeyPairInfo() *types.KeyPair {
	publicKeyBytes, privateKeyBytes := secp256.GenerateKeyPair()
	publicKey := hex.EncodeToString(publicKeyBytes)
	privateKey := hex.EncodeToString(privateKeyBytes)
	privateKeyAddress := c.GetAddressFromPrivateKey(privateKey)
	var keyPair = &types.KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    privateKeyAddress,
	}
	return keyPair
}
