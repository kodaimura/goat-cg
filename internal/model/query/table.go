package query

import (
	"database/sql"

	"goat-cg/internal/shared/dto"
	"goat-cg/internal/core/db"
)


type TableQuery interface {
	QueryTableLog(id int) ([]dto.QueOutTableLog, error)
}


type tableQuery struct {
	db *sql.DB
}


func NewTableQuery() TableQuery {
	db := db.GetDB()
	return &tableQuery{db}
}


func (que *tableQuery)QueryTableLog(id int) ([]dto.QueOutTableLog, error){
	var ret []dto.QueOutTableLog	
	rows, err := que.db.Query(
		`SELECT 
			td.table_id,
			td.table_name,
			td.table_name_logical,
			td.del_flg,
			td.create_user_id,
			u1.user_name create_user_name,
			td.update_user_id,
			u2.user_name update_user_name,
			td.create_at,
			td.update_at
		 FROM 
		 	table_def td
		 	LEFT OUTER JOIN users u1 ON td.create_user_id = u1.user_id
		 	LEFT OUTER JOIN users u2 ON td.update_user_id = u2.user_id
		 WHERE 
		 	td.table_id = ?
		 ORDER BY td.update_at`, 
		 id,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		x := dto.QueOutTableLog{}
		err = rows.Scan(
			&x.TableId, 
			&x.TableName,
			&x.TableNameLogical,
			&x.DelFlg,
			&x.CreateUserId,
			&x.CreateUserName,
			&x.UpdateUserId,
			&x.UpdateUserName,
			&x.CreateAt,
			&x.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, x)
	}

	return ret, err
}