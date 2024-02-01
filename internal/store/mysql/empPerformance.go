package mysql

import (
	"api-iad-ams/internal/store"
	"api-iad-ams/internal/store/httppayload"
	"context"
	"database/sql"

	"github.com/bfi-finance/bfi-go-pkg/sqldb"
)

type EmpPerformance struct {
	db *sql.DB
}

func NewEmpPerformance(db *sql.DB) *EmpPerformance {
	return &EmpPerformance{db: db}
}

// LIMIT 7, is because to get last six month (delay 1 month)
const empPerformanceFindAllByNIKSQL = `
	SELECT 
		NIK,
		Period,
		ScoreMonthly
	FROM infolmpperformance 
	WHERE NIK = ? 
	ORDER BY Period DESC
	LIMIT 7 
`

func (e *EmpPerformance) Find(ctx context.Context,
	data *httppayload.EmpPerformanceGetRequest) ([]store.EmpPerformanceData, error) {
	rows, err := sqldb.WithinTxContextOrDB(ctx, e.db).QueryContext(ctx, empPerformanceFindAllByNIKSQL, data.Nik)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	listEmpPerformanceData := []store.EmpPerformanceData{}
	var temp store.EmpPerformanceData
	for rows.Next() {
		err = rows.Scan(&temp.Nik, &temp.Period, &temp.ScoreMonthly)
		if err != nil {
			return nil, err
		}
		listEmpPerformanceData = append(listEmpPerformanceData, temp)
	}
	defer rows.Close()
	if len(listEmpPerformanceData) == 0 {
		// Handle case when no rows were returned (empty result set)
		return nil, sql.ErrNoRows
	}
	return listEmpPerformanceData, nil
}
