package main

type KycData struct {
	UserId        string `json:"UserId"`
	BankName      string `json:"BankName"`
	UserName      string `json:"UserName"`
	CreateDate    string `json:"CreateDate"`
	ValidTillDate string `json:"ValidTillDate"`
	KycDocument   string `json:"KycDocument"`
}
