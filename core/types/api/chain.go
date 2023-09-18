package api

type SearchChainInfo struct {
	TotalCount int         `json:"totalCount"`
	Items      []ChainInfo `json:"items"`
}

type ChainInfo struct {
	Id                string `json:"id"`
	ChainId           string `json:"chainId"`
	ChainName         string `json:"chainName"`
	Endpoint          string `json:"endpoint"`
	ExplorerUrl       string `json:"explorerUrl"`
	CaContractAddress string `json:"caContractAddress"`
	LastModifyTime    string `json:"lastModifyTime"`
	DefaultToken      struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		ImageUrl string `json:"imageUrl"`
		Symbol   string `json:"symbol"`
		Decimal  string `json:"decimal"`
	} `json:"defaultToken"`
}
