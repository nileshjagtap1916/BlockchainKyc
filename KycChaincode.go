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
	} else if function == "InsertBankDetails" {
		// save BankDetails in blockchain
		return SaveBankDetails(stub, args)
	} else if function == "UpdateKycDetails" {
		// save BankDetails in blockchain
		return UpdateKyc(stub, args)
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KycChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions of KnowYourCust 
	if function == "GetAllKyc" {
		// get User's KYC data by UserId from blockchain
		return GetAllKyc(stub, args)
	} else if function == "GetKycByUserId" {
		// get User's KYC data by UserId from blockchain
		return GetKycByUserId(stub, args)
	} else if function == "GetKycByBankName" {
		// get User's KYC data by BankName from blockchain
		return GetKycByBankName(stub, args)
	} else if function == "GetKycByExpiringMonth" {
		// get User's KYC data by ExpiringMonth from blockchain
		return GetKycByExpiringMonth(stub, args)
	} else if function == "GetKycByCreatedMonth" {
		// get User's KYC data by CreatedMonth from blockchain
		return GetKycByCreatedMonth(stub, args)
	} else if function == "GetKycCount" {
		// get User's KYC Count
		return GetKycCount(stub, args)
	}

	return nil, nil
}

func main() {
	err := shim.Start(new(KycChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
