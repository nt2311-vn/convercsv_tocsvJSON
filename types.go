package main

type Invoice struct {
	InternalID string  `json:"internal_id"`
	AppliedAmt float64 `json:"applied_amt"`
}

type Journal struct {
	InternalID string  `json:"internal_id"`
	AppliedAmt float64 `json:"applied_amt"`
}

type PaymentInfo struct {
	PaymentRef    string  `json:"payment_ref"`
	PaymentDate   string  `json:"payment_date"`
	CustomerCode  string  `json:"customer_code"`
	Memo          string  `json:"memo"`
	BankName      string  `json:"bank_name"`
	LocationCode  string  `json:"location_code"`
	PaymentAmount float64 `json:"payment_amount"`
}

type GroupRecords struct {
	PaymentObj PaymentInfo `json:"payment_info"`
	Invoices   []Invoice   `json:"invoices"`
	Journals   []Journal   `json:"journals"`
}
