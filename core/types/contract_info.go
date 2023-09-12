package types

import client "ale/protobuf/generated"

type ContractInfo struct {
	ContractName string
	Info         *client.ContractInfo
	Address      string
}
