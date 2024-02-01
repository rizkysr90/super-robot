package store

import (
	"api-iad-ams/internal/store/httppayload"
	"context"
	"database/sql"
)

type EmpPerformanceData struct {
	Nik          string
	Name         sql.NullString
	Region       sql.NullString
	JobTitle     sql.NullString
	WorkLocation sql.NullString
	Period       string
	ScoreMonthly sql.NullFloat64
}

func (e *EmpPerformanceData) ToEmpPerformancesGetResponse() *httppayload.EmpPerformanceGetResponse {
	res := &httppayload.EmpPerformanceGetResponse{
		Nik:          e.Nik,
		Period:       e.Period,
		ScoreMonthly: e.ScoreMonthly.Float64,
	}
	return res
}

type EmpPerformance interface {
	Find(ctx context.Context,
		data *httppayload.EmpPerformanceGetRequest) ([]EmpPerformanceData, error)
}
