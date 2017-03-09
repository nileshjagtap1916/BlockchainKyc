package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func InitializeChaincode(stub shim.ChaincodeStubInterface) error {
	return CreateDatabase(stub)
}

func SaveKycDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//return InsertKYCDetails(stub, args)
	var KycDetails KycData
	var err error
	var ok bool

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 4 argument")
	}

	//get data from middle layer
	KycDetails.USER_ID = args[0]
	KycDetails.KYC_BANK_NAME = args[1]
	KycDetails.USER_NAME = args[2]
	KycDetails.KYC_DOC_BLOB = args[3]
	CurrentDate := time.Now().Local()
	KycDetails.KYC_CREATE_DATE = CurrentDate.Format("02-01-2006")
	KycDetails.KYC_VALID_TILL_DATE = CurrentDate.AddDate(2, 0, 0).Format("02-01-2006")

	//save data into blockchain
	ok, err = InsertKYCDetails(stub, KycDetails)
	if !ok && err == nil {
		return nil, errors.New("Error in adding KycDetails record.")
	}

	// Update Userlist with current UserId
	UserList, _ := GetUserList(stub, KycDetails.KYC_BANK_NAME)
	UserList = append(UserList, KycDetails.USER_ID)

	//Update Bank details on blockchain
	ok, err = UpdateBankDetails(stub, args[1], UserList)
	if !ok && err == nil {
		return nil, errors.New("Error in Updating User ContractList")
	}

	return nil, nil
}

func SaveBankDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var UserList []string
	var err error
	var ok bool

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	BankName := args[0]

	//save data into blockchain
	ok, err = InsertBankDetails(stub, BankName, UserList)
	if !ok && err == nil {
		return nil, errors.New("Error in adding BankDetails record.")
	}

	return nil, nil
}

func GetKycByBankName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//return GetKYCDetails(stub, args)

	var KycList []KycData
	var KycDetails KycData

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	BankName := args[0]
	//UserId := args[1]

	//get data from blockchain
	UserList, _ := GetUserList(stub, BankName)

	for _, UserId := range UserList {
		KycDetails, _ = GetKYCDetails(stub, UserId, BankName)
		KycList = append(KycList, KycDetails)
	}

	JsonAsBytes, _ := json.Marshal(KycList)

	return JsonAsBytes, nil
}

func GetKycCount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	All := 0
	Expering := 0
	Created := 0
	var KycDetails KycData
	var KycCountObj KycCount

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	BankName := args[0]

	UserList, _ := GetUserList(stub, BankName)

	for _, UserId := range UserList {
		All = All + 1
		KycDetails, _ = GetKYCDetails(stub, UserId, BankName)

		CurrentMonth := time.Now().Month()
		ValidTillDate, _ := time.Parse("02-01-2006", KycDetails.KYC_VALID_TILL_DATE)
		CreateDate, _ := time.Parse("02-01-2006", KycDetails.KYC_CREATE_DATE)

		if CurrentMonth == ValidTillDate.Month() {
			Expering = Expering + 1
		}
		if CurrentMonth == CreateDate.Month() {
			Created = Created + 1
		}
	}

	KycCountObj.AllContracts = All
	KycCountObj.ExperingContracts = Expering
	KycCountObj.CreatedContracts = Created

	JsonAsBytes, _ := json.Marshal(KycCountObj)

	return JsonAsBytes, nil
}
