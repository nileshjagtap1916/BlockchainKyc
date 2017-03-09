package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Region Chaincode implementation
type KycChaincode struct {
}

func (t *KycChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	//Create database on blockchain
	err = InitializeChaincode(stub)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Add user KYC data in Blockchain
func (t *KycChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "InsertKycDetails" {
		// Insert User's KYC data in blockchain
		return SaveKycDetails(stub, args)
	}
	/*if function == "InsertBankDetails" {
		// save BankDetails in blockchain
		return SaveBankDetails(stub, args)
	}*/

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "GetKycDetails" {
		// get User's KYC data from blockchain
		return GetKyc(stub, args)
	}

	return nil, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
