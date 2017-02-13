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

type KycData struct {
	USER_PAN_NO  string `json:"USER_PAN_NO"`
	USER_NAME    string `json:"USER_NAME"`
	USER_KYC_PDF string `json:"USER_KYC_PDF"`
}

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("tblKycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_PAN_NO", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "USER_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "USER_KYC_PDF", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating KYC table.")
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

	return nil, nil
}

func (t *KycChaincode) InsertKycDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	UserPanNumber := args[0]
	UserName := args[1]
	KycDoc := args[2]
	ok, err := stub.InsertRow("tblKycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UserPanNumber}},
			&shim.Column{Value: &shim.Column_String_{String_: UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDoc}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding KYC record.")
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var KycDataObj KycData

	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting enrollId to query")
	}

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
	columns = append(columns, col1)

	row, err := stub.GetRow("tblKycDetails", columns)
	if err != nil {
		return nil, errors.New("Failed to query")
	}

	KycDataObj.USER_PAN_NO = row.Columns[0].GetString_()
	KycDataObj.USER_NAME = row.Columns[1].GetString_()
	KycDataObj.USER_KYC_PDF = row.Columns[2].GetString_()

	jsonAsBytes, _ := json.Marshal(KycDataObj)

	return jsonAsBytes, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
