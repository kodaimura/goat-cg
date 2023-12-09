package query

import (
	"database/sql"

	"goat-cg/internal/dto"
	"goat-cg/internal/core/db"
)


type ColumnQuery interface {
	QueryColumnLog(id int) ([]dto.QueOutColumnLog, error)
}


type columnQuery struct {
	db *sql.DB
}


func NewColumnQuery() ColumnQuery {
	db := db.GetDB()
	return &columnQuery{db}
}


func (que *columnQuery)QueryColumnLog(id int) ([]dto.QueOutColumnLog, error){
	var ret []dto.QueOutColumnLog

	rows, err := que.db.Query(
		`SELECT 
			cl.column_id,
			cl.table_id, 
			cl.column_name,
			cl.column_name_logical,
			cl.data_type_cls,
			cl.precision,
			cl.scale,
			cl.primary_key_flg,
			cl.not_null_flg,
			cl.unique_flg,
			cl.default_value,
			cl.remark,
			cl.align_seq,
			cl.del_flg,
			cl.create_user_id,
			u1.username create_username,
			cl.update_user_id,
			u2.username update_username,
			cl.created_at ,
			cl.updated_at
		 FROM
			 column_def_log cl
			 LEFT OUTER JOIN users u1 ON cl.create_user_id = u1.user_id
			 LEFT OUTER JOIN users u2 ON cl.update_user_id = u2.user_id
		 WHERE 
			 cl.column_id = ?
		 ORDER BY cl.updated_at`,
		 id,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		x := dto.QueOutColumnLog{}
		err = rows.Scan(
			&x.ColumnId,
			&x.TableId,
			&x.ColumnName, 
			&x.ColumnNameLogical,
			&x.DataTypeCls,
			&x.Precision,
			&x.Scale,
			&x.PrimaryKeyFlg,
			&x.NotNullFlg,
			&x.UniqueFlg,
			&x.DefaultValue,
			&x.Remark,
			&x.AlignSeq,
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