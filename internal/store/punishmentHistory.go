package store

import (
	"api-iad-ams/internal/store/httppayload"
	"context"
	"database/sql"
)

type PunishmentHistoryData struct {
	RegNo          string
	Remarks        string
	Status         string
	RegDate        sql.NullString
	ID             sql.NullString
	JoinDate       sql.NullString
	ResignDate     sql.NullString
	Directorate    sql.NullString
	Division       sql.NullString
	Department     sql.NullString
	JobTitle       sql.NullString
	JG             sql.NullString
	PunishmentType sql.NullString
	Start          sql.NullTime
	End            sql.NullTime
	Description    sql.NullString
	Classification sql.NullString
}

func (e *PunishmentHistoryData) ToGetAllPunishmentHistoryByNik() *httppayload.PunishmentHistoriesGetResponse {
	res := &httppayload.PunishmentHistoriesGetResponse{
		RegNo:          e.RegNo,
		RegDate:        e.RegDate.String,
		Nik:            e.ID.String,
		PunishmentType: e.PunishmentType.String,
		StartDate:      e.Start.Time.String(),
		EndDate:        e.End.Time.String(),
		JG:             e.JG.String,
		Status:         e.Status,
	}
	return res
}

type PunishmentHistory interface {
	FindByID(ctx context.Context,
		id string) ([]PunishmentHistoryData, error)
}
