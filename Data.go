package main

type KycData struct {
	UserId        string `json:"USER_ID"`
	BankName      string `json:"BANK_NAME"`
	UserName      string `json:"USER_NAME"`
	CreateDate    string `json:"KYC_CREATE_DATE"`
	ValidTillDate string `json:"KYC_VALID_TILL_DATE"`
	KycDocument   string `json:"KYC_DOC_BLOB"`
}
