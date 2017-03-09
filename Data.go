package main

type KycData struct {
	USER_NAME           string `json:"USER_NAME"`
	USER_ID             string `json:"USER_ID"`
	KYC_BANK_NAME       string `json:"KYC_BANK_NAME"`
	KYC_CREATE_DATE     string `json:"KYC_CREATE_DATE"`
	KYC_VALID_TILL_DATE string `json:"KYC_VALID_TILL_DATE"`
	KYC_DOC_BLOB        string `json:"KYC_DOC_BLOB"`
}

type KycCount struct {
	AllContracts      int
	ExperingContracts int
	CreatedContracts  int
}
