package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func main() {
	db, err := ConnectSqlx()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var rows []Model

	qryStr := `SELECT * FROM olympic_winners LIMIT 1`
	err = db.Select(&rows, qryStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, err := json.Marshal(rows)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Println(string(result))
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
