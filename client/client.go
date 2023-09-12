package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/hex"
	"errors"
	secp256 "github.com/haltingstate/secp256k1-go"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// AElfClient AElf Client.
type AElfClient struct {
	Host       string
	Version    string
	PrivateKey string
	HttpClient *utils.HttpClient
	endpoints  []url.URL
	rand       *rand.Rand
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

	ExamplePrivateKey = "680afd630d82ae5c97942c4141d60b8a9fedfa5b2864fca84072c17ee1f72d9d"
)

var (
	ErrNoEndpoints = errors.New("client: no endpoints available")
)

type Config struct {
	Endpoints  []string
	Version    string
	PrivateKey string
	Timeout    time.Duration
}

func New(cfg *Config) (*AElfClient, error) {
	c := &AElfClient{
		PrivateKey: cfg.PrivateKey,
		HttpClient: utils.NewHttpClient(cfg.Version, cfg.Timeout),
	}
	if err := c.SetEndpoints(cfg.Endpoints); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AElfClient) SetEndpoints(eps []string) error {
	neps, err := parseEndpoints(eps)
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.endpoints = shuffleEndpoints(c.rand, neps)
	return nil
}

func parseEndpoints(eps []string) ([]url.URL, error) {
	if len(eps) == 0 {
		return []url.URL{}, ErrNoEndpoints
	}

	neps := make([]url.URL, len(eps))
	for i, ep := range eps {
		u, err := url.Parse(ep)
		if err != nil {
			return []url.URL{}, err
		}
		neps[i] = *u
	}
	return neps, nil
}

func shuffleEndpoints(r *rand.Rand, eps []url.URL) []url.URL {
	// copied from Go 1.9<= rand.Rand.Perm
	n := len(eps)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		j := r.Intn(i + 1)
		p[i] = p[j]
		p[j] = i
	}
	neps := make([]url.URL, n)
	for i, k := range p {
		neps[i] = eps[k]
	}
	return neps
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
