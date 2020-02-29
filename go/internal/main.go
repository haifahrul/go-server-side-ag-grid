package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Model struct
type Model struct {
	Athlete      string    `json="athlete"`
	Age          int32     `json="age"`
	Country      string    `json="country"`
	CountryGroup *string   `json="country_group"`
	Year         int32     `json="year"`
	Date         time.Time `json="date"`
	Sport        string    `json="sport"`
	Gold         int64     `json="gold"`
	Silver       int64     `json="silver"`
	Bronze       int64     `json="bronze"`
	Total        int64     `json="total"`
}

func main() {
	db, err := ConnectSqlx()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var rows []Model

	qryStr := `SELECT athlete FROM olympic_winners`
	err = db.Select(&rows, qryStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println("FAHRUL", rows)
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
