package httppayload

type EmpPerformanceGetRequest struct {
	Nik string `json:"nik"`
}
type EmpPerformanceGetResponse struct {
	Nik          string  `json:"nik"`
	Period       string  `json:"period"`
	ScoreMonthly float64 `json:"score_monthly"`
}
