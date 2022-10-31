package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

type SmartContract struct {
	contractapi.Contract
}

type Transfer struct {
	ID    string `json:"ID"`
	Type  string `json:"Type"`
	Chain string `json:"Chain"`
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	Trans := []Transfer{
		{ID: "Trans1", Type: "Milk", Chain: "tb1chain"},
		{ID: "Trans2", Type: "Milk", Chain: "tb2chain"},
		{ID: "Trans3", Type: "Milk", Chain: "tb3chain"},
		{ID: "Trans4", Type: "Cake", Chain: "tb1chain"},
		{ID: "Trans5", Type: "Cake", Chain: "tb2chain"},
		{ID: "Trans6", Type: "Cake", Chain: "tb3chain"},
		{ID: "Trans7", Type: "Apple", Chain: "tb1chain"},
		{ID: "Trans8", Type: "Apple", Chain: "tb2chain"},
		{ID: "Trans9", Type: "Apple", Chain: "tb3chain"},
		{ID: "Trans10", Type: "Banana", Chain: "tb4chain"},
		{ID: "Trans11", Type: "Banana", Chain: "tb5chain"},
		{ID: "Trans12", Type: "Banana", Chain: "tb6chain"},
		{ID: "Trans13", Type: "Car", Chain: "tb4chain"},
		{ID: "Trans14", Type: "Car", Chain: "tb5chain"},
		{ID: "Trans16", Type: "Car", Chain: "tb6chain"},
		{ID: "Trans16", Type: "Orange", Chain: "tb4chain"},
		{ID: "Trans17", Type: "Orange", Chain: "tb5chain"},
		{ID: "Trans18", Type: "Orange", Chain: "tb6chain"},
	}
	for _, asset := range Trans {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CreateTrans(ctx contractapi.TransactionContextInterface, ID string, Type string, Chain string) error {
	exists, err := s.AssetExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}

	Trans := Transfer{
		ID:    ID,
		Type:  Type,
		Chain: Chain,
	}
	rcroJSON, err := json.Marshal(Trans)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, rcroJSON)
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
func (s *SmartContract) GetAllTrans(ctx contractapi.TransactionContextInterface) ([]*Transfer, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("Trans0", "Trans999")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Transfer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Transfer
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, Type string) string {
	Trans, err := s.GetAllTrans(ctx)
	if err != nil {
		return err.Error()
	}
	var result string
	for _, trans := range Trans {
		if trans.Type == Type {
			result = result + trans.Chain + " "
		}
	}
	return result
}

func (s *SmartContract) WODQuery(ctx contractapi.TransactionContextInterface, UPC string) string {
	ChannelName := []string{"tb1chain", "tb2chain", "tb3chain", "tb4chain", "tb5chain", "tb6chain"}
	var result string
	for j := 0; j < len(ChannelName); j++ {
		ChaincodeName := ChannelName[j] + "code"
		function := "IntraChainTransfer"
		queryArgs := toChaincodeArgs(function, UPC)
		response := ctx.GetStub().InvokeChaincode(ChaincodeName, queryArgs, ChannelName[j])
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
			return errStr
		}
		result = string(response.Payload) + "," + result
	}
	return result
}

func (s *SmartContract) WDQuery(ctx contractapi.TransactionContextInterface, Type string, UPC string) string {
	Trans, err := s.GetAllTrans(ctx)
	if err != nil {
		return err.Error()
	}
	ChannelName := make([]string, 0)
	for _, trans := range Trans {
		if trans.Type == Type {
			ChannelName = append(ChannelName, trans.Chain)
		}
	}
	var result string
	for j := 0; j < len(ChannelName); j++ {
		ChaincodeName := ChannelName[j] + "code"
		function := "IntraChainTransfer"
		queryArgs := toChaincodeArgs(function, UPC)
		response := ctx.GetStub().InvokeChaincode(ChaincodeName, queryArgs, ChannelName[j])
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
			return errStr
		}
		result = string(response.Payload) + "," + result
	}
	return result
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
