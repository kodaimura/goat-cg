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
			tl.table_id,
			tl.table_name,
			tl.table_name_logical,
			tl.del_flg,
			tl.create_user_id,
			u1.user_name create_user_name,
			tl.update_user_id,
			u2.user_name update_user_name,
			tl.created_at ,
			tl.updated_at
		 FROM 
			 table_def_log tl
			 LEFT OUTER JOIN users u1 ON tl.create_user_id = u1.user_id
			 LEFT OUTER JOIN users u2 ON tl.update_user_id = u2.user_id
		 WHERE 
			 tl.table_id = ?
		 ORDER BY tl.updated_at`, 
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
			&x.CreatedAt,
			&x.UpdatedAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, x)
	}

	return ret, err
}