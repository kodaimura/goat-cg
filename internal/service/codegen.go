package service

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"os/exec"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/utils"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type CodegenService interface {
	GenerateGoat(rdbms string, tableIds []int) string
}


type codegenService struct {
	columnRepository repository.ColumnRepository
	tableRepository repository.TableRepository
}


func NewCodegenService() CodegenService {
	columnRepository := repository.NewColumnRepository()
	tableRepository := repository.NewTableRepository()
	return &codegenService{columnRepository, tableRepository}
}


// Generate goat source and return zip path.
// param rdbms: "sqlite3" or "postgresql" 
func (srv *codegenService) GenerateGoat(rdbms string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7)

	srv.generateSource(rdbms, tableIds, path)

	if err := exec.Command("zip", "-rm", path + ".zip", path).Run(); err != nil {
		logger.Error(err.Error())
	}

	return path + ".zip"
}


func (srv *codegenService) generateSource(rdbms string, tableIds []int, rootPath string) {
	path := rootPath + "/scripts"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return
	}
	srv.generateScriptsSource(rdbms, tableIds, path)

	path = rootPath + "/internal"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return
	}
	srv.generateInternalSource(rdbms, tableIds, path)
}


//dataTypeMapSqlite3 map DataTypeCls and sqlite3 data types.
var dataTypeMapSqlite3 = map[string]string {
	constant.DATA_TYPE_CLS_SERIAL: "INTEGER PRIMARY KEY AUTOINCREMENT",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "TEXT",
	constant.DATA_TYPE_CLS_CHAR: "TEXT",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TEXT",
	constant.DATA_TYPE_CLS_DATE: "TEXT",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dataTypeMapPostgresql map DataTypeCls and postgresql data types.
var dataTypeMapPostgresql = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "SERIAL PRIMARY KEY",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TIMESTAMP",
	constant.DATA_TYPE_CLS_DATE: "DATE",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dataTypeMapMysql map DataTypeCls and mysql data types.
var dataTypeMapMysql = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "INT PRIMARY KEY AUTO_INCREMENT",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INT",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "DATETIME",
	constant.DATA_TYPE_CLS_DATE: "DATE",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dbDataTypeGoTypeMap map DataTypeCls and Golang types.
var dbDataTypeGoTypeMap = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "int",
	constant.DATA_TYPE_CLS_TEXT: "string",
	constant.DATA_TYPE_CLS_VARCHAR: "string",
	constant.DATA_TYPE_CLS_CHAR: "string",
	constant.DATA_TYPE_CLS_INTEGER: "int",
	constant.DATA_TYPE_CLS_NUMERIC: "float64",
	constant.DATA_TYPE_CLS_TIMESTAMP: "string",
	constant.DATA_TYPE_CLS_DATE: "string",
	constant.DATA_TYPE_CLS_BLOB: "string",
}


func (srv *codegenService) writeFile(path, content string) {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.Error(err.Error())
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.Error(err.Error())
	}
}


func (srv *codegenService) extractPrimaryKeys(columns []model.Column) []model.Column {
	var ret []model.Column

	for _, c := range columns {
		if c.PrimaryKeyFlg == constant.FLG_ON {
			ret = append(ret, c)
		}
	}

	return ret
}


// generateScriptsSource generate ddl(create table) source.
// main processing of GenerateDdl.
func (srv *codegenService) generateScriptsSource(rdbms string, tableIds []int, path string) {
	s := srv.generateDdlCreateTables(rdbms, tableIds) + "\n" +
		srv.generateDdlCreateTriggers(rdbms, tableIds)

	srv.writeFile(path + "/create-table.sql", s)
}


func (srv *codegenService) generateDdlCreateTables(rdbms string, tableIds []int) string {
	s := ""
	for _, tid := range tableIds {
		s += srv.generateDdlCreateTable(rdbms, tid) + "\n\n"
	}

	return s
}


func (srv *codegenService) generateDdlCreateTable(rdbms string, tid int) string {
	s := ""
	table, err := srv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	s += "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n" +
		srv.generateDdlColumns(rdbms, tid) + "\n);"

	return s
}


func (srv *codegenService) getValidColumns(tid int) ([]model.Column, error) {
	columns, err := srv.columnRepository.GetByTableId(tid)
	if err != nil {
		logger.Error(err.Error())
		return columns, err
	}

	var cols []model.Column
	for _, col := range columns {
		if (col.DelFlg != 1) {
			cols = append(cols, col)
		} 
	}

	return cols, nil
}


func (srv *codegenService) generateDdlColumns(rdbms string, tid int) string {
	s := ""
	columns, err := srv.getValidColumns(tid)
	if err != nil {
		logger.Error(err.Error())
		return s
	}

	for _, col := range columns {
		s += srv.generateDdlColumn(rdbms, col)
	}
	s += srv.generateDdlCommonColumns(rdbms)
	s += srv.generateDdlPrymaryKey(rdbms, columns)

	return strings.TrimRight(s, ",\n")
}


func (srv *codegenService) generateDdlCommonColumns(rdbms string) string {
	s := ""
	if rdbms == "sqlite3" {
		s = "\tcreated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n" + 
			"\tupdated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n"

	} else if rdbms == "postgresql" {
		s = "\tcreated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"

	} else if rdbms == "mysql" {
		s = "\tcreated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n"
	}

	return s
}


func (srv *codegenService) generateDdlPrymaryKey(rdbms string, columns []model.Column) string {
	s := "" 
	pkcolumns := srv.extractPrimaryKeys(columns)

	for i, c := range pkcolumns {
		if c.DataTypeCls == constant.DATA_TYPE_CLS_SERIAL {
			return ""
		}
		if i == 0 {
			s += "\tPRIMARY KEY("
		} else {
			s += ", "
		}
		s += c.ColumnName
	}
	
	if s != "" {
		s += "),\n"
	}

	return s
}


func (srv *codegenService) generateDdlColumn(rdbms string, column model.Column) string {
	s := "\t" + column.ColumnName + " " + srv.generateDdlColumnDataType(rdbms, column)
	if cts := srv.generateDdlColumnConstraints(column); cts != "" {
		s += " " + cts
	}
	if dflt := srv.generateDdlColumnDefault(column); dflt != "" {
		s += " " + dflt
	}

	return s + ",\n"
}


func (srv *codegenService) generateDdlColumnConstraints(column model.Column) string {
	s := ""
	if column.NotNullFlg == constant.FLG_ON {
		s += "NOT NULL "
	}
	if column.UniqueFlg == constant.FLG_ON {
		s += "UNIQUE"
	} 
	
	return strings.TrimRight(s, " ")
}


func (srv *codegenService) generateDdlColumnDefault(column model.Column) string {
	s := ""
	if column.DefaultValue != "" {
		if column.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC ||
		column.DataTypeCls == constant.DATA_TYPE_CLS_INTEGER {
			s = "DEFAULT " + column.DefaultValue
		} else {
			s = "DEFAULT '" + column.DefaultValue + "'"
		}
	}

	return s
}


func (srv *codegenService) generateDdlColumnDataType(rdbms string, column model.Column) string {
	s := ""
	if rdbms == "sqlite3" {
		s = dataTypeMapSqlite3[column.DataTypeCls]

	} else if rdbms == "postgresql" {
		s = srv.generateDdlColumnDataTypePostgresql(column)
	
	} else if rdbms == "mysql" {
		s = srv.generateDdlColumnDataTypeMysql(column)
	}

	return s
}


func (srv *codegenService) generateDdlColumnDataTypePostgresql(column model.Column) string {
	s := dataTypeMapPostgresql[column.DataTypeCls]

	if column.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	column.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if column.Precision != 0 {
			s += "(" + strconv.Itoa(column.Precision) + ")"
		}	
	}
	if column.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if column.Precision != 0 {
			s += "(" + strconv.Itoa(column.Precision) + "," + strconv.Itoa(column.Scale) + ")"
		}
	}

	return s
}


func (srv *codegenService) generateDdlColumnDataTypeMysql(column model.Column) string {
	s := dataTypeMapMysql[column.DataTypeCls]

	if column.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	column.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if column.Precision != 0 {
			s += "(" + strconv.Itoa(column.Precision) + ")"
		}	
	}
	if column.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if column.Precision != 0 {
			s += "(" + strconv.Itoa(column.Precision) + "," + strconv.Itoa(column.Scale) + ")"
		}
	}

	return s
}


func (srv *codegenService) generateDdlCreateTriggers(rdbms string, tableIds []int) string {
	s := ""
	if rdbms == "postgresql" {
		s += "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	for _, tid := range tableIds {
		s += srv.generateDdlCreateTrigger(rdbms, tid) + "\n\n"
	}

	return s
}


func (srv *codegenService) generateDdlCreateTrigger(rdbms string, tid int) string {
	s := ""
	table, err := srv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	if rdbms == "sqlite3" {
		s = "CREATE TRIGGER IF NOT EXISTS trg_" + table.TableName + "_upd " + 
			"AFTER UPDATE ON " + table.TableName + "\n" +
			"BEGIN\n\tUPDATE " + table.TableName + "\n" +
			"\tSET updated_at = DATETIME('now', 'localtime')\n" + 
			"\tWHERE rowid == NEW.rowid;\nEND;"

	} else if rdbms == "postgresql" {
		s = "CREATE TRIGGER trg_" + table.TableName + "_upd " + 
			"BEFORE UPDATE ON " + table.TableName + " FOR EACH ROW\n" + 
			"\texecute procedure set_update_time();" 
	}

	return s
}


// tableNameToFileName get file name from table name
// user => user.go / USER_TABLE => user_table.go
func (srv *codegenService) tableNameToFileName(tn string) string {
	n := strings.ToLower(tn)
	return n + ".go"
}


//xxx -> Xxx / xxx_yyy -> XxxYyy
func SnakeToPascal(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}


//xxx -> xxx / xxx_yyy -> xxxYyy
func SnakeToCamel(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		if i != 0 {
			ls[i] = strings.ToUpper(s[0:1]) + s[1:]
		}
	}
	return strings.Join(ls, "")
}


//xxx -> x / xxx_yyy -> xy
func GetSnakeInitial(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	ret := ""
	for _, s := range ls {
		ret += s[0:1]
	}
	return ret
}


func (srv *codegenService) generateInternalSource(rdbms string, tableIds []int, path string) {
	modelPath := path + "/model"
	if err := os.MkdirAll(modelPath, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

	repositoryPath := path + "/repository"
	if err := os.MkdirAll(repositoryPath, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

	for _, tid := range tableIds {
		table, err := srv.tableRepository.GetById(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		columns, err := srv.getValidColumns(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		srv.generateModelFile(&table, columns, modelPath)
		srv.generateRepositoryFile(rdbms, &table, columns, repositoryPath)
	}
}


func (srv *codegenService) generateModelFile(table *model.Table, columns []model.Column, path string) {
	path += "/" + srv.tableNameToFileName(table.TableName)
	code := srv.generateModelCode(table, columns)
	srv.writeFile(path, code)
}


func (srv *codegenService) generateModelCode(table *model.Table, columns []model.Column) string {
	s := "package model\n\n\n"

	s += fmt.Sprintf("type %s struct {\n", SnakeToPascal(table.TableName))
	for _, c := range columns {
		s += fmt.Sprintf(
			"\t%s %s `db:\"%s\" json:\"%s\"`\n", 
			SnakeToPascal(c.ColumnName),
			dbDataTypeGoTypeMap[c.DataTypeCls],
			strings.ToLower(c.ColumnName),
			strings.ToLower(c.ColumnName),
		)
	}
	s += "\tCreatedAt string `db:\"created_at\" json:\"created_at\"`\n"
	s += "\tUpdatedAt string `db:\"updated_at\" json:\"updated_at\"`\n"

	return s + "}"
}


func (srv *codegenService) generateRepositoryFile(rdbms string, table *model.Table, columns []model.Column, path string) {
	path += "/" + srv.tableNameToFileName(table.TableName)
	code := srv.generateRepositoryCode(rdbms, table, columns)
	srv.writeFile(path, code)
}


func (srv *codegenService) generateRepositoryCode(rdbms string, table *model.Table, columns []model.Column) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	s := "package repository\n\n\nimport (\n" + 
		"\t\"database/sql\"\n\n\t\"xxxxx/internal/core/db\"\n\t\"xxxxx/internal/model\"\n)\n\n\n"

	s += srv.generateRepositoryInterfaceCode(table, columns)
	
	s += "\n\n" +
		fmt.Sprintf("type %sRepository struct {\n\tdb *sql.DB\n}\n\n\n", tnc) +
		fmt.Sprintf("func New%sRepository() *%sRepository {\n", tnp, tnc) +
		fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sRepository{db}\n}\n\n\n", tnc)

	s += srv.generateRepositoryGet(table, columns) + "\n\n\n"
	if len(srv.extractPrimaryKeys(columns)) > 0 {
		s += srv.generateRepositoryGetByPk(rdbms, table, columns) + "\n\n\n"
		s += srv.generateRepositoryInsertPk(rdbms, table, columns) + "\n\n\n"
		s += srv.generateRepositoryUpdate(rdbms, table, columns) + "\n\n\n"
		s += srv.generateRepositoryDelete(rdbms, table, columns)
	} else {
		s += srv.generateRepositoryInsert(rdbms, table, columns)
	}

	return s
}


func (srv *codegenService) generateInsertReturnType(columns []model.Column) string {
	pkcolumns := srv.extractPrimaryKeys(columns)
	s := "error"
	if len(pkcolumns) > 0 {
		s = "("
		isFirst := true
		for _, column := range pkcolumns {
			if isFirst {
				s += dbDataTypeGoTypeMap[column.DataTypeCls]
				isFirst = false
			} else {
				s += ", " + dbDataTypeGoTypeMap[column.DataTypeCls] 
			}
		}
		s += ", error)"
	}

	return s
}


// return "type *Repository interface { ... }"
func (srv *codegenService) generateRepositoryInterfaceCode(table *model.Table, columns []model.Column) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)

	s := fmt.Sprintf("type %sRepository interface {\n", tnp) +
		fmt.Sprintf("\tGet() ([]model.%s, error)\n", tnp)

	if len(srv.extractPrimaryKeys(columns)) > 0 {
		s += fmt.Sprintf("\tGetByPk(%s *model.%s) (model.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tInsert(%s *model.%s) %s\n", tni, tnp, srv.generateInsertReturnType(columns)) +
		fmt.Sprintf("\tUpdate(%s *model.%s) error\n", tni, tnp) +
		fmt.Sprintf("\tDelete(%s *model.%s) error\n", tni, tnp)
	} else {
		s += fmt.Sprintf("\tInsert(%s *model.%s) error\n", tni, tnp,)
	}
	s += "}\n"
	return s
		
}


// generateRepositoryGet generate repository function 'Get'.
// return "func (ur *userRepository) Get() ([]model.User, error) {...}"
func (srv *codegenService) generateRepositoryGet(table *model.Table, columns []model.Column) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf("func (%sr *%sRepository) Get() ([]model.%s, error) {\n", tni, tnc, tnp) +
		fmt.Sprintf("\trows, err := %sr.db.Query(\n", tni)

	s += "\t\t`SELECT\n"
	for i, c := range columns {
		if i == 0 {
			s += fmt.Sprintf("\t\t\t%s", c.ColumnName)
		} else {
			s += fmt.Sprintf("\n\t\t\t,%s", c.ColumnName)
		}
	}
	s += "\n\t\t\t,created_at\n\t\t\t,updated_at" +
		fmt.Sprintf("\n\t\t FROM %s`,\n\t)\n\tdefer rows.Close()\n\n", tn) +
		"\tif err != nil {\n\t\treturn nil, err\n\t}\n\n" +
		fmt.Sprintf("\tret := []model.%s{}\n", tnp) +
		"\tfor rows.Next() {\n" +
		fmt.Sprintf("\t\t%s := model.%s{}\n\t\terr = rows.Scan(\n", tni, tnp)

	for _, c := range columns {
		s += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
	}
	s += fmt.Sprintf("\t\t\t&%s.CreatedAt,\n", tni) + fmt.Sprintf("\t\t\t&%s.UpdatedAt,\n", tni)

	s += fmt.Sprintf("\t\t)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tret = append(ret, %s)\n", tni) +
		"\t}\n\n\treturn ret, nil\n}"

	return s
}


// generateRepositoryGetByPk generate repository function 'GetByPk'.
// return "func (ur *userRepository) GetByPk(u *model.User) (model.User, error) {...}"
func (srv *codegenService) generateRepositoryGetByPk(
	rdbms string, table *model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (%sr *%sRepository) GetByPk(%s *model.%s) (model.%s, error) {\n", 
		tni, tnc, tni, tnp, tnp,
	) + fmt.Sprintf("\tvar ret model.%s\n\n\terr := %sr.db.QueryRow(\n", tnp, tni)

	bindCount := 0
	s += "\t\t`SELECT\n"
	for i, c := range columns {
		if i == 0 {
			s += fmt.Sprintf("\t\t\t%s", c.ColumnName)
		} else {
			s += fmt.Sprintf("\n\t\t\t,%s", c.ColumnName)
		}
	}
	s += "\n\t\t\t,created_at\n\t\t\t,updated_at" + fmt.Sprintf("\n\t\t FROM %s\n", tn)
	s += srv.generateRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"
	s += srv.generateRepositoryWhereClauseBindVals(table, columns)
	s += "\t).Scan(\n"
	for _, c := range columns {
		s += fmt.Sprintf("\t\t&ret.%s,\n", SnakeToPascal(c.ColumnName))
	}
	s += "\t\t&ret.CreatedAt,\n\t\t&ret.UpdatedAt,\n" +
		"\t)\n\n\treturn ret, err\n}"

	return s
}


func (srv *codegenService) getBindVar(rdbms string, n int) string {
	if rdbms == "postgresql" {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}


// concatBindVariableWithCommas return ?,?,?,?,... or $1,$2,$3,$4,...
func (srv *codegenService) concatBindVariableWithCommas(rdbms string, bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, srv.getBindVar(rdbms, i))
	}
	return strings.Join(ls, ",")
}


// generateRepositoryInsert generate repository function 'Insert'.
// return "func (ur *userRepository) Insert(u *model.User) error {...}"
func (srv *codegenService) generateRepositoryInsert(
	rdbms string, table *model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (%sr *%sRepository) Insert(%s *model.%s) error {\n", 
		tni, tnc, tni, tnp,
	) + fmt.Sprintf("\t_, err := %sr.db.Exec(\n", tni) +fmt.Sprintf("\t\t`INSERT INTO %s (\n", tn)

	bindCount := 0
	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("\t\t\t%s", c.ColumnName)
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s", c.ColumnName)
			}
		}	
	}
	s += fmt.Sprintf("\n\t\t ) VALUES(%s)`,\n", srv.concatBindVariableWithCommas(rdbms, bindCount))

	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateRepositoryInsertPk generate repository function 'Insert'.
// return "func (ur *userRepository) Insert(u *model.User) (int, error) {...}"
func (srv *codegenService) generateRepositoryInsertPk(
	rdbms string, table *model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (%sr *%sRepository) Insert(%s *model.%s) %s {\n", 
		tni, tnc, tni, tnp, srv.generateInsertReturnType(columns),
	)

	pkcolumns := srv.extractPrimaryKeys(columns)
	for _, column := range pkcolumns {
		s += fmt.Sprintf(
			"\tvar %s %s\n", 
			SnakeToCamel(column.ColumnName), dbDataTypeGoTypeMap[column.DataTypeCls],
		)
	}
	s += fmt.Sprintf("\n\terr := %sr.db.QueryRow(\n", tni) +fmt.Sprintf("\t\t`INSERT INTO %s (\n", tn)

	bindCount := 0
	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("\t\t\t%s", c.ColumnName)
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s", c.ColumnName)
			}
		}	
	}
	s += fmt.Sprintf("\n\t\t ) VALUES(%s)\n\t\t RETURNING ", srv.concatBindVariableWithCommas(rdbms, bindCount))
	isFirst := true
	for _, column := range pkcolumns {
		if isFirst {
			s += column.ColumnName
			isFirst = false
		} else {
			s += ", " + column.ColumnName
		}
	}
	s += "`,\n"

	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	s += "\t).Scan(\n"
	for _, column := range pkcolumns {
		s += fmt.Sprintf("\t\t&%s,\n", SnakeToCamel(column.ColumnName))
	}

	s += "\t)\n\n\treturn "
	for _, column := range pkcolumns {
		s += fmt.Sprintf("%s, ", SnakeToCamel(column.ColumnName))
	}
	s += "err\n}"

	return s
}


// generateRepositoryUpdate generate repository function 'Update'.
// return "func (ur *userRepository) Update(u *model.User) error {...}"
func (srv *codegenService) generateRepositoryUpdate(
	rdbms string, table *model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	
	s := fmt.Sprintf(
		"func (%sr *%sRepository) Update(%s *model.%s) error {\n", 
		tni, tnc, tni, tnp,
	) + fmt.Sprintf("\t_, err := %sr.db.Exec(\n", tni) + fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET ", tn)

	bindCount := 0
	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		c.PrimaryKeyFlg != constant.FLG_ON {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("%s = %s", c.ColumnName, srv.getBindVar(rdbms, bindCount))
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s = %s", c.ColumnName, srv.getBindVar(rdbms, bindCount))
			}
		}
	}
	s += "\n"
	s += srv.generateRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"

	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		c.PrimaryKeyFlg != constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	s += srv.generateRepositoryWhereClauseBindVals(table, columns)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateRepositoryDelete generate repository function 'Delete'.
// return "func (ur *userRepository) Delete(u *model.User) error {...}"
func (srv *codegenService) generateRepositoryDelete(
	rdbms string, table *model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	
	s := fmt.Sprintf(
		"func (%sr *%sRepository) Delete(%s *model.%s) error {\n", 
		tni, tnc, tni, tnp,
	) + fmt.Sprintf("\t_, err := %sr.db.Exec(\n", tni) + fmt.Sprintf("\t\t`DELETE FROM %s\n", tn)

	bindCount := 0
	s += srv.generateRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"
	s += srv.generateRepositoryWhereClauseBindVals(table, columns)
	s += "\t)\n\n\treturn err\n}"

	return s
}


func (srv *codegenService) generateRepositoryWhereClause(
	rdbms string, columns []model.Column, bindCount *int,
) string {
	s := "\t\t WHERE "

	isFirst := true
	for _, c := range columns {
		if c.PrimaryKeyFlg == constant.FLG_ON {
			*bindCount += 1
			if isFirst {
				s += fmt.Sprintf("%s = %s", c.ColumnName, srv.getBindVar(rdbms, *bindCount))
				isFirst = false
			} else {
				s += fmt.Sprintf("\n\t\t   AND %s = %s", c.ColumnName, srv.getBindVar(rdbms, *bindCount))
			}
		}
	}

	return s
}


func (srv *codegenService) generateRepositoryWhereClauseBindVals(
	table *model.Table, columns []model.Column,
) string {
	s := ""
	tni := GetSnakeInitial(table.TableName)

	for _, c := range columns {
		if c.PrimaryKeyFlg == constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	return s
}
