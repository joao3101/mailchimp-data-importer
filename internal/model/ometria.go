package model

type Users struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

type OmetriaResponse struct {
	Status   string `json:"status"`
	Response int64  `json:"response"`
}
