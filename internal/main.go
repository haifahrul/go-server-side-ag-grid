package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haifahrul/go-server-side-ag-grid/builder"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
var dbMysql *sqlx.DB
var dbPgsql *sqlx.DB

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		mysqlHost     = os.Getenv("MYSQL_HOST")
		mysqlPort     = os.Getenv("MYSQL_PORT")
		mysqlUser     = os.Getenv("MYSQL_USER")
		mysqlPassword = os.Getenv("MYSQL_PASSWORD")
		mysqlDbname   = os.Getenv("MYSQL_DBNAME")
	)
	connStr := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDbname,
	)
	dbMysql, err = ConnectMySqlx(connStr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer dbMysql.Close()

	var (
		pgHost     = os.Getenv("PG_HOST")
		pgPort     = os.Getenv("PG_PORT")
		pgUser     = os.Getenv("PG_USER")
		pgPassword = os.Getenv("PG_PASSWORD")
		pgDbname   = os.Getenv("PG_DBNAME")
		pgSslmode  = os.Getenv("PG_SSLMODE")
	)
	connPgStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUser, pgPassword, pgHost, pgPort, pgDbname, pgSslmode,
	)
	dbPgsql, err = ConnectPgSqlx(connPgStr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer dbPgsql.Close()

	http.HandleFunc("/mysql-olympic-winners", ListMySQL)         // Query For MySQL
	http.HandleFunc("/postgre-olympic-winners", ListViaPostgres) // Query For PostgresSQL

	// TODO: using query Mongo
	// http.HandleFunc("/mongo-olympic-winners", ListMongo) // Query For Mongo

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

// ConnectMySqlx connection
func ConnectMySqlx(c string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", c)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

// ConnectPgSqlx connection
func ConnectPgSqlx(c string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", c)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

// ListMySQL with method post
func ListMySQL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rows []Model
		var err error
		var req builder.RequestAgGrid

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// buildSQL
		SQL := builder.MySQL.BuildQuery(req, "olympic_winners")
		// log.Println("\n\n------ START QUERY BUILDER -----")
		// log.Println(SQL)
		// log.Println("======= END QUERY BUILDER ======")

		err = dbMysql.Select(&rows, SQL)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsLength := len(rows)
		rowCount := builder.MySQL.GetRowCount(req, rowsLength)
		// resultsForPage := builder.MySQL.CutResultsToPageSize(req, rows)
		// log.Println(resultsForPage)

		response := ResponseAgGrid{
			LastRow: rowCount,
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

// ListViaPostgres with method post
func ListViaPostgres(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rows []Model
		var err error
		var req builder.RequestAgGrid

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// buildSQL
		SQL := builder.PostgreSQL.BuildQuery(req, `"public"."olympic_winners"`)
		log.Println("\n\n======= POSTGRE SQL =======")
		log.Println(SQL)
		log.Println("======= END POSTGRE SQL ======")

		err = dbPgsql.Select(&rows, SQL)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsLength := len(rows)
		rowCount := builder.MySQL.GetRowCount(req, rowsLength)

		response := ResponseAgGrid{
			LastRow: rowCount,
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
