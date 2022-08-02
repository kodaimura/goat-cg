package queryservice

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ColumnQueryService interface {
	QueryColumns(tableId int) ([]entity.Column, error)
	QueryColumn(id int) (entity.Column, error)
}


type columnQueryService struct {
	db *sql.DB
}


func NewColumnQueryService() ColumnQueryService {
	db := db.GetDB()
	return &columnQueryService{db}
}


func (qs *columnQueryService) QueryColumns(
	tableId int,
) ([]entity.Column, error) {
	
	var ret []entity.Column
	rows, err := qs.db.Query(
		`SELECT 
			column_id,
			table_id, 
			column_name,
			column_name_logical,
			data_type_cls,
			data_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
			create_at,
			update_at
		 FROM
		 	columns
		 WHERE
		 	table_id = ?`,
		 tableId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := entity.Column{}
		err = rows.Scan(
			&c.ColumnId,
			&c.TableId,
			&c.ColumnName, 
			&c.ColumnNameLogical,
			&c.DataTypeCls,
			&c.DataByte,
			&c.PrimaryKeyFlg,
			&c.NotNullFlg,
			&c.UniqueFlg,
			&c.Remark,
			&c.CreateUserId,
			&c.CreateUserId,
			&c.CreateAt,
			&c.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, c)
	}

	return ret, err
}


func (qs *columnQueryService) QueryColumn(id int) (entity.Column, error) {
	var ret entity.Column
	err := qs.db.QueryRow(
		`SELECT 
			column_id,
			table_id, 
			column_name,
			column_name_logical,
			data_type_cls,
			data_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
			create_at,
			update_at
		 FROM
		 	columns
		 WHERE
		 	column_id = ?`,
		 id,
	).Scan(
		&ret.ColumnId,
		&ret.TableId,
		&ret.ColumnName, 
		&ret.ColumnNameLogical,
		&ret.DataTypeCls,
		&ret.DataByte,
		&ret.PrimaryKeyFlg,
		&ret.NotNullFlg,
		&ret.UniqueFlg,
		&ret.Remark,
		&ret.CreateUserId,
		&ret.CreateUserId,
		&ret.CreateAt,
		&ret.UpdateAt,
	)

	return ret, err
}