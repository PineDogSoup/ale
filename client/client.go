package client

import (
	"ale/utils"
	"context"
	"net/url"
	"time"
)

type Client struct {
	ctx  context.Context
	AElf AElfClient
}

type AElfClient struct {
	Host           string
	Version        string
	PrivateKey     string
	Address        string
	HttpClient     *utils.HttpClient
	endpoints      []url.URL
	PortkeyAddress string
}

type Config struct {
	ChainId    string
	Env        string
	Endpoints  []string
	Version    string
	PrivateKey string
	Timeout    time.Duration
}

func New(cfg *Config) (*Client, error) {
	c := &Client{
		AElf: AElfClient{
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
	if len(cfg.Endpoints) > 0 {
		client.AElf.Host = cfg.Endpoints[0]
	}
}
