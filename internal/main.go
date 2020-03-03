package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haifahrul/go-server-side-ag-grid/builder"
	"github.com/jmoiron/sqlx"
)

// Model struct
type Model struct {
	Description   *string `json:"description"`
	Athlete       string  `json:"athlete"`
	Age           int32   `json:"age"`
	Country       string  `json:"country"`
	Country_Group *string `json:"country_group"`
	Year          int32   `json:"year"`
	Date          string  `json:"date"`
	Sport         string  `json:"sport"`
	Gold          int64   `json:"gold"`
	Silver        int64   `json:"silver"`
	Bronze        int64   `json:"bronze"`
	Total         int64   `json:"total"`
}

// ResponseAgGrid struct for Ag Grid
type ResponseAgGrid struct {
	LastRow int64   `json:"lastRow"`
	Rows    []Model `json:"rows"`
}

// DBConn connection
var db *sqlx.DB

func main() {
	http.HandleFunc("/mysql-olympic-winners", List) // Query For SQL

	// TODO: using query Mongo
	// http.HandleFunc("/mongo-olympic-winners", ListMongo) // Query For Mongo

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

// ConnectSqlx connection
func ConnectSqlx() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "guest:guest@(127.0.0.1:3306)/sample_data")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

// List with method post
func List(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rows []Model
		var err error
		var req builder.RequestAgGrid

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		db, err = ConnectSqlx()
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer db.Close()

		w.Header().Set("Content-Type", "application/json")

		// buildSQL
		SQL := builder.MySQL.BuildQuery(req, "olympic_winners")
		log.Println("\n\n------ START QUERY BUILDER -----")
		log.Println(SQL)
		log.Println("======= END QUERY BUILDER ======")

		err = db.Select(&rows, SQL)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// rowCount := builder.MySQL.GetRowCount(req, rw)
		// log.Println("rowCount : ", rowCount)
		// resultsForPage := builder.MySQL.CutResultsToPageSize(req, rows)
		// log.Println(resultsForPage)

		response := ResponseAgGrid{
			LastRow: 100,
			Rows:    rows,
		}

		result, err := json.Marshal(response)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
