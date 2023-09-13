package types

import client "ale/protobuf/generated"

type ContractInfo struct {
	Info    *client.ContractInfo
	Address string
}
