package main

import (
	"encoding/json"
	"errors"

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
	/*KycDetails.UserId = args[0]
	KycDetails.BankName = args[1]
	KycDetails.UserName = args[2]
	KycDetails.KycDocument = args[3]
	CurrentDate := time.Now().Local()
	KycDetails.CreateDate = CurrentDate.Format("02-01-2006")
	KycDetails.ValidTillDate = CurrentDate.AddDate(2, 0, 0).Format("02-01-2006")*/

	//save data into blockchain
	ok, err = InsertKYCDetails(stub, args)
	if !ok && err == nil {
		return nil, errors.New("Error in adding KycDetails record.")
	}

	// Update Userlist with current UserId
	UserList, _ := GetUserList(stub, args[1])
	UserList = append(UserList, KycDetails.UserId)

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

func GetKyc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var KycList []KycData
	var KycDetails KycData

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}

	//get data from middle layer
	BankName := args[0]
	UserId := args[1]

	//get data from blockchain
	//UserList, _ := GetUserList(stub, BankName)

	//for _, UserId := range UserList {
	KycDetails, _ = GetKYCDetails(stub, UserId, BankName)
	//KycList = append(KycList, KycDetails)
	//}

	JsonAsBytes, _ := json.Marshal(KycDetails)

	return JsonAsBytes, nil
}
