package main

type KycData struct {
	USER_NAME           string `json:"userName"`
	USER_ID             string `json:"userId"`
	KYC_BANK_NAME       string `json:"kycBankName"`
	KYC_CREATE_DATE     string `json:"kycCreateDate"`
	KYC_VALID_TILL_DATE string `json:"kycValidTillDate"`
	KYC_DOC_BLOB        string `json:"kycDocBlob"`
	KYC_INFO_1          string `json:"kycInfo1"`
	KYC_INFO_2          string `json:"kycInfo2"`
	KYC_INFO_3          string `json:"kycInfo3"`
	KYC_INFO_4          string `json:"kycInfo4"`
}

type KycCount struct {
	AllContracts      int
	ExpiringContracts int
	CreatedContracts  int
}
