package models

import (
	"strings"
	"testing"
)

func BatchInserts() {
	sqlStr := "INSERT INTO test(n1, n2, n3) VALUES "
	vals := []interface{}{}

	for _, row := range data {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row["v1"], row["v2"], row["v3"])
	}
	//trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)

	//format all vals at once
	res, _ := stmt.Exec(vals...)
}

func TestGetProducts(t *testing.T) {

}
