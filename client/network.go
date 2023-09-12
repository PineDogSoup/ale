package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/json"
	"errors"
)

// GetNetworkInfo Get the node's network information.
func (c *AElfClient) GetNetworkInfo() (*types.NetworkInfo, error) {
	url := c.Host + NETWORKINFO
	networkBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get Network Info error:" + err.Error())
	}
	var network = new(types.NetworkInfo)
	json.Unmarshal(networkBytes, &network)
	return network, nil
}

// RemovePeer Attempt to remove a node from the connected network nodes by given the ipAddress.
//func (ac *AElfClient) RemovePeer(ipAddress string) (bool, error) {
//	url := ac.Host + REMOVEPEER
//	combine := ac.UserName + ":" + ac.Password
//	combineToBase64 := "Basic " + base64.StdEncoding.EncodeToString([]byte(combine))
//	params := map[string]interface{}{"address": ipAddress}
//	peerBytes, err := utils.GetRequestWithAuth("DELETE", url, ac.Version, params, combineToBase64)
//	if err != nil {
//		return false, errors.New("Remove Peer error:" + err.Error())
//	}
//	var data interface{}
//	json.Unmarshal(peerBytes, &data)
//	return data.(bool), nil
//}
//
// AddPeer Attempt to add a node to the connected network nodes.Input parameter contains the ipAddress of the node.
//func (ac *AElfClient) AddPeer(ipAddress string) (bool, error) {
//	url := ac.Host + ADDPEER
//	combine := ac.UserName + ":" + ac.Password
//	combineToBase64 := "Basic " + base64.StdEncoding.EncodeToString([]byte(combine))
//	params := map[string]interface{}{"Address": ipAddress}
//	peerBytes, err := utils.PostRequestWithAuth(url, ac.Version, params, combineToBase64)
//	if err != nil {
//		return false, errors.New("Add Peer error:" + err.Error())
//	}
//	var data interface{}
//	json.Unmarshal(peerBytes, &data)
//	return data.(bool), nil
//}

// GetPeers Gets information about the peer nodes of the current node.Optional whether to include metrics.
func (c *AElfClient) GetPeers(withMetrics bool) ([]*types.Peer, error) {
	url := c.Host + PEERS
	params := map[string]interface{}{"withMetrics": withMetrics}
	peerBytes, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get Peers error:" + err.Error())
	}
	var datas interface{}
	var peers []*types.Peer
	json.Unmarshal(peerBytes, &datas)
	for _, data := range datas.([]interface{}) {
		var peer = new(types.Peer)
		peerBytes, _ := json.Marshal(data)
		json.Unmarshal(peerBytes, &peer)
		peers = append(peers, peer)
	}
	return peers, nil
}
