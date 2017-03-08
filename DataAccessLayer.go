package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CreateDatabase(stub shim.ChaincodeStubInterface) error {
	var err error

	//Create table "tblKycDetails"
	err = stub.CreateTable("KycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UserId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "BankName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UserName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CreateDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ValidTillDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KycDocument", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating KycDetails table.")
	}

	//Create table "BankDetails"
	err = stub.CreateTable("BankDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "BankName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UserList", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating BankDetails table.")
	}

	return nil
}

func InsertKYCDetails(stub shim.ChaincodeStubInterface, KYCDetails KycData) (bool, error) {
	var err error
	var ok bool
	ok, err = stub.InsertRow("KycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.UserId}},
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.BankName}},
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.CreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.ValidTillDate}},
			&shim.Column{Value: &shim.Column_String_{String_: KYCDetails.KycDocument}},
		},
	})
	return ok, err
}

func InsertBankDetails(stub shim.ChaincodeStubInterface, BankName string, UserList []string) (bool, error) {
	var err error
	var ok bool
	JsonAsBytes, _ := json.Marshal(UserList)
	ok, err = stub.InsertRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})
	return ok, err
}

func GetKYCDetails(stub shim.ChaincodeStubInterface, UserId string, BankName string) (KycData, error) {
	var KYCDetails KycData
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: BankName}}
	columns = append(columns, col2)

	row, err := stub.GetRow("KYCDetails", columns)
	if err != nil {
		return KYCDetails, errors.New("Failed to query table ContractDetails")
	}

	KYCDetails.UserId = row.Columns[0].GetString_()
	KYCDetails.BankName = row.Columns[1].GetString_()
	KYCDetails.UserName = row.Columns[2].GetString_()
	KYCDetails.CreateDate = row.Columns[3].GetString_()
	KYCDetails.ValidTillDate = row.Columns[4].GetString_()
	KYCDetails.KycDocument = row.Columns[5].GetString_()

	return KYCDetails, nil
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

	ok, err := stub.ReplaceRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})

	if !ok && err == nil {
		return false, errors.New("Error in updating Bank record.")
	}
	return true, nil
}
