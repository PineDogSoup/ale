package client

import (
	"ale/utils"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	AElf AElfClient
	AElfCmd
}

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

type Config struct {
	Endpoints  []string
	Version    string
	PrivateKey string
	Timeout    time.Duration
}

func New(cfg *Config) (*Client, error) {
	c := &Client{
		AElf: AElfClient{
			Host:       cfg.Endpoints[0],
			PrivateKey: cfg.PrivateKey,
			HttpClient: utils.NewHttpClient(cfg.Version, cfg.Timeout),
		},
	}

	SetAElfClientAddress(cfg, c)
	SetAElfClientEndpoint(cfg, c)
	return c, nil
}

func SetAElfClientAddress(cfg *Config, client *Client) {
	if cfg.PrivateKey != "" {
		client.AElf.Address = utils.GetAddressFromPrivateKey(cfg.PrivateKey)
	}
}
func SetAElfClientEndpoint(cfg *Config, client *Client) {
}
