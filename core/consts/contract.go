package consts

const (
	// contract
	TokenContractSystemName      = "AElf.ContractNames.Token"
	PortkeyContractSystemName    = "Portkey.ContractNames.CA"
	CrossChainContractSystemName = "AElf.ContractNames.CrossChain"
	ForestContractSystemName     = "AElf.ContractNames.Forest"
	EconomicContractSystemName   = "AElf.ContractNames.Economic"
	ProfitContractSystemName     = "AElf.ContractNames.Profit"
	TreasuryContractSystemName   = "AElf.ContractNames.Treasury"
	ElectionContractSystemName   = "AElf.ContractNames.Election"
	VoteContractSystemName       = "AElf.ContractNames.Vote"
	ConsensusContractSystemName  = "AElf.ContractNames.Consensus"
	GenesisContractSystemName    = "AElf.ContractNames.Genesis"

	//  contract LogEvent
	TransferredLogEventName           = "Transferred"
	CrossChainTransferredLogEventName = "CrossChainTransferred"
	CrossChainReceivedLogEventName    = "CrossChainReceived"

	// contract Method
	GenesisContractGetContractInfo      = "GetContractInfo"
	TokenContractGetBalance             = "GetBalance"
	TokenContractGetTokenInfo           = "GetTokenInfo"
	TokenContractTransfer               = "Transfer"
	TokenContractCrossChainTransfer     = "CrossChainTransfer"
	CrossChainContractCrossChainReceive = "CrossChainReceive"
	PortkeyContractGetHolderInfo        = "GetHolderInfo"
	PortkeyContractManagerApprove       = "ManagerApprove"
)
