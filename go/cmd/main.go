package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/haifahrul/go-server-side-ag-grid/go/internal"
)

// Model struct
type Model struct {
	Athlete       string  `json="athlete"`
	Age           int32   `json="age"`
	Country       string  `json="country"`
	Country_Group *string `json="country_group"`
	Year          int32   `json="year"`
	Date          string  `json="date"`
	Sport         string  `json="sport"`
	Gold          int64   `json="gold"`
	Silver        int64   `json="silver"`
	Bronze        int64   `json="bronze"`
	Total         int64   `json="total"`
}

// DBConn connection
var db *sqlx.DB

func main() {
	db, err := ConnectSqlx()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	http.HandleFunc("/olympic-winners", OlympicWinners.FindAll)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

// ConnectSqlx connection
func ConnectSqlx() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "guest:guest@(127.0.0.1:3306)/sample_data")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}

// func buildSql(request) (q string, error) {

// selectSql := this.createSelectSql(request);
// const fromSql = " FROM sample_data.olympic_winners ";
// const whereSql = this.createWhereSql(request);
// const limitSql = this.createLimitSql(request);

// const orderBySql = this.createOrderBySql(request);
// const groupBySql = this.createGroupBySql(request);

// const SQL = selectSql + fromSql + whereSql + groupBySql + orderBySql + limitSql;

// console.log(request)
// console.log(SQL);

// 	return q, nil;
// }

// func createSelectSql(request) (q string) {
// const rowGroupCols = request.rowGroupCols;
// const valueCols = request.valueCols;
// const groupKeys = request.groupKeys;

// if (this.isDoingGrouping(rowGroupCols, groupKeys)) {
// 	const colsToSelect = [];

// 	const rowGroupCol = rowGroupCols[groupKeys.length];
// 	colsToSelect.push(rowGroupCol.field);

// 	valueCols.forEach(function (valueCol) {
// 		colsToSelect.push(valueCol.aggFunc + '(' + valueCol.field + ') as ' + valueCol.field);
// 	});

// 	return ' select ' + colsToSelect.join(', ');
// }

// return ' select *';
// }
