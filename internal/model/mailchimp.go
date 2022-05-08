package model

type ApiResp struct {
	Members    []MailchimpMembers `json:"members"`
	TotalItems int64              `json:"total_items"`
}

type MailchimpMembers struct {
	ID           string      `json:"id"`
	LastChanged  string      `json:"last_changed"`
	EmailAddress string      `json:"email_address"`
	Status       string      `json:"status"`
	MergeFields  MergeFields `json:"merge_fields"`
}

type MergeFields struct {
	FirstName string `json:"FNAME"`
	LastName  string `json:"LNAME"`
}

type APIReq struct {
	Limit       int64
	Offset      int64
	LastChanged string
	URL         string
	APIKey      string
	ListID      string
}
