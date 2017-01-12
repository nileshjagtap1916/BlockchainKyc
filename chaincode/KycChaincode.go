package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
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

	var err error
	// Initialize the chaincode

	fmt.Printf("Deployment of KYC is completed\n")

	var EmptyKYC KycData
	jsonAsBytes, _ := json.Marshal(EmptyKYC)
	err = stub.PutState(KycIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Add user KYC data in Blockchain
func (t *KycChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == KycIndexTxStr {
		return t.RegisterKYC(stub, args)
	}
	return nil, nil
}

func (t *KycChaincode) RegisterKYC(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var KycDataObj KycData
	//var KycDataList []KycData
	var err error
	var UserPanNumber string

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode
	UserPanNumber = args[0]
	KycDataObj.USER_NAME = args[1]
	KycDataObj.USER_KYC_PDF = args[2]

	fmt.Printf("Input from user:%s\n", KycDataObj)

	//regionTxsAsBytes, err := stub.GetState(UserPanNumber)
	//if err != nil {
	//return nil, errors.New("Failed to get consumer Transactions")
	//}
	//json.Unmarshal(regionTxsAsBytes, &KycDataObj)

	//KycDataList = append(KycDataList, KycDataObj)
	jsonAsBytes, _ := json.Marshal(KycDataObj)

	err = stub.PutState(UserPanNumber, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var resAsBytes []byte
	var UserPanNumber string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	UserPanNumber = args[0]

	resAsBytes, err = t.GetKycDetails(stub, UserPanNumber)

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *KycChaincode) GetKycDetails(stub shim.ChaincodeStubInterface, UserPanNumber string) ([]byte, error) {

	//var requiredObj KycData
	KycTxAsBytes, err := stub.GetState(UserPanNumber)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	//var KycTxObject KycData
	//json.Unmarshal(KycTxAsBytes, &KycTxObject)
	//fmt.Printf("Output from chaincode: %s\n", KycTxObject)

	//res, err := json.Marshal(KycTxAsBytes)
	//if err != nil {
	//return nil, errors.New("Failed to Marshal the required Obj")
	//}
	return KycTxAsBytes, nil

}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
