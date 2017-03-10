package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CreateDatabase(stub shim.ChaincodeStubInterface) error {
	var err error

	//Create table "KycDetails"
	err = stub.CreateTable("KycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "KYC_BANK_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "USER_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_CREATE_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_VALID_TILL_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_DOC_BLOB", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating KycDetails table.")
	}

	//Create table "BankDetails"
	err = stub.CreateTable("BankDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "BankName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UserList", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating BankDetails table.")
	}

	return nil
}

func InsertKYCDetails(stub shim.ChaincodeStubInterface, Kycdetails KycData) (bool, error) {
	return stub.InsertRow("KycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.USER_ID}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_BANK_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.USER_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_CREATE_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_VALID_TILL_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_DOC_BLOB}},
		},
	})
}

func InsertBankDetails(stub shim.ChaincodeStubInterface, BankName string, UserList []string) (bool, error) {
	JsonAsBytes, _ := json.Marshal(UserList)
	return stub.InsertRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})
}

func GetKYCDetails(stub shim.ChaincodeStubInterface, UserId string) (KycData, error) {
	var KycDataObj KycData

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("KycDetails", columns)
	if err != nil {
		return KycDataObj, errors.New("Failed to query")
	}

	KycDataObj.USER_ID = row.Columns[0].GetString_()
	KycDataObj.KYC_BANK_NAME = row.Columns[1].GetString_()
	KycDataObj.USER_NAME = row.Columns[2].GetString_()
	KycDataObj.KYC_CREATE_DATE = row.Columns[3].GetString_()
	KycDataObj.KYC_VALID_TILL_DATE = row.Columns[4].GetString_()
	KycDataObj.KYC_DOC_BLOB = row.Columns[5].GetString_()

	return KycDataObj, nil
}

func GetUserList(stub shim.ChaincodeStubInterface, BankName string) ([]string, error) {
	var UserList []string
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: BankName}}
	columns = append(columns, col1)

	row, err := stub.GetRow("BankDetails", columns)
	if err != nil {
		return UserList, errors.New("Failed to query table BankDetails")
	}

	UsersAsBytes := row.Columns[1].GetBytes()
	json.Unmarshal(UsersAsBytes, &UserList)

	return UserList, nil
}

func UpdateBankDetails(stub shim.ChaincodeStubInterface, BankName string, Userlist []string) (bool, error) {

	JsonAsBytes, _ := json.Marshal(Userlist)

	return stub.ReplaceRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})
}

func UpdateKycDetails(stub shim.ChaincodeStubInterface, KycDetails KycData) (bool, error) {

	return stub.ReplaceRow("KycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.USER_ID}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_BANK_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.USER_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_CREATE_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_VALID_TILL_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_DOC_BLOB}},
		},
	})
}
