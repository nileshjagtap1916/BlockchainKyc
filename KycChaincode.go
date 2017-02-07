package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Region Chaincode implementation
type KycChaincode struct {
}

var KycIndexTxStr = "_KycIndexTxStr"

/*type KycData struct {
	USER_PAN_NO  string `json:"USER_PAN_NO"`
	USER_NAME    string `json:"USER_NAME"`
	USER_KYC_PDF string `json:"USER_KYC_PDF"`
}*/

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode

	fmt.Printf("Deployment of KYC is completed\n")

	/*var EmptyKYC KycData
	jsonAsBytes, _ := json.Marshal(EmptyKYC)
	err = stub.PutState(KycIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}*/

	// Create ownership table
	err = stub.CreateTable("tblKycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_PAN_NO", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "USER_NAME", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "USER_KYC_PDF", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating tblKycDetails table, [%v]", err)
	}

	return nil, nil
}

// Add user KYC data in Blockchain
func (t *KycChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "InsertKycDetails" {
		// Insert User's KYC data in blockchain
		return t.InsertKycDetails(stub, args)
	}
	/*else if function == "UpdateKycDetails" {
		// Update User's KYC data in blockchain
		return t.UpdateKycDetails(stub, args)
	}*/

	return nil, errors.New("Received unknown function invocation")
}

func (t *KycChaincode) InsertKycDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var ok bool
	var UserPanNumber string
	var UserName string
	var UserKycDoc string

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	// Initialize the chaincode
	UserPanNumber = args[0]
	UserName = args[1]
	UserKycDoc = args[2]

	ok, err = stub.InsertRow("tblKycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UserPanNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: UserKycDoc}}},
	})

	if !ok && err == nil {
		fmt.Println("Error inserting row")
		return nil, errors.New("Kyc Details already on blockchain.")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var UserPanNumber string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	UserPanNumber = args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: UserPanNumber}}
	columns = append(columns, col1)

	row, err := stub.GetRow("tblKycDetails", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed retrieving data for " + UserPanNumber + ". Error " + err.Error() + ". \"}"
		return nil, errors.New(jsonResp)
	}

	/*if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"no data present for " + UserPanNumber + " on blockchain. \"}"
		return nil, errors.New(jsonResp)
	}*/

	jsonResp := "{\"KYC_DOC\":\"" + row.Columns[2].GetString_() + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	res, _ := json.Marshal(row)

	return res, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
