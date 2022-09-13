package service

import (
	"os"
	"time"
	"strconv"
	"strings"
	"regexp"
	"unicode"
	"os/exec"

	"goat-cg/pkg/utils"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
)


type CodegenService interface {
	CodeGenerateGoat(dbType string, tableIds []int) string
	CodeGenerateDdl(dbType string, tableIds []int) string
}


type codegenService struct {
	cRep repository.ColumnRepository
	tRep repository.TableRepository
}


func NewCodegenService() CodegenService {
	cRep := repository.NewColumnRepository()
	tRep := repository.NewTableRepository()
	return &codegenService{cRep, tRep}
}


// CodeGenerateDdl generate ddl(create table) source 
// and return file path.
// param dbType: "sqlite3" or "postgresql" 
func (serv *codegenService) CodeGenerateDdl(dbType string, tableIds []int) string {
	path := "./tmp/ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
    	"-" + utils.RandomString(7) + ".sql"

	serv.cgDdlSource(dbType, tableIds, path)

	return path
}


// CodeGenerateGoat generate programs(entity, repository for goat) 
// and return zip path .
// param dbType: "sqlite3" or "postgresql" 
func (serv *codegenService) CodeGenerateGoat(dbType string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
    	"-" + utils.RandomString(7)

	serv.cgGoatSource(dbType, tableIds, path)

	if err := exec.Command("zip", "-r", path + ".zip", path).Run(); err != nil {
		logger.LogError(err.Error())
	}

	return path + ".zip"
}


//dataTypeMapSqlite3 map DataTypeCls and sqlite3 data types.
var dataTypeMapSqlite3 = map[string]string {
	constant.DATA_TYPE_CLS_SERIAL: "INTEGER AUTOINCREMENT",
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
	constant.DATA_TYPE_CLS_SERIAL: "SERIAL",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TIMESTAMP",
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
		logger.LogError(err.Error())
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.LogError(err.Error())
	}
}


func (serv *codegenService) extractPrimaryKeys(columns []entity.Column) []entity.Column {
	var ret []entity.Column

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			ret = append(ret, col)
		}
	}

	return ret
}


///////////////////////
/// CodeGenerateDdl ///
///////////////////////

// cgDdlSource generate ddl(create table) source.
// main processing of CodeGenerateDdl.
func (serv *codegenService) cgDdlSource(dbType string, tableIds []int, path string) {
	s := serv.cgDdlCreateTables(dbType, tableIds) + "\n" +
		serv.cgDdlCreateTriggers(dbType, tableIds)

    serv.writeFile(path, s)
}


func (serv *codegenService) cgDdlCreateTables(dbType string, tableIds []int) string {
	s := ""
	for _, tid := range tableIds {
		s += serv.cgDdlCreateTable(dbType, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) cgDdlCreateTable(dbType string, tid int) string {
	s := ""
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
		return s
	}

	s += "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n" +
		serv.cgDdlColumns(dbType, tid) + "\n);"

	return s
}


func (serv *codegenService) cgDdlColumns(dbType string, tid int) string {
	s := ""
	columns, err := serv.cRep.SelectByTableId(tid)

	if err != nil {
		logger.LogError(err.Error())
		return s
	}

	for _, col := range columns {
		s += serv.cgDdlColumn(dbType, col)
	}
	s += serv.cgDdlCommonColumns(dbType)
	s += serv.cgDdlPrymaryKey(dbType, columns)

	return strings.TrimRight(s, ",\n")
}


func (serv *codegenService) cgDdlCommonColumns(dbType string) string {
	s := ""
	if dbType == "sqlite3" {
		s = "\tcreate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n" + 
			"\tupdate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n"

	} else if dbType == "postgresql" {
		s = "\tcreate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"
	}

	return s
}


func (serv *codegenService) cgDdlPrymaryKey(dbType string, columns []entity.Column) string {
	s := "" 
	pkcols := serv.extractPrimaryKeys(columns)

	for i, col := range pkcols {
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


func (serv *codegenService) cgDdlColumn(dbType string, col entity.Column) string {
	s := "\t" + col.ColumnName + " " + serv.cgDdlColumnDataType(dbType, col)
	if cts := serv.cgDdlColumnConstraints(col); cts != "" {
		s += " " + cts
	}
	if dflt := serv.cgDdlColumnDefault(col); dflt != "" {
		s += " " + dflt
	}

	return s + ",\n"
}


func (serv *codegenService) cgDdlColumnConstraints(col entity.Column) string {
	s := ""
	if col.NotNullFlg == constant.FLG_ON {
		s += "NOT NULL "
	}
	if col.UniqueFlg == constant.FLG_ON {
		s += "UNIQUE"
	} 
	
	return strings.TrimRight(s, " ")
}


func (serv *codegenService) cgDdlColumnDefault(col entity.Column) string {
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


func (serv *codegenService) cgDdlColumnDataType(dbType string, col entity.Column) string {
	s := ""
	if dbType == "sqlite3" {
		s = dataTypeMapSqlite3[col.DataTypeCls]

	} else if dbType == "postgresql" {
		s = serv.cgDdlColumnDataTypePostgresql(col)
	}

	return s
}


func (serv *codegenService) cgDdlColumnDataTypePostgresql(col entity.Column) string {
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


func (serv *codegenService) cgDdlCreateTriggers(dbType string, tableIds []int) string {
	s := ""
	if dbType == "postgresql" {
		s = "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	for _, tid := range tableIds {
		s = serv.cgDdlCreateTrigger(dbType, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) cgDdlCreateTrigger(dbType string, tid int) string {
	s := ""
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
		return s
	}

	if dbType == "sqlite3" {
		s = "CREATE TRIGGER IF NOT EXISTS " + table.TableName + "_update_trg " + 
			"AFTER UPDATE ON " + table.TableName + 
			"\nBEGIN\n\tUPDATE " + table.TableName + 
			"\n\tSET update_at = DATETIME('now', 'localtime')\n" + 
			"\n\tWHERE rowid == NEW.rowid;\nEND;"

	} else if dbType == "postgresql" {
		s = "CREATE TRIGGER " + table.TableName + "_update_trg " + 
			"AFTER UPDATE ON " + table.TableName + " FOR EACH ROW" + 
			"\n\texecute procedure set_update_time()" 
	}

	return s
}


////////////////////////
/// CodeGenerateGoat ///
////////////////////////

// tableNameToFileName get file name from tabel name
// user => user.go / user_name => user-name.go
func (serv *codegenService) tableNameToFileName(tableName string) string {
	n := strings.ToLower(tableName)
	n = strings.Replace(n, "_", "-", -1)
	return n + ".go"
}


// snakeToCamelCase
// user => User / user_name => UserName
func (serv *codegenService) snakeToCamelCase(snake string) string {
	n := strings.ToLower(snake)
	ls := strings.Split(n, "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}


// snakeToLowerCamelCase
// user => user / user_name => userName
func (serv *codegenService) snakeToLowerCamelCase(snake string) string {
	n := strings.ToLower(snake)
	ls := strings.Split(n, "_")
	for i, s := range ls {
		if i == 0 {
	        continue
	    }
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}


// cgGoatSource generate programs(entity, repository for goat).
// main processing of CodeGenerateGoat.
func (serv *codegenService) cgGoatSource(dbType string, tableIds []int, path string) {
	mePath := path + "/model/entity"
	if err := os.MkdirAll(mePath, 0777); err != nil {
		logger.LogError(err.Error())
		return
	}

	mrPath := path + "/model/repository"
	if err := os.MkdirAll(mrPath, 0777); err != nil {
		logger.LogError(err.Error())
		return
	}

	for _, tid := range tableIds {
		table, err := serv.tRep.Select(tid)
		if err != nil {
			logger.LogError(err.Error())
			break
		}

		columns, err := serv.cRep.SelectByTableId(tid)
		if err != nil {
			logger.LogError(err.Error())
			break
		}

		serv.cgGoatEntitySource(table.TableName, columns, mePath)
		serv.cgGoatRepositorySource(dbType, table.TableName, columns, mrPath)
		//serv.cgGoatController(tableName, cPath)
		//serv.cgGoatService(tableName, sPath)
	}	
}


/////////////////////////////////
/// CodeGenerateGoat (Entity) ///
/////////////////////////////////

// cgGoatEntitySource generate entity program for goat.
func (serv *codegenService) cgGoatEntitySource(
	tableName string, columns []entity.Column, path string,
) {
	entityName := serv.snakeToCamelCase(tableName)
	path += "/" + entityName + ".go"
	s := serv.cgGoatEntity(entityName, columns)
	serv.writeFile(path, s)
}


func (serv *codegenService) cgGoatEntity(
	entityName string, columns []entity.Column,
) string {
	s := "package entity\n\n\n"

	s += "type " + entityName + " struct {\n"
	for _, col := range columns {
		s += "\t" + serv.snakeToCamelCase(col.ColumnName) + " " +
			dbDataTypeGoTypeMap[col.DataTypeCls] + " " +
			"`db:\"" + strings.ToLower(col.ColumnName) + "\" " +
			"json:\"" + strings.ToLower(col.ColumnName) + "\"`\n"
	}
	s += "\tCreateAt string `db:\"create_at\" json:\"create_at\"`\n"
	s += "\tUpdateAt string `db:\"update_at\" json:\"update_at\"`\n"

	return s + "}"
}


/////////////////////////////////////
/// CodeGenerateGoat (Repository) ///
/////////////////////////////////////

// columnNameToVariableName get shorten variable name from column name.
// tableName: user
// columnName: user_id => id
// columnName: user_name => name
// columnName: age => age
// columnName: company_id => companyId
// columnName: user_second_name => secondName
func (serv *codegenService) columnNameToVariableName(
	tableName, columnName string,
) string {
	match, _ := regexp.MatchString("^" + tableName + "_.+", columnName)
	if match {
		columnName = strings.TrimLeft(columnName, tableName + "_")
	}

	return serv.snakeToLowerCamelCase(columnName)
}


// entityNameToVariableName get shorten variable name from entity name.
// entityName: User => u
// entityName: ProjectUser => up
func (serv *codegenService) entityNameToVariableName(
	entityName string,
) string {
	ret := ""
	for _, r := range entityName {
		if unicode.IsUpper(r) {
			ret += string(r)
		}
	}
	return strings.ToLower(ret)
}


// cgGoatRepositorySource generate repository program for goat.
func (serv *codegenService) cgGoatRepositorySource(
	dbType, tableName string, columns []entity.Column, path string,
) {
	path += "/" + serv.tableNameToFileName(tableName)
	s := serv.cgGoatRepository(dbType, tableName, columns)
	serv.writeFile(path, s)
}


func (serv *codegenService) cgGoatRepository(
	dbType, tableName string, columns []entity.Column,
) string {
	// EntityName
	en := serv.snakeToCamelCase(tableName)
	// RepositoryInterfaceName
	rin := en + "Repository"
	// RepositoryName
	rn := serv.snakeToLowerCamelCase(tableName) + "Repository" 

	s := "package repository\n\n\n" +
		"import (\n" + 
		"\t\"database/sql\"\n\n" +
		"\t\"xxxxx/internal/core/db\"\n" +
		"\t\"xxxxx/internal/model/entity\"\n)\n\n\n"

	s += serv.cgGoatRepositoryInterface(tableName, en, rin, columns)
	
	s += "\n\n\n" +
		"type " + rn + " struct {\n" + "\tdb *sql.DB\n}" +
		"\n\n\n" +
		"func New" + rin + "() " + rin + " {\n" +
		"\tdb := db.GetDB()\n" + 
		"\treturn &" + rn + "{db}\n}" +
		"\n\n\n"

	rep := serv.cgGoatRepositoryInsert(dbType, tableName, en, rn, columns)
	if rep != "" {
		s += rep + "\n\n\n"
	} 
	rep = serv.cgGoatRepositorySelect(dbType, tableName, en, rn, columns)
	if rep != "" {
		s += rep + "\n\n\n"
	}
	rep = serv.cgGoatRepositoryUpdate(dbType, tableName, en, rn, columns)
	if rep != "" {
		s += rep + "\n\n\n"
	} 
	rep = serv.cgGoatRepositoryDelete(dbType, tableName, en, rn, columns)
	s += rep

	return s
}


func (serv *codegenService) cgGoatRepositoryInterface(
	tableName, entityName, repoIName string, columns []entity.Column,
) string {
	s := "type " + repoIName + " interface {\n"

	s += "\t" + serv.cgGoatRepositoryInterfaceInsert(entityName) + "\n"

	args := serv.cgGoatRepositoryInterfaceCommonArgs(tableName, columns)
	if args != "" {
		s += "\t" + serv.cgGoatRepositoryInterfaceSelect(args, entityName) + "\n"
		s += "\t" + serv.cgGoatRepositoryInterfaceUpdate(args, entityName) + "\n"
		s += "\t" + serv.cgGoatRepositoryInterfaceDelete(args, entityName) + "\n"
	}

	return s + "}"
}


// cgGoatRepositoryInterfaceCommonArgs generate repository common args
// from primary key columns.
// pk: user_id, user_name => "userId int, userName string"
// if tableName is "user" omit "user". => "id int, name string"
func (serv *codegenService) cgGoatRepositoryInterfaceCommonArgs(
	tableName string, columns []entity.Column,
) string {
	s := ""
	pkcols := serv.extractPrimaryKeys(columns)

	for i, col := range pkcols {
		if i > 0 {
			s += ", "
		}
		s += serv.columnNameToVariableName(tableName, col.ColumnName)
		s += " " + dbDataTypeGoTypeMap[col.DataTypeCls]
	}

	return s
}


// cgGoatRepositoryInterfaceInsert
// return "Insert(x entity.Entity) error"
func (serv *codegenService) cgGoatRepositoryInterfaceInsert(
	entityName string,
) string {
	return "Insert(" + serv.entityNameToVariableName(entityName) + 
		" *entity." + entityName + ") error"
}


// cgGoatRepositoryInterfaceSelect
// return "Select(commonArgs) (entity.Entity, error)"
func (serv *codegenService) cgGoatRepositoryInterfaceSelect(
	commonArgs string, entityName string,
) string {
	return "Select(" + commonArgs + ") (entity." + entityName + ", error)"
}


// cgGoatRepositoryInterfaceUpdate
// return "Update(commonArgs, x entity.Entity) error"
func (serv *codegenService) cgGoatRepositoryInterfaceUpdate(
	commonArgs string, entityName string,
) string {
	return "Update(" + commonArgs + ", " + serv.entityNameToVariableName(entityName) + 
		" *entity." + entityName + ") error"
}


// cgGoatRepositoryInterfaceDelete
// return "Delete(commonArgs) error"
func (serv *codegenService) cgGoatRepositoryInterfaceDelete(
	commonArgs string, entityName string,
) string {
	return "Delete(" + commonArgs + ") error"
}


// cgGoatRepositoryInsert generate repository function 'Insert'.
// return "func (rep *repoName) Insert(x entity.Entity) error {...}"
func (serv *codegenService) cgGoatRepositoryInsert(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
	cols := []string{}
	bvars := []string{}
	bvals := []string{}

	c := 0
	for _, col := range columns {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_SERIAL {
			continue
		}
		c++

		if dbType == "sqlite3" {
			bvars = append(bvars, "?")
		} else if dbType == "postgresql" {
			bvars = append(bvars, "$" + strconv.Itoa(c))
		}
		
		cols = append(cols, col.ColumnName)
		bvals = append(bvals, serv.snakeToCamelCase(col.ColumnName))
	}

	s := "func (rep *" + repoName + ") " +
		serv.cgGoatRepositoryInterfaceInsert(entityName) + " {\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`INSERT INTO " + tableName + " (\n"

	for _, col := range cols {
		s += "\t\t\t" + col + ",\n"
	}

	s = strings.TrimRight(s, ",\n")
	s += "\n\t\t ) VALUES("

	for i, bvar := range bvars {
		if i > 0 {
			s += ","
		} 
		s += bvar
	}

	s += ")`,\n"

	ev := serv.entityNameToVariableName(entityName)
	for _, bval := range bvals {
		s += "\t\t" + ev + "." + bval + ",\n"
	}

	return s + "\t)\n\n\treturn err\n}"
}

// cgGoatRepositorySelect generate repository function 'Select'.
// return "func (rep *repoName) Select(commonArgs) (entity.Entity, error) {...}"
func (serv *codegenService) cgGoatRepositorySelect(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
	args := serv.cgGoatRepositoryInterfaceCommonArgs(tableName, columns)
	//pkがない場合
	if args == "" {
		return ""
	}

	cols := []string{}
	conds := []string{}
	bvals := []string{}
	scans := []string{}

	pkc := 0
	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			pkc++
			if dbType == "sqlite3" {
				conds = append(conds, col.ColumnName + " = ?")
			} else if dbType == "postgresql" {
				conds = append(conds, col.ColumnName + " = $" + strconv.Itoa(pkc))
			}
			bvals = append(bvals, serv.columnNameToVariableName(tableName, col.ColumnName))
		}

		cols = append(cols, col.ColumnName)
		scans = append(scans, serv.snakeToCamelCase(col.ColumnName))
	}

	cols = append(cols, "create_at")
	cols = append(cols, "update_at")
	scans = append(scans, "CreateAt")
	scans = append(scans, "UpdateAt")

	s := "func (rep *" + repoName + ") " +
		serv.cgGoatRepositoryInterfaceSelect(args, entityName) + " {\n" +
		"\tvar ret entity." + entityName + "\n\n" +
		"\terr := rep.db.QueryRow(\n" + 
		"\t\t`SELECT\n"

	for _, col := range cols {
		s += "\t\t\t" + col + ",\n"
	}

	s = strings.TrimRight(s, ",\n")
	s += "\n\t\t FROM " + tableName + "\n" +
	"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			s += cond + "\n"
		} else {
			s += "\t\t   AND " + cond + "\n"
		}
	}

	s = strings.TrimRight(s, "\n")
	s += "`,\n"

	for _, bval := range bvals {
		s += "\t\t" + bval + ",\n"
	}

	s += "\t).Scan(\n"

	for _, scan := range scans {
		s += "\t\t&ret." + scan + ",\n"
	}

	return s + "\t)\n\n\treturn ret, err\n}"
}


// cgGoatRepositoryUpdate generate repository function 'Update'.
// return "func (rep *repoName) Update(commonArgs, x entity.Entity) error {...}"
func (serv *codegenService) cgGoatRepositoryUpdate(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
	args := serv.cgGoatRepositoryInterfaceCommonArgs(tableName, columns)
	//pkがない場合
	if args == "" {
		return ""
	}

	sets := []string{}
	bsets := []string{}
	conds := []string{}
	bconds := []string{}

	c := 0
	for _, col := range columns {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_SERIAL {
			continue
		}
		c ++

		if dbType == "sqlite3" {
			sets = append(sets, col.ColumnName + " = ?")
		} else if dbType == "postgresql" {
			sets = append(sets, col.ColumnName + " = $" + strconv.Itoa(c))
		}

		bsets = append(bsets, serv.snakeToCamelCase(col.ColumnName))
	}

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			c++
			if dbType == "sqlite3" {
				conds = append(conds, col.ColumnName + " = ?")
			} else if dbType == "postgresql" {
				conds = append(conds, col.ColumnName + " = $" + strconv.Itoa(c))
			}
			bconds = append(bconds, serv.columnNameToVariableName(tableName, col.ColumnName))
		}
	}

	s := "func (rep *" + repoName + ") " +
		serv.cgGoatRepositoryInterfaceUpdate(args, entityName) + " {\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`UPDATE " + tableName + "\n" +
		"\t\t SET\n"

	for _, set := range sets {
		s += "\t\t\t" + set + ",\n"
	}

	s = strings.TrimRight(s, ",\n")
	s += "\n\t\t FROM " + tableName + "\n" +
	"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			s += cond + "\n"
		} else {
			s += "\t\t   AND " + cond + "\n"
		}
	}

	s = strings.TrimRight(s, "\n")
	s += "`,\n"

	ev := serv.entityNameToVariableName(entityName)
	for _, bset := range bsets {
		s += "\t\t" + ev + "." + bset + ",\n"
	}

	for _, bcond := range bconds {
		s += "\t\t" + bcond + ",\n"
	}

	return s + "\t)\n\n\treturn err\n}"	
}

// cgGoatRepositoryDelete generate repository function 'Delete'.
// return "func (rep *repoName) Delete(commonArgs) error {...}"
func (serv *codegenService) cgGoatRepositoryDelete(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
	args := serv.cgGoatRepositoryInterfaceCommonArgs(tableName, columns)
	//pkがない場合
	if args == "" {
		return ""
	}

	conds := []string{}
	bvals := []string{}

	pkc := 0
	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			pkc++
			if dbType == "sqlite3" {
				conds = append(conds, col.ColumnName + " = ?")
			} else if dbType == "postgresql" {
				conds = append(conds, col.ColumnName + " = $" + strconv.Itoa(pkc))
			}
			bvals = append(bvals, serv.columnNameToVariableName(tableName, col.ColumnName))
		}
	}
	
	s := "func (rep *" + repoName + ") " +
		serv.cgGoatRepositoryInterfaceDelete(args, entityName) + " {\n" +
		"\tvar ret entity." + entityName + "\n\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`DELETE FROM " + tableName + "\n" +
		"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			s += cond + "\n"
		} else {
			s += "\t\t   AND " + cond + "\n"
		}
	}

	s = strings.TrimRight(s, "\n")
	s += "`,\n"

	for _, bval := range bvals {
		s += "\t\t" + bval + ",\n"
	}

	return s + "\t)\n\n\treturn err\n}"	
}
