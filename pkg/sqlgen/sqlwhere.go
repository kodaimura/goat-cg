package pkg

import (
    "strconv"
)


type SqlWhere string


func (s SqlWhere) Where(field , op string, value interface{}) SqlWhere {
    return "where " + makeSqlCondition(field, op, value)
}


func (s SqlWhere) And(field , op string, value interface{}) SqlWhere {
    return s + " And " + makeSqlCondition(field, op, value)
}

func (s SqlWhere) Or(field , op string, value interface{}) SqlWhere {
    return s + " Or " + makeSqlCondition(field, op, value)
}

func makeSqlCondition(field , op string, value interface{}) SqlWhere {
    var v string
    switch value.(type) {
        case int:
            v = strconv.Itoa(value.(int))
        case string:
            v = "'" + value.(string) + "'"
        case nil:
            v = "NULL"
    }
    return SqlWhere(field) + " " + SqlWhere(op) + " " + SqlWhere(v) 
}