package mysql

import (
	"api-iad-ams/internal/store"
	"context"
	"database/sql"

	"github.com/bfi-finance/bfi-go-pkg/sqldb"
)

type PunishmentHistory struct {
	db *sql.DB
}

func NewPunishmentHistory(db *sql.DB) *PunishmentHistory {
	return &PunishmentHistory{db: db}
}

const punishmentHistoryFindAllByNIKSQL = `
	SELECT 
		Reg_No,
		Reg_Date,
		ID,
		Start,
		End,
		JG,
		Status,
		Punishment_Type
	FROM punishment 
	WHERE Status = 'Approved' AND
	ID = ? 
	ORDER BY Reg_Date DESC
`

func (p *PunishmentHistory) FindByID(ctx context.Context,
	ID string) ([]store.PunishmentHistoryData, error) {
	rows, err := sqldb.WithinTxContextOrDB(ctx, p.db).QueryContext(ctx,
		punishmentHistoryFindAllByNIKSQL, ID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	listPunishmentHistoryData := []store.PunishmentHistoryData{}
	var data store.PunishmentHistoryData
	for rows.Next() {
		err = rows.Scan(&data.RegNo, &data.RegDate, &data.ID,
			&data.Start, &data.End, &data.JG, &data.Status, &data.PunishmentType,
		)
		if err != nil {
			return nil, err
		}
		listPunishmentHistoryData = append(listPunishmentHistoryData, data)
	}
	defer rows.Close()
	if len(listPunishmentHistoryData) == 0 {
		// Handle case when no rows were returned (empty result set)
		return nil, sql.ErrNoRows
	}
	return listPunishmentHistoryData, nil
}
