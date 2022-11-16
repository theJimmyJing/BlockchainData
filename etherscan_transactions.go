package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var Etherscan_API_KEY = "BN94Y44F1PU9U7X2ZMVXJFYF4TUJ5AAKRE" // dk's API KEY for DEBUG

// 从http://api.etherscan.io/api 获取指定地址的交易记录
func EtherscanAddrTransactions(address string) EtherscanTransactionsResp {
	// 查询指定地址的所有交易记录
	var data = EtherscanTransactionsResp{}

	url := "http://api.etherscan.io/api?module=account&action=txlist&address=" + address + "&startblock=0&endblock=99999999&sort=asc&apikey=" + Etherscan_API_KEY
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("EtherscanAddrTransactions -1 err: ", err)
		data.Status = "-1"
		data.Message = err.Error()
		return data
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("EtherscanAddrTransactions -2 err: ", err)
		// handle err
		data.Status = "-2"
		data.Message = err.Error()
		return data
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("EtherscanAddrTransactions -3 err: ", err)
			data.Status = "-3"
			data.Message = err.Error()
			return data
		}
	}

	rbytes, _ := json.Marshal(data)
	fmt.Println("EtherscanAddrTransactions : ", string(rbytes))

	return data
}

type EtherscanTransactionsResp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []EtherscanTransaction `json:"result"`
}
type EtherscanTransaction struct {
	BlockNumber       string `json:"blockNumber,omitempty"`
	TimeStamp         string `json:"timeStamp,omitempty"`
	Hash              string `json:"hash,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	BlockHash         string `json:"blockHash,omitempty"`
	TransactionIndex  string `json:"transactionIndex,omitempty"`
	From              string `json:"from,omitempty"`
	To                string `json:"to,omitempty"`
	Value             string `json:"value,omitempty"`
	Gas               string `json:"gas,omitempty"`
	GasPrice          string `json:"gasPrice,omitempty"`
	IsError           string `json:"isError,omitempty"`
	TxreceiptStatus   string `json:"txreceipt_status,omitempty"`
	Input             string `json:"input,omitempty"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	CumulativeGasUsed string `json:"cumulativeGasUsed,omitempty"`
	GasUsed           string `json:"gasUsed,omitempty"`
	Confirmations     string `json:"confirmations,omitempty"`
	MethodID          string `json:"methodId,omitempty"`
	FunctionName      string `json:"functionName,omitempty"`
}
