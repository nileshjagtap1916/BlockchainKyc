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
	KycDetails.KYC_CREATE_DATE = CurrentDate.Format("02 Jan 2006")
	KycDetails.KYC_VALID_TILL_DATE = CurrentDate.AddDate(2, 0, 0).Format("02 Jan 2006")

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

func GetKycByUserId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var KycDetails KycData

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	UserId := args[0]
	KycDetails, _ = GetKYCDetails(stub, UserId)
	JsonAsBytes, _ := json.Marshal(KycDetails)

	return JsonAsBytes, nil
}

func GetKycByBankName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
		KycDetails, _ = GetKYCDetails(stub, UserId)
		KycList = append(KycList, KycDetails)
	}

	JsonAsBytes, _ := json.Marshal(KycList)

	return JsonAsBytes, nil
}

func GetKycByExpiringMonth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var KycList []KycData
	var KycDetails KycData

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	BankName := args[0]

	//get data from blockchain
	UserList, _ := GetUserList(stub, BankName)

	for _, UserId := range UserList {
		KycDetails, _ = GetKYCDetails(stub, UserId)

		CurrentDate := time.Now()
		ValidTillDate, _ := time.Parse("02 Jan 2006", KycDetails.KYC_VALID_TILL_DATE)

		if CurrentDate.Month() == ValidTillDate.Month() && CurrentDate.Year() == ValidTillDate.Year() {
			KycList = append(KycList, KycDetails)
		}
	}

	JsonAsBytes, _ := json.Marshal(KycList)

	return JsonAsBytes, nil
}

func GetKycByCreatedMonth(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var KycList []KycData
	var KycDetails KycData

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	BankName := args[0]

	//get data from blockchain
	UserList, _ := GetUserList(stub, BankName)

	for _, UserId := range UserList {
		KycDetails, _ = GetKYCDetails(stub, UserId)

		CurrentDate := time.Now()
		CreateDate, _ := time.Parse("02 Jan 2006", KycDetails.KYC_CREATE_DATE)
		if CurrentDate.Month() == CreateDate.Month() && CurrentDate.Year() == CreateDate.Year() {
			KycList = append(KycList, KycDetails)
		}
	}

	JsonAsBytes, _ := json.Marshal(KycList)

	return JsonAsBytes, nil
}

func GetKycCount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	All := 0
	Expiring := 0
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
		KycDetails, _ = GetKYCDetails(stub, UserId)

		CurrentDate := time.Now()
		ValidTillDate, _ := time.Parse("02 Jan 2006", KycDetails.KYC_VALID_TILL_DATE)
		CreateDate, _ := time.Parse("02 Jan 2006", KycDetails.KYC_CREATE_DATE)

		if CurrentDate.Month() == ValidTillDate.Month() && CurrentDate.Year() == ValidTillDate.Year() {
			Expiring = Expiring + 1
		}
		if CurrentDate.Month() == CreateDate.Month() && CurrentDate.Year() == CreateDate.Year() {
			Created = Created + 1
		}
	}

	KycCountObj.AllContracts = All
	KycCountObj.ExpiringContracts = Expiring
	KycCountObj.CreatedContracts = Created

	JsonAsBytes, _ := json.Marshal(KycCountObj)

	return JsonAsBytes, nil
}

func UpdateKyc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var ok bool

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	//get data from middle layer
	KycDetails, _ := GetKYCDetails(stub, args[0])

	KycDetails.USER_NAME = args[1]
	KycDetails.KYC_DOC_BLOB = args[2]
	CurrentDate := time.Now().Local()
	KycDetails.KYC_VALID_TILL_DATE = CurrentDate.AddDate(2, 0, 0).Format("02 Jan 2006")

	//Update data into blockchain
	ok, err = UpdateKycDetails(stub, KycDetails)
	if !ok && err == nil {
		return nil, errors.New("Error in updating KycDetails record.")
	}

	return nil, nil
}
