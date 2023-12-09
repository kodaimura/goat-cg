package query

import (
	"database/sql"

	"goat-cg/internal/dto"
	"goat-cg/internal/core/db"
)


type TableQuery interface {
	GetTableLog(id int) ([]dto.TableLog, error)
}


type tableQuery struct {
	db *sql.DB
}


func NewTableQuery() TableQuery {
	db := db.GetDB()
	return &tableQuery{db}
}


func (que *tableQuery)GetTableLog(id int) ([]dto.TableLog, error){
	var ret []dto.TableLog	
	rows, err := que.db.Query(
		`SELECT 
			tl.table_id,
			tl.table_name,
			tl.table_name_logical,
			tl.del_flg,
			tl.create_user_id,
			u1.username create_username,
			tl.update_user_id,
			u2.username update_username,
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
		x := dto.TableLog{}
		err = rows.Scan(
			&x.TableId, 
			&x.TableName,
			&x.TableNameLogical,
			&x.DelFlg,
			&x.CreateUserId,
			&x.CreateUsername,
			&x.UpdateUserId,
			&x.UpdateUsername,
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