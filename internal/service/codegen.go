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
	GenerateDdl(rdbms string, tableIds []int) string
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


// Generate ddl source and return file path.
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) GenerateDdl(rdbms string, tableIds []int) string {
	path := "./tmp/ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7) + ".sql"

	serv.generateDdlSource(rdbms, tableIds, path)

	return path
}


// Generate goat source and return zip path.
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) GenerateGoat(rdbms string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7)

	serv.generateGoatSource(rdbms, tableIds, path)

	if err := exec.Command("zip", "-rm", path + ".zip", path).Run(); err != nil {
		logger.Error(err.Error())
	}

	return path + ".zip"
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


func (serv *codegenService) writeFile(path, content string) {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.Error(err.Error())
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.Error(err.Error())
	}
}


func (serv *codegenService) extractPrimaryKeys(columns []model.Column) []model.Column {
	var ret []model.Column

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			ret = append(ret, col)
		}
	}

	return ret
}


///////////////////////
/// GenerateDdl     ///
///////////////////////

// generateDdlSource generate ddl(create table) source.
// main processing of GenerateDdl.
func (serv *codegenService) generateDdlSource(rdbms string, tableIds []int, path string) {
	s := serv.generateDdlCreateTables(rdbms, tableIds) + "\n" +
		serv.generateDdlCreateTriggers(rdbms, tableIds)

	serv.writeFile(path, s)
}


func (serv *codegenService) generateDdlCreateTables(rdbms string, tableIds []int) string {
	s := ""
	for _, tid := range tableIds {
		s += serv.generateDdlCreateTable(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) generateDdlCreateTable(rdbms string, tid int) string {
	s := ""
	table, err := serv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	s += "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n" +
		serv.generateDdlColumns(rdbms, tid) + "\n);"

	return s
}


func (serv *codegenService) generateDdlColumns(rdbms string, tid int) string {
	s := ""
	columns, err := serv.columnRepository.GetByTableId(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	for _, col := range columns {
		s += serv.generateDdlColumn(rdbms, col)
	}
	s += serv.generateDdlCommonColumns(rdbms)
	s += serv.generateDdlPrymaryKey(rdbms, columns)

	return strings.TrimRight(s, ",\n")
}


func (serv *codegenService) generateDdlCommonColumns(rdbms string) string {
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


func (serv *codegenService) generateDdlPrymaryKey(rdbms string, columns []model.Column) string {
	s := "" 
	pkcolumns := serv.extractPrimaryKeys(columns)

	for i, col := range pkcolumns {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_SERIAL {
			return ""
		}

		if i == 0 {
			s += "\tPRIMARY KEY("
		} else {
			s += ", "
		}
		s += col.ColumnName
	}
	
	if s != "" {
		s += "),\n"
	}

	return s
}


func (serv *codegenService) generateDdlColumn(rdbms string, col model.Column) string {
	s := "\t" + col.ColumnName + " " + serv.generateDdlColumnDataType(rdbms, col)
	if cts := serv.generateDdlColumnConstraints(col); cts != "" {
		s += " " + cts
	}
	if dflt := serv.generateDdlColumnDefault(col); dflt != "" {
		s += " " + dflt
	}

	return s + ",\n"
}


func (serv *codegenService) generateDdlColumnConstraints(col model.Column) string {
	s := ""
	if col.NotNullFlg == constant.FLG_ON {
		s += "NOT NULL "
	}
	if col.UniqueFlg == constant.FLG_ON {
		s += "UNIQUE"
	} 
	
	return strings.TrimRight(s, " ")
}


func (serv *codegenService) generateDdlColumnDefault(col model.Column) string {
	s := ""
	if col.DefaultValue != "" {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC ||
		col.DataTypeCls == constant.DATA_TYPE_CLS_INTEGER {
			s = "DEFAULT " + col.DefaultValue
		} else {
			s = "DEFAULT '" + col.DefaultValue + "'"
		}
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataType(rdbms string, col model.Column) string {
	s := ""
	if rdbms == "sqlite3" {
		s = dataTypeMapSqlite3[col.DataTypeCls]

	} else if rdbms == "postgresql" {
		s = serv.generateDdlColumnDataTypePostgresql(col)
	
	} else if rdbms == "mysql" {
		s = serv.generateDdlColumnDataTypeMysql(col)
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataTypePostgresql(col model.Column) string {
	s := dataTypeMapPostgresql[col.DataTypeCls]

	if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + ")"
		}	
	}
	if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + "," + strconv.Itoa(col.Scale) + ")"
		}
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataTypeMysql(col model.Column) string {
	s := dataTypeMapMysql[col.DataTypeCls]

	if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + ")"
		}	
	}
	if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + "," + strconv.Itoa(col.Scale) + ")"
		}
	}

	return s
}


func (serv *codegenService) generateDdlCreateTriggers(rdbms string, tableIds []int) string {
	s := ""
	if rdbms == "postgresql" {
		s += "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	for _, tid := range tableIds {
		s += serv.generateDdlCreateTrigger(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) generateDdlCreateTrigger(rdbms string, tid int) string {
	s := ""
	table, err := serv.tableRepository.GetById(tid)

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


////////////////////////
/// GenerateGoat     ///
////////////////////////

// tableNameToFileName get file name from table name
// user => user.go / USER_TABLE => user_table.go
func (serv *codegenService) tableNameToFileName(tn string) string {
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
		ret = s[0:1]
	}
	return ret
}


func (serv *codegenService) generateGoatSource(rdbms string, tableIds []int, rootPath string) {
	path := rootPath + "/internal"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

	serv.generateGoatInternalSource(rdbms, tableIds, path)
}


func (serv *codegenService) generateGoatInternalSource(rdbms string, tableIds []int, path string) {
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
		table, err := serv.tableRepository.GetById(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		columns, err := serv.columnRepository.GetByTableId(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		serv.generateGoatModelFile(table, columns, modelPath)
		serv.generateGoatRepositoryFile(rdbms, table, columns, repositoryPath)
	}
}


func (serv *codegenService) generateGoatModelFile(table model.Table, columns []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(table.TableName)
	code := serv.generateGoatModelCode(table, columns)
	serv.writeFile(path, code)
}


func (serv *codegenService) generateGoatModelCode(table model.Table, columns []model.Column) string {
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
	s += "\tCreatedAt string `db:\"created_at \" json:\"created_at \"`\n"
	s += "\tUpdatedAt string `db:\"updated_at\" json:\"updated_at\"`\n"

	return s + "}"
}


func (serv *codegenService) generateGoatRepositoryFile(rdbms string, table model.Table, columns []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(table.TableName)
	code := serv.generateGoatRepositoryCode(rdbms, table, columns)
	serv.writeFile(path, code)
}


func (serv *codegenService) generateGoatRepositoryCode(rdbms string, table model.Table, columns []model.Column) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	s := "package repository\n\n\nimport (\n" + 
		"\t\"database/sql\"\n\n\t\"xxxxx/internal/core/db\"\n\t\"xxxxx/internal/model\"\n)\n\n\n"

	s += serv.generateGoatRepositoryInterfaceCode(table)
	
	s += "\n\n" +
		fmt.Sprintf("type %sRepository struct {\n\tdb *sql.DB\n}\n\n\n", tnc) +
		fmt.Sprintf("func New%sRepository() *%sRepository {\n", tnp, tnc) +
		fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sRepository{db}\n}\n\n\n", tnc)

	s += serv.generateGoatRepositoryGet(table, columns) + "\n\n\n"
	s += serv.generateGoatRepositoryGetByPk(rdbms, table, columns) + "\n\n\n"
	s += serv.generateGoatRepositoryInsert(rdbms, table, columns) + "\n\n\n"
	s += serv.generateGoatRepositoryUpdate(rdbms, table, columns) + "\n\n\n"
	s += serv.generateGoatRepositoryDelete(rdbms, table, columns) + "\n\n\n"

	return s
}


// return "type *Repository interface { ... }"
func (serv *codegenService) generateGoatRepositoryInterfaceCode(table model.Table) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)
	return fmt.Sprintf("type %sRepository interface {\n", tnp) +
		fmt.Sprintf("\tGet() ([]model.%s, error)\n", tnp) +
		fmt.Sprintf("\tGetByPk(%s *model.%s) (model.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tInsert(%s *model.%s) (model.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%s *model.%s) (model.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%s *model.%s) error\n", tni, tnp) + "}\n"
}


// generateGoatRepositoryGet generate repository function 'Get'.
// return "func (ur *userRepository) Get() ([]model.User, error) {...}"
func (serv *codegenService) generateGoatRepositoryGet(table model.Table, columns []model.Column) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf("func (%sr *%sRepository) Get() ([]model.%s, error) {\n", tni, tnc, tnp) +
		fmt.Sprintf("\tvar ret []model.%s\n\n\trows, err := %sr.db.Query(\n", tnp, tni)

	s += "\t\t`SELECT\n"
	for i, c := range columns {
		if i == 0 {
			s += fmt.Sprintf("\t\t\t%s", c.ColumnName)
		} else {
			s += fmt.Sprintf("\n\t\t\t,%s", c.ColumnName)
		}
	}
	s += "\n\t\t\t,created_at\n\t\t\t,updated_at" +
		fmt.Sprintf("\n\t\t FROM %s`,\n\t)\n\n", tn) +
		"\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\tfor rows.Next() {\n" +
		fmt.Sprintf("\t\t%s := model.%s{}\n\t\terr = rows.Scan(\n", tni, tnp)

	for _, c := range columns {
		s += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
	}
	s += fmt.Sprintf("\t\t\t&%s.CreatedAt,\n", tni) + fmt.Sprintf("\t\t\t&%s.UpdatedAt,\n", tni)

	s += fmt.Sprintf("\t\t)\n\t\tif err != nil {\n\t\t\tbreak\n\t\t}\n\t\tret = append(ret, %s)\n", tni) +
		"\t}\n\n\treturn ret, err\n}"

	return s
}


// generateGoatRepositoryGetByPk generate repository function 'GetByPk'.
// return "func (ur *userRepository) GetByPk(u *model.User) (model.User, error) {...}"
func (serv *codegenService) generateGoatRepositoryGetByPk(
	rdbms string, table model.Table, columns []model.Column,
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (%sr *%sRepository) GetByPk(%s *model.%s) (entity.%s, error) {\n", 
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
	s += serv.generateGoatRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"
	s += serv.generateGoatRepositoryWhereClauseBindVals(table, columns)
	s += "\t).Scan(\n"
	for _, c := range columns {
		s += fmt.Sprintf("\t\t&ret.%s,\n", SnakeToPascal(c.ColumnName))
	}
	s += "\t\t&retCreatedAt,\n\t\t&retUpdatedAt,\n" +
		"\t)\n\n\treturn ret, err\n}"

	return s
}


func (serv *codegenService) getBindVar(rdbms string, n int) string {
	if rdbms == "postgresql" {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}


// concatBindVariableWithCommas return ?,?,?,?,... or $1,$2,$3,$4,...
func (serv *codegenService) concatBindVariableWithCommas(rdbms string, bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, serv.getBindVar(rdbms, i))
	}
	return strings.Join(ls, ",")
}


// generateGoatRepositoryInsert generate repository function 'Insert'.
// return "func (ur *userRepository) Insert(u *entity.User) error {...}"
func (serv *codegenService) generateGoatRepositoryInsert(
	rdbms string, table model.Table, columns []model.Column
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (%sr *%sRepository) Insert(%s *model.%s) (model.%s, error) {\n", 
		tni, tnc, tni, tnp, tnp,
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
	s += fmt.Sprintf("\n\t\t ) VALUES(%s)`,\n", serv.concatBindVariableWithCommas(rdbms, bindCount))

	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateGoatRepositoryUpdate generate repository function 'Update'.
// return "func (ur *userRepository) Update(u *entity.User) error {...}"
func (serv *codegenService) generateGoatRepositoryUpdate(
	rdbms string, table model.Table, columns []model.Column
) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	
	s := fmt.Sprintf(
		"func (%sr *%sRepository) Update(%s *model.%s) (model.%s, error) {\n", 
		tni, tnc, tni, tnp, tnp,
	) + fmt.Sprintf("\t_, err := %sr.db.Exec(\n", tni) + fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET\n", tn)

	bindCount := 0
	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		c.PrimaryKeyFlg != constant.FLG_ON {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("\t\t\t%s = %s", c.ColumnName, serv.getBindVar(rdbms, bindCount))
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s = %s", c.ColumnName, serv.getBindVar(rdbms, bindCount))
			}
		}
	}

	s += "\n"
	s += serv.generateGoatRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"

	for _, c := range columns {
		if c.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		c.PrimaryKeyFlg != constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(c.ColumnName))
		}
	}
	s += serv.generateGoatRepositoryWhereClauseBindVals(table, columns)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateGoatRepositoryDelete generate repository function 'Delete'.
// return "func (ur *userRepository) Delete(u *entity.User) error {...}"
func (serv *codegenService) generateGoatRepositoryDelete(rdbms string, table model.Table, columns []model.Column) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	
	s := fmt.Sprintf(
		"func (%sr *%sRepository) Delete(%s *model.%s) (model.%s, error) {\n", 
		tni, tnc, tni, tnp, tnp,
	) + fmt.Sprintf("\t_, err := %sr.db.Exec(\n", tni) + fmt.Sprintf("\t\t`DELETE FROM %s\n", tn)

	bindCount := 0
	s += serv.generateGoatRepositoryWhereClause(rdbms, columns, &bindCount)
	s += "`,\n"
	s += serv.generateGoatRepositoryWhereClauseBindVals(table, columns)
	s += "\t)\n\n\treturn err\n}"

	return s
}


func (serv *codegenService) generateGoatRepositoryWhereClause(
	rdbms string, columns []model.Column, bindCount *int,
) string {
	s := "\t\t WHERE "

	isFirst := true
	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			*bindCount += 1
			if isFirst {
				s += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVar(rdbms, *bindCount))
				isFirst = false
			} else {
				s += fmt.Sprintf("\n\t\t   AND %s = %s", col.ColumnName, serv.getBindVar(rdbms, *bindCount))
			}
		}
	}

	return s
}


func (serv *codegenService) generateGoatRepositoryWhereClauseBindVals(
	table model.Table, columns []model.Column,
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


func (serv *codegenService) generateGoatDaoSelectScanVars(columns []model.Column, prefix string,) string {
	s := ""
	for _, col := range columns {
		s += fmt.Sprintf("%s%s,\n", prefix, SnakeToPascal(col.ColumnName))
	}

	s += fmt.Sprintf("%sCreatedAt,\n", prefix)
	s += fmt.Sprintf("%sUpdatedAt,\n", prefix)

	return s
}
