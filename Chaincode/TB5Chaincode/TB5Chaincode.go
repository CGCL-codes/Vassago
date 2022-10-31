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

// Transaction describes basic details of what makes up a simple asset
type Transaction struct {
	ID   string `json:"ID"`
	Type string `json:"Type"`
	UPC  string `json:"UPC"`
	From string `json:"From"`
	To   string `json:"To"`
}
type Rcro struct {
	ID         string `json:"ID"`
	Downstream string `json:"Downstream"`
	Upstream   string `json:"Upstream"`
	UPC        string `json:"UPC"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	TX := []Transaction{
		{ID: "Tx78", UPC: "B1613", Type: "Banana", From: "Bob", To: "Carol"},
		{ID: "Tx79", UPC: "B1613", Type: "Banana", From: "Carol", To: "Dave"},
		{ID: "Tx80", UPC: "B1613", Type: "Banana", From: "Dave", To: "Eva"},
		{ID: "Tx81", UPC: "B1612", Type: "Banana", From: "Bob", To: "Carol"},
		{ID: "Tx82", UPC: "B1612", Type: "Banana", From: "Carol", To: "Dave"},
		{ID: "Tx83", UPC: "B1612", Type: "Banana", From: "Dave", To: "Eva"},
		{ID: "Tx84", UPC: "B1614", Type: "Banana", From: "Bob", To: "Carol"},
		{ID: "Tx85", UPC: "B1614", Type: "Banana", From: "Carol", To: "Dave"},
		{ID: "Tx86", UPC: "B1614", Type: "Banana", From: "Dave", To: "Eva"},
		{ID: "Tx87", UPC: "C1711", Type: "Car", From: "Bob", To: "Carol"},
		{ID: "Tx88", UPC: "C1711", Type: "Car", From: "Carol", To: "Dave"},
		{ID: "Tx89", UPC: "C1711", Type: "Car", From: "Dave", To: "Eva"},
		{ID: "Tx90", UPC: "C1712", Type: "Car", From: "Bob", To: "Carol"},
		{ID: "Tx91", UPC: "C1712", Type: "Car", From: "Carol", To: "Dave"},
		{ID: "Tx92", UPC: "C1712", Type: "Car", From: "Dave", To: "Eva"},
		{ID: "Tx93", UPC: "C1713", Type: "Car", From: "Bob", To: "Carol"},
		{ID: "Tx94", UPC: "C1713", Type: "Car", From: "Carol", To: "Dave"},
		{ID: "Tx95", UPC: "C1713", Type: "Car", From: "Dave", To: "Eva"},
		{ID: "Tx96", UPC: "O1613", Type: "Orange", From: "Bob", To: "Carol"},
		{ID: "Tx97", UPC: "O1613", Type: "Orange", From: "Carol", To: "Dave"},
		{ID: "Tx98", UPC: "O1613", Type: "Orange", From: "Dave", To: "Eva"},
	}
	rcro := []Rcro{
		{ID: "Rcro15", Downstream: "Tx78", Upstream: "Tx59", UPC: "B1613"},
		{ID: "Rcro16", Downstream: "Tx81", Upstream: "Tx62", UPC: "B1612"},
		{ID: "Rcro17", Downstream: "Tx84", Upstream: "Tx65", UPC: "B1614"},
		{ID: "Rcro18", Downstream: "Tx87", Upstream: "Tx68", UPC: "C1711"},
		{ID: "Rcro19", Downstream: "Tx90", Upstream: "Tx71", UPC: "C1712"},
		{ID: "Rcro20", Downstream: "Tx93", Upstream: "Tx74", UPC: "C1713"},
		{ID: "Rcro21", Downstream: "Tx96", Upstream: "T77", UPC: "O1613"},
	}

	for _, Tx := range TX {
		TxJSON, err := json.Marshal(Tx)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(Tx.ID, TxJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	for _, Tx := range rcro {
		TxJSON, err := json.Marshal(Tx)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(Tx.ID, TxJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, ID string, Type string, UPC string, From string, To string) error {
	exists, err := s.AssetExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}

	Tx := Transaction{
		ID:   ID,
		Type: Type,
		UPC:  UPC,
		From: From,
		To:   To,
	}
	TxJSON, err := json.Marshal(Tx)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, TxJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Transaction, error) {
	TxJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if TxJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var Tx Transaction
	err = json.Unmarshal(TxJSON, &Tx)
	if err != nil {
		return nil, err
	}

	return &Tx, nil
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) CreateRcro(ctx contractapi.TransactionContextInterface, ID string, Downsteam string, Upstream string, UPC string) error {
	exists, err := s.AssetExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}

	rcro := Rcro{
		ID:         ID,
		Downstream: Downsteam,
		Upstream:   Upstream,
		UPC:        UPC,
	}
	rcroJSON, err := json.Marshal(rcro)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, rcroJSON)
}

func (s *SmartContract) GetAllRcros(ctx contractapi.TransactionContextInterface) ([]*Rcro, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("Rcro0", "Rcro99")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Rcro
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Rcro
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Transaction, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("Tx0", "Tx999")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Transaction
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) IntraChainTransfer(ctx contractapi.TransactionContextInterface, UPC string) ([]*Transaction, error) {
	Txs, err := s.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	var res []*Transaction
	for _, tx := range Txs {
		if tx.UPC == UPC {
			res = append(res, tx)
		}
	}

	return res, nil
}

func (s *SmartContract) DependQuery(ctx contractapi.TransactionContextInterface, Type string, UPC string) string {
	rcros, err := s.GetAllRcros(ctx)
	if err != nil {
		return err.Error()
	}
	result := "不存在跨链依赖"
	for _, rcro := range rcros {
		if rcro.UPC == UPC {
			result = "存在跨链依赖"
		}
	}
	if result == "存在跨链依赖" {
		result = result + ",存在于" + s.CrossChainTransfer(ctx, Type)
	}
	return result
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (s *SmartContract) CrossChainTransfer(ctx contractapi.TransactionContextInterface, Type string) string {
	ChaincodeName := "dbchaincode"
	ChannelName := "dbchain"
	function := "Query"
	queryArgs := toChaincodeArgs(function, Type)
	response := ctx.GetStub().InvokeChaincode(ChaincodeName, queryArgs, ChannelName)
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		return errStr
	}
	return string(response.Payload)
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
