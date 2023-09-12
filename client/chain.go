package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/json"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

// GetChainStatus Get the current status of the block chain.
func (c *AElfClient) GetChainStatus() (*types.ChainStatus, error) {
	url := c.Host + CHAINSTATUS
	chainBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get ChainStatus error:" + err.Error())
	}
	var chain = new(types.ChainStatus)
	json.Unmarshal(chainBytes, &chain)
	return chain, nil
}

// GetContractFileDescriptorSet Get the definitions of proto-buff related to a contract.
func (c *AElfClient) GetContractFileDescriptorSet(address string) ([]byte, error) {
	url := c.Host + FILEDESCRIPTOR
	params := map[string]interface{}{"address": address}
	data, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get ContractFile Descriptor Set error:" + err.Error())
	}
	return data, err
}

// GetChainID Get id of the chain.
func (c *AElfClient) GetChainID() (int, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return 0, errors.New("Get Chain Status error:" + err.Error())
	}
	chainIDBytes := base58.Decode(chainStatus.ChainId)
	if len(chainIDBytes) < 4 {
		var bs [4]byte
		for i := 0; i < 4; i++ {
			bs[i] = 0
			if len(chainIDBytes) > i {
				bs[i] = chainIDBytes[i]
			}
		}
		chainIDBytes = bs[:]
	}
	return utils.BytesToInt(chainIDBytes), nil
}

// GetTaskQueueStatus Get the status information of the task queue.
func (c *AElfClient) GetTaskQueueStatus() ([]*types.TaskQueueInfo, error) {
	url := c.Host + TASKQUEUESTATUS
	queues, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get Task Queue Status error:" + err.Error())
	}
	var datas interface{}
	json.Unmarshal(queues, &datas)
	var queueInfos []*types.TaskQueueInfo
	for _, data := range datas.([]interface{}) {
		var queue = new(types.TaskQueueInfo)
		queueBytes, err := json.Marshal(data)
		if err != nil {
			return nil, errors.New("json Marshal data error:" + err.Error())
		}
		json.Unmarshal(queueBytes, &queue)
		queueInfos = append(queueInfos, queue)
	}
	return queueInfos, nil
}
