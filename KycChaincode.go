package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Region Chaincode implementation
type KycChaincode struct {
}

type KycData struct {
	USER_NAME           string `json:"USER_NAME"`
	USER_ID             string `json:"USER_ID"`
	KYC_BANK_NAME       string `json:"KYC_BANK_NAME"`
	KYC_CREATE_DATE     string `json:"KYC_CREATE_DATE"`
	KYC_VALID_TILL_DATE string `json:"KYC_VALID_TILL_DATE"`
	KYC_DOC_BLOB        string `json:"KYC_DOC_BLOB"`
}

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("tblKycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "KYC_BANK_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_CREATE_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_VALID_TILL_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_DOC_BLOB", Type: shim.ColumnDefinition_STRING, Key: false},
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
	} else if function == "UpdateKycDetails" {
		// Update User's KYC data in blockchain
		return t.UpdateKycDetails(stub, args)
	}

	return nil, nil
}

func (t *KycChaincode) InsertKycDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 4 arguments")
	}

	UserName := args[0]
	UserId := args[1]
	BankName := args[2]
	KycDoc := args[3]
	CurrentDate := time.Now().Local()
	CreateDate := CurrentDate.Format("02-01-2006")
	ValidTillDate := CurrentDate.AddDate(2, 0, 0).Format("02-01-2006")

	ok, err := stub.InsertRow("tblKycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: UserId}},
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_String_{String_: CreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ValidTillDate}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDoc}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding KYC record.")
	}
	return nil, nil
}

func (t *KycChaincode) UpdateKycDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	UserName := args[0]
	UserId := args[1]
	BankName := args[2]
	KycDoc := args[3]
	CurrentDate := time.Now().Local()
	CreateDate := CurrentDate.Format("02-01-2006")
	ValidTillDate := CurrentDate.AddDate(2, 0, 0).Format("02-01-2006")

	ok, err := stub.ReplaceRow("tblKycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: UserId}},
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_String_{String_: CreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ValidTillDate}},
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
	var jsonAsBytes []byte
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting enrollId to query")
	}

	UserId := args[0]

	if UserId == "" {
		table, err := stub.GetTable("tblKycDetails")
		if err != nil {
			return nil, errors.New("Failed to query")
		}
		output := table.String()
		jsonAsBytes, _ = json.Marshal(output)
	} else {
		var columns []shim.Column

		col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
		columns = append(columns, col1)

		row, err := stub.GetRow("tblKycDetails", columns)
		if err != nil {
			return nil, errors.New("Failed to query")
		}

		KycDataObj.USER_NAME = row.Columns[0].GetString_()
		KycDataObj.USER_ID = row.Columns[1].GetString_()
		KycDataObj.KYC_BANK_NAME = row.Columns[2].GetString_()
		KycDataObj.KYC_CREATE_DATE = row.Columns[3].GetString_()
		KycDataObj.KYC_VALID_TILL_DATE = row.Columns[4].GetString_()
		KycDataObj.KYC_DOC_BLOB = row.Columns[5].GetString_()

		jsonAsBytes, _ = json.Marshal(KycDataObj)
	}

	return jsonAsBytes, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
