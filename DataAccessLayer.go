package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CreateDatabase(stub shim.ChaincodeStubInterface) error {
	var err error

	//Create table "tblKycDetails"
	err = stub.CreateTable("tblKycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "KYC_BANK_NAME", Type: shim.ColumnDefinition_STRING, Key: true},
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

/*func InsertKYCDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 4 arguments")
	}

	UserId := args[0]
	BankName := args[1]
	UserName := args[2]
	KycDoc := args[3]
	CurrentDate := time.Now().Local()
	CreateDate := CurrentDate.Format("02-01-2006")
	ValidTillDate := CurrentDate.AddDate(2, 0, 0).Format("02-01-2006")

	ok, err := stub.InsertRow("tblKycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UserId}},
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_String_{String_: UserName}},
			&shim.Column{Value: &shim.Column_String_{String_: CreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ValidTillDate}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDoc}},
		},
	})

	if !ok && err == nil {
		return nil, errors.New("Error in adding KYC record.")
	}
	return nil, nil
}*/

func InsertKYCDetails(stub shim.ChaincodeStubInterface, Kycdetails KycData) (bool, error) {
	return stub.InsertRow("tblKycDetails", shim.Row{
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

func GetKYCDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var KycDataObj KycData
	//var jsonAsBytes []byte
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting UserId and UserBank to query")
	}

	UserId := args[0]
	UserBank := args[1]

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: UserBank}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("tblKycDetails", columns)
	if err != nil {
		return nil, errors.New("Failed to query")
	}

	KycDataObj.USER_ID = row.Columns[0].GetString_()
	KycDataObj.KYC_BANK_NAME = row.Columns[1].GetString_()
	KycDataObj.USER_NAME = row.Columns[2].GetString_()
	KycDataObj.KYC_CREATE_DATE = row.Columns[3].GetString_()
	KycDataObj.KYC_VALID_TILL_DATE = row.Columns[4].GetString_()
	KycDataObj.KYC_DOC_BLOB = row.Columns[5].GetString_()

	jsonAsBytes, _ := json.Marshal(KycDataObj)
	return jsonAsBytes, nil
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
