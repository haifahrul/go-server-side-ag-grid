package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Model struct {
	Country string `json="country"`
	Gold    int64  `json="gold"`
	Silver  int64  `json="silver"`
	Bronze  int64  `json="bronze"`
}

func main() {
	// var err error

	db, err := sqlx.Connect("mysql", "username=guest password=guest dbname=sample_data sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	}

	model := Model{}
	rows, err := db.Queryx("SELECT * FROM olympic_winners")
	for rows.Next() {
		err := rows.StructScan(&model)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Printf("%#v\n", model)
	}

	log.Println(model)
}
