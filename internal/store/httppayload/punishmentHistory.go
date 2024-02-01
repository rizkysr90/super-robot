package httppayload

type PunishmentHistoriesGetRequest struct {
	Nik string `json:"nik"`
}

type PunishmentHistoriesGetResponse struct {
	RegNo          string `json:"reg_no"`
	RegDate        string `json:"reg_date"`
	Nik            string `json:"nik"`
	PunishmentType string `json:"punishment_type"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	JG             string `json:"jg"`
	Status         string `json:"status"`
}
