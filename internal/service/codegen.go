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


func (serv *codegenService) CodeGenerateDdl(dbType string, tableIds []int) string {
	path := "./tmp/ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
    	"-" + utils.RandomString(7) + ".sql"

	serv.generateDdlSource(dbType, tableIds, path)

	return path
}


func (serv *codegenService) CodeGenerateGoat(dbType string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
    	"-" + utils.RandomString(7)

	serv.generateGoatSource(dbType, tableIds, path)

	err := exec.Command("zip", "-r", path + ".zip", path).Run() 

	if err != nil {
		logger.LogError(err.Error())
	}

	return path + ".zip"
}


/* ############################################## */
/* ############## generateDdlSource ############# */
/* ############################################## */

func (serv *codegenService) generateDdlSource(dbType string, tableIds []int, path string) {
	ddl := serv.generateDdlCreateTables(dbType, tableIds) + 
	serv.generateDdlCreateTriggers(dbType, tableIds)

    serv.writeFile(path, ddl)
}


func (serv *codegenService) generateDdlCreateTables(dbType string, tableIds []int) string {
	ddl := ""

	for _, tid := range tableIds {
		ddl += serv.generateDdlTable(dbType, tid) + "\n\n"
	}

	return ddl
}


func (serv *codegenService) generateDdlTable(dbType string, tid int) string {
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	ddl := "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n"
	ddl += serv.generateDdlColumns(dbType, tid)
	ddl += "\n);"

	return ddl
}


func (serv *codegenService) generateDdlColumns(dbType string, tid int) string {
	columns, err := serv.cRep.SelectByTableId(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	ddl := ""
	for _, col := range columns {
		ddl += "\t" + serv.generateDdlColumn(dbType, col) + ",\n"
	}

	ddl += serv.generateDdlCommonColumns(dbType)

	pk := serv.generateDdlPrymaryKey(dbType, columns)
	if pk == "" {
		ddl = strings.TrimRight(ddl, ",\n")
	} else {
		ddl += "\t" + pk
	}

	return ddl
}


func (serv *codegenService) generateDdlColumn(dbType string, col entity.Column) string {
	ddl := col.ColumnName
	ddl += " "
	ddl += serv.generateDdlColumnDataType(dbType, col)
	ddl += " "
	ddl += serv.generateDdlColumnConstraints(dbType, col)
	ddl = strings.TrimRight(ddl, " ")
	ddl += " "
	ddl += serv.generateDdlColumnDefault(dbType, col)
	return strings.TrimRight(ddl, " ")
}


func (serv *codegenService) generateDdlColumnDataType(dbType string, col entity.Column) string {
	ddl := ""

	if dbType == "sqlite3" {
		ddl = dataTypeMapSqlite3[col.DataTypeCls]

	} else if dbType == "postgresql" {
		ddl = dataTypeMapPostgresql[col.DataTypeCls]

		if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
		col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {

			if col.Precision != 0 {
				ddl += "(" + strconv.Itoa(col.Precision) + ")"
			}	
		}

		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
			if col.Precision != 0 {
				ddl += "(" + strconv.Itoa(col.Precision)
				ddl += "," + strconv.Itoa(col.Scale) + ")"
			}
		}
	}
	return ddl
}


func (serv *codegenService) generateDdlColumnConstraints(dbType string, col entity.Column) string {
	ddl := ""

	if col.NotNullFlg == constant.FLG_ON {
		ddl += "NOT NULL"
	}

	if col.UniqueFlg == constant.FLG_ON {
		if ddl != "" {
			ddl += " "
		} 
		ddl += "UNIQUE"
	} 
	
	return ddl
}


func (serv *codegenService) generateDdlPrymaryKey(dbType string, columns []entity.Column) string {
	ddl := "PRIMARY KEY("
	pkc := 0

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			if pkc != 0 {
				ddl += ", "
			}
			ddl += col.ColumnName
			pkc++
		}
	}
	
	if pkc == 0 {
		ddl = ""
	} else {
		ddl += ")"
	}

	return ddl
}


func (serv *codegenService) generateDdlColumnDefault(dbType string, col entity.Column) string {
	if col.DefaultValue != "" {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC ||
		col.DataTypeCls == constant.DATA_TYPE_CLS_INTEGER {
			return "DEFAULT " + col.DefaultValue
		} else {
			return "DEFAULT '" + col.DefaultValue + "'"
		}
	}

	return ""
}


func (serv *codegenService) generateDdlCommonColumns(dbType string) string {
	if dbType == "sqlite3" {
		return "\tcreate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n" + 
			"\tupdate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n"

	} else if dbType == "postgresql" {
		return "\tcreate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"
	}

	return ""
}



func (serv *codegenService) generateDdlCreateTriggers(dbType string, tableIds []int) string {
	ddl := ""

	ddl += serv.generateDdlFunction(dbType)

	for _, tid := range tableIds {
		ddl += serv.generateDdlTrigger(dbType, tid) + "\n\n"
	}

	return ddl
}


func (serv *codegenService) generateDdlFunction(dbType string) string {
	if dbType == "postgresql" {
		return "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	return ""
}


func (serv *codegenService) generateDdlTrigger(dbType string, tid int) string {
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	if dbType == "sqlite3" {
		return "CREATE TRIGGER IF NOT EXISTS " + table.TableName + "_update_trg " + 
			"AFTER UPDATE ON " + table.TableName + 
			"\nBEGIN\n\tUPDATE " + table.TableName + 
			"\n\tSET update_at = DATETIME('now', 'localtime')\n" + 
			"\n\tWHERE rowid == NEW.rowid;\nEND;"

	} else if dbType == "postgresql" {
		return "CREATE TRIGGER " + table.TableName + "_update_trg " + 
			"AFTER UPDATE ON " + table.TableName + " FOR EACH ROW" + 
			"\n\texecute procedure set_update_time()" 
	}

	return ""
}


/* ############################################## */
/* ############# generateGoatSource ############# */
/* ############################################## */

func (serv *codegenService) generateGoatSource(dbType string, tableIds []int, path string) {
	mePath := path + "/model/entity"
	if err := os.MkdirAll(mePath, 0777); err != nil {
		logger.LogError(err.Error())
	}

	mrPath := path + "/model/repository"
	if err := os.MkdirAll(mrPath, 0777); err != nil {
		logger.LogError(err.Error())
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

		serv.generateGoatEntitySource(table.TableName, columns, mePath)
		serv.generateGoatRepositorySource(dbType, table.TableName, columns, mrPath)
		//serv.generateGoatController(tableName, cPath)
		//serv.generateGoatService(tableName, sPath)
	}	
}

// (user) => user.go  (user_name) => user-name.go
func (serv *codegenService) tableNameToFileName(tableName string) string {
	n := strings.ToLower(tableName)
	n = strings.Replace(n, "_", "-", -1)
	return n + ".go"
}


// (user) => User  (user_name) => UserName
func (serv *codegenService) snakeToCamelCase(snake string) string {
	n := strings.ToLower(snake)
	ls := strings.Split(n, "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}


// (user) => user  (user_name) => userName
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


func (serv *codegenService) generateGoatEntitySource(
	tableName string, columns []entity.Column, path string,
) {
	entityName := serv.snakeToCamelCase(tableName)
	path += "/" + entityName + ".go"
	entity := serv.generateGoatEntity(entityName, columns)
	serv.writeFile(path, entity)
}


func (serv *codegenService) generateGoatEntity(
	entityName string, columns []entity.Column,
) string {
	entity := "package entity\n\n\n" +
	"type " + entityName + " struct {\n"

	for _, col := range columns {
		entity += "\t" + serv.snakeToCamelCase(col.ColumnName) + " " +
			dbDataTypeGoTypeMap[col.DataTypeCls] + " " +
			"`db:\"" + strings.ToLower(col.ColumnName) + "\" " +
			"json:\"" + strings.ToLower(col.ColumnName) + "\"`\n"
	}

	return entity + "}"
}


func (serv *codegenService) generateGoatRepositorySource(
	dbType, tableName string, columns []entity.Column, path string,
) {
	path += "/" + serv.tableNameToFileName(tableName)
	repository := serv.generateGoatRepository(dbType, tableName, columns)
	serv.writeFile(path, repository)
}


func (serv *codegenService) generateGoatRepository(
	dbType, tableName string, columns []entity.Column,
) string {
	entityName := serv.snakeToCamelCase(tableName)
	repoIName := entityName + "Repository" 
	repoName := serv.snakeToLowerCamelCase(tableName) + "Repository" 

	repository := "package repository\n\n\n" +
		"import (\n" + 
		"\t\"database/sql\"\n\n" +
		"\t\"xxxxx/internal/core/db\"\n" +
		"\t\"xxxxx/internal/model/entity\"\n)\n\n\n"
	
	repository += serv.generateGoatRepositoryInterface(tableName, entityName, repoIName, columns)

	repository += "\n\n\n" +
		"type " + repoName + " struct {\n" + "\tdb *sql.DB\n}" +
		"\n\n\n" +
		"func New" + repoIName + "() " + repoIName + " {\n" +
		"\tdb := db.GetDB()\n" + 
		"\treturn &" + repoName + "{db}\n}" +
		"\n\n\n"

	rep := serv.generateGoatRepositoryInsert(dbType, tableName, entityName, repoName, columns)
	if rep != "" {
		repository += rep + "\n\n\n"
	} 


	rep = serv.generateGoatRepositorySelect(dbType, tableName, entityName, repoName, columns)
	if rep != "" {
		repository += rep + "\n\n\n"
	}

	rep = serv.generateGoatRepositoryUpdate(dbType, tableName, entityName, repoName, columns)
	if rep != "" {
		repository += rep + "\n\n\n"
	} 

	rep = serv.generateGoatRepositoryDelete(dbType, tableName, entityName, repoName, columns)
	repository += rep


	return repository
}


func (serv *codegenService) getPrimaryKeys(
	columns []entity.Column,
) []entity.Column {
	var pkcols []entity.Column

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			pkcols = append(pkcols, col)
		}
	}

	return pkcols
}

/*
tableName: user 

columnName: user_id => id
columnName: user_name => name
columnName: age => age
columnName: company_id => companyId
columnName: user_second_name => secondName
*/
func (serv *codegenService) columnNameToVariableName(
	tableName, columnName string,
) string {
	match, _ := regexp.MatchString("^" + tableName + "_.+", columnName)
	if match {
		columnName = strings.TrimLeft(columnName, tableName + "_")
	}

	return serv.snakeToLowerCamelCase(columnName)
}


/*
entityName: User => u
entityName: UserProject => up
*/
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


func (serv *codegenService) generateGoatRepositoryInterface(
	tableName, entityName, repoIName string, columns []entity.Column,
) string {
	ret := "type " + repoIName + " interface {\n"

	ret += "\t" + serv.generateGoatRepositoryInterfaceInsert(entityName) + "\n"

	args := serv.generateGoatRepositoryInterfaceCommonArgs(tableName, columns)
	if args != "" {
		ret += "\t" + serv.generateGoatRepositoryInterfaceSelect(args, entityName) + "\n"
		ret += "\t" + serv.generateGoatRepositoryInterfaceUpdate(args, entityName) + "\n"
		ret += "\t" + serv.generateGoatRepositoryInterfaceDelete(args, entityName) + "\n"
	}

	return ret + "}"
}

func (serv *codegenService) generateGoatRepositoryInterfaceCommonArgs(
	tableName string, columns []entity.Column,
) string {
	pkcols := serv.getPrimaryKeys(columns)

	args := ""
	for i, col := range pkcols {
		if i > 0 {
			args += ", "
		}
		args += serv.columnNameToVariableName(tableName, col.ColumnName)
		args += " " + dbDataTypeGoTypeMap[col.DataTypeCls]
	}

	return args
}


func (serv *codegenService) generateGoatRepositoryInterfaceInsert(
	entityName string,
) string {
	return "Insert(" + serv.entityNameToVariableName(entityName) + 
		" *entity." + entityName + ") error"
}

func (serv *codegenService) generateGoatRepositoryInterfaceSelect(
	commonArgs string, entityName string,
) string {
	return "Select(" + commonArgs + ") (entity." + entityName + ", error)"
}


func (serv *codegenService) generateGoatRepositoryInterfaceUpdate(
	commonArgs string, entityName string,
) string {
	return "Update(" + commonArgs + ", " + serv.entityNameToVariableName(entityName) + 
		" *entity." + entityName + ") error"
}

func (serv *codegenService) generateGoatRepositoryInterfaceDelete(
	commonArgs string, entityName string,
) string {
	return "Delete(" + commonArgs + ") error"
}


func (serv *codegenService) generateGoatRepositoryInsert(
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
		c ++

		if dbType == "sqlite3" {
			bvars = append(bvars, "?")
		} else if dbType == "postgresql" {
			bvars = append(bvars, "$" + strconv.Itoa(c))
		}
		
		cols = append(cols, col.ColumnName)
		bvals = append(bvals, serv.snakeToCamelCase(col.ColumnName))
	}

	ret := "func (rep *" + repoName + ") " +
		serv.generateGoatRepositoryInterfaceInsert(entityName) + " {\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`INSERT INTO " + tableName + " (\n"

	for _, col := range cols {
		ret += "\t\t\t" + col + ",\n"
	}

	ret = strings.TrimRight(ret, ",\n")
	ret += "\n\t\t ) VALUES("

	for i, bvar := range bvars {
		if i > 0 {
			ret += ","
		} 
		ret += bvar
	}

	ret += ")`,\n"

	ev := serv.entityNameToVariableName(entityName)
	for _, bval := range bvals {
		ret += "\t\t" + ev + "." + bval + ",\n"
	}

	ret += "\t)\n\n\treturn err\n}"

	return ret
}


func (serv *codegenService) generateGoatRepositorySelect(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
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

	args := serv.generateGoatRepositoryInterfaceCommonArgs(tableName, columns)
	ret := "func (rep *" + repoName + ") " +
		serv.generateGoatRepositoryInterfaceSelect(args, entityName) + " {\n" +
		"\tvar ret entity." + entityName + "\n\n" +
		"\terr := rep.db.QueryRow(\n" + 
		"\t\t`SELECT\n"

	for _, col := range cols {
		ret += "\t\t\t" + col + ",\n"
	}

	ret = strings.TrimRight(ret, ",\n")
	ret += "\n\t\t FROM " + tableName + "\n" +
	"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			ret += cond + "\n"
		} else {
			ret += "\t\t   AND " + cond + "\n"
		}
	}

	ret = strings.TrimRight(ret, "\n")
	ret += "`,\n"

	for _, bval := range bvals {
		ret += "\t\t" + bval + ",\n"
	}

	ret += "\t).Scan(\n"

	for _, scan := range scans {
		ret += "\t\t&ret." + scan + ",\n"
	}

	ret += "\t)\n\n\treturn ret, err\n}"	

	return ret

}


func (serv *codegenService) generateGoatRepositoryUpdate(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
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

	args := serv.generateGoatRepositoryInterfaceCommonArgs(tableName, columns)
	ret := "func (rep *" + repoName + ") " +
		serv.generateGoatRepositoryInterfaceUpdate(args, entityName) + " {\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`UPDATE " + tableName + "\n" +
		"\t\t SET\n"

	for _, set := range sets {
		ret += "\t\t\t" + set + ",\n"
	}

	ret = strings.TrimRight(ret, ",\n")
	ret += "\n\t\t FROM " + tableName + "\n" +
	"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			ret += cond + "\n"
		} else {
			ret += "\t\t   AND " + cond + "\n"
		}
	}

	ret = strings.TrimRight(ret, "\n")
	ret += "`,\n"

	ev := serv.entityNameToVariableName(entityName)
	for _, bset := range bsets {
		ret += "\t\t" + ev + "." + bset + ",\n"
	}

	for _, bcond := range bconds {
		ret += "\t\t" + bcond + ",\n"
	}

	ret += "\t)\n\n\treturn err\n}"	

	return ret
}


func (serv *codegenService) generateGoatRepositoryDelete(
	dbType, tableName, entityName, repoName string, columns []entity.Column,
) string {
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

	args := serv.generateGoatRepositoryInterfaceCommonArgs(tableName, columns)
	ret := "func (rep *" + repoName + ") " +
		serv.generateGoatRepositoryInterfaceDelete(args, entityName) + " {\n" +
		"\tvar ret entity." + entityName + "\n\n" +
		"\t_, err := rep.db.Exec(\n" + 
		"\t\t`DELETE FROM " + tableName + "\n" +
		"\t\t WHERE "

	for i, cond := range conds {
		if i == 0 {
			ret += cond + "\n"
		} else {
			ret += "\t\t   AND " + cond + "\n"
		}
	}

	ret = strings.TrimRight(ret, "\n")
	ret += "`,\n"

	for _, bval := range bvals {
		ret += "\t\t" + bval + ",\n"
	}

	ret += "\t)\n\n\treturn err\n}"	

	return ret

}