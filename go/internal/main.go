package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

// RequestAgGrid struct
type RequestAgGrid struct {
	StartRow     int32                    `json:"startRow"`
	EndRow       int32                    `json:"endRow"`
	RowGroupCols []map[string]interface{} `json:"rowGroupCols"`
	ValueCols    []map[string]interface{} `json:"valueCols"`
	PivotCols    []map[string]interface{} `json:"pivotCols"`
	PivotMode    bool                     `json:"pivotMode"`
	GroupKeys    []map[string]interface{} `json:"groupKeys"`
	FilterModel  interface{}              `json:"filterModel"`
	SortModel    []map[string]interface{} `json:"sortModel"`
}

// DBConn connection
var db *sqlx.DB

func main() {
	http.HandleFunc("/olympic-winners", List)
	http.HandleFunc("/olympic-winners-2", List2)

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

// List with method get
func List(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var rows []Model
		var err error

		db, err = ConnectSqlx()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		w.Header().Set("Content-Type", "application/json")

		qryStr := `SELECT * FROM olympic_winners LIMIT 10`
		err = db.Select(&rows, qryStr)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(rows)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// List2 with method post
func List2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rows []Model
		var err error
		var req RequestAgGrid

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println(err.Error())
		}

		db, err = ConnectSqlx()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		w.Header().Set("Content-Type", "application/json")

		// buildSQL
		SQL := buildSQL(req, "olympic_winners")

		err = db.Select(&rows, SQL)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(rows)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func buildSQL(r RequestAgGrid, table string) (q string) {
	selectSQL := createSelectSQL(r)
	fromSQL := fmt.Sprintf("FROM %s ", table)
	whereSQL := createWhereSQL(r)
	limitSQL := createLimitSQL(r)
	orderBySQL := createOrderBySQL(r)
	groupBySQL := createGroupBySQL(r)

	SQL := fmt.Sprintf("%s %s %s %s %s %s", selectSQL, fromSQL, whereSQL, groupBySQL, orderBySQL, limitSQL)

	log.Println("------ START QUERY BUILDER -----")
	log.Println(SQL)
	log.Println("======= END QUERY BUILDER ======")

	return SQL
}

func createSelectSQL(r RequestAgGrid) (q string) {
	rowGroupCols := r.RowGroupCols
	valueCols := r.ValueCols
	groupKeys := r.GroupKeys

	isDoingGrouping := isDoingGrouping(rowGroupCols, groupKeys)
	if isDoingGrouping {
		groupKeysLength := len(groupKeys)
		rowGroupCol := rowGroupCols[groupKeysLength]
		colsToSelect := make([]interface{}, 0)
		colsToSelect = append(colsToSelect, rowGroupCol["field"])

		for _, v := range valueCols {
			s := fmt.Sprintf(`%s(%s) AS %s`, v["aggFunc"], v["field"], v["field"])
			colsToSelect = append(colsToSelect, s)
		}

		strs := make([]string, len(colsToSelect))
		for i, v := range colsToSelect {
			strs[i] = v.(string)
		}
		q := strings.Join(strs, ", ")

		return fmt.Sprintf("SELECT %s", q)
	}

	return "SELECT *"
}

// createFilterSql(key, item) {
// 	switch (item.filterType) {
// 		case 'text':
// 			return this.createTextFilterSql(key, item);
// 		case 'number':
// 			return this.createNumberFilterSql(key, item);
// 		default:
// 			console.log('unkonwn filter type: ' + item.filterType);
// 	}
// }

// createNumberFilterSql(key, item) {
// 	switch (item.type) {
// 		case 'equals':
// 			return key + ' = ' + item.filter;
// 		case 'notEqual':
// 			return key + ' != ' + item.filter;
// 		case 'greaterThan':
// 			return key + ' > ' + item.filter;
// 		case 'greaterThanOrEqual':
// 			return key + ' >= ' + item.filter;
// 		case 'lessThan':
// 			return key + ' < ' + item.filter;
// 		case 'lessThanOrEqual':
// 			return key + ' <= ' + item.filter;
// 		case 'inRange':
// 			return '(' + key + ' >= ' + item.filter + ' and ' + key + ' <= ' + item.filterTo + ')';
// 		default:
// 			console.log('unknown number filter type: ' + item.type);
// 			return 'true';
// 	}
// }

// createTextFilterSql(key, item) {
// 	switch (item.type) {
// 		case 'equals':
// 			return key + ' = "' + item.filter + '"';
// 		case 'notEqual':
// 			return key + ' != "' + item.filter + '"';
// 		case 'contains':
// 			return key + ' like "%' + item.filter + '%"';
// 		case 'notContains':
// 			return key + ' not like "%' + item.filter + '%"';
// 		case 'startsWith':
// 			return key + ' like "' + item.filter + '%"';
// 		case 'endsWith':
// 			return key + ' like "%' + item.filter + '"';
// 		default:
// 			console.log('unknown text filter type: ' + item.type);
// 			return 'true';
// 	}
// }

func createWhereSQL(r RequestAgGrid) (q string) {
	// 	const rowGroupCols = request.rowGroupCols;
	// 	const groupKeys = request.groupKeys;
	// 	const filterModel = request.filterModel;

	// 	const that = this;
	// 	const whereParts = [];

	// 	if (groupKeys.length > 0) {
	// 		groupKeys.forEach(function (key, index) {
	// 			const colName = rowGroupCols[index].field;
	// 			whereParts.push(colName + ' = "' + key + '"')
	// 		});
	// }

	// 	if (filterModel) {
	// 		const keySet = Object.keys(filterModel);
	// 		keySet.forEach(function (key) {
	// 			const item = filterModel[key];
	// 			whereParts.push(that.createFilterSql(key, item));
	// 		});
	// 	}

	// 	if (whereParts.length > 0) {
	// 		return ' where ' + whereParts.join(' and ');
	// 	} else {
	return ""
	// 	}
}

func createGroupBySQL(r RequestAgGrid) (q string) {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys

	isDoingGrouping := isDoingGrouping(rowGroupCols, groupKeys)
	if isDoingGrouping {
		// colsToGroupBy := []map[string]interface{};

		// rowGroupCol := rowGroupCols[groupKeys.length];
		// colsToGroupBy.push(rowGroupCol.field);

		// return " group by " + colsToGroupBy.join(", ");
	}

	// select all columns
	return ""
}

func createOrderBySQL(r RequestAgGrid) (q string) {
	// 	const rowGroupCols = request.rowGroupCols;
	// 	const groupKeys = request.groupKeys;
	// 	const sortModel = request.sortModel;

	// 	const grouping = this.isDoingGrouping(rowGroupCols, groupKeys);

	// 	const sortParts = [];
	// 	if (sortModel) {

	// 		const groupColIds =
	// 			rowGroupCols.map(groupCol => groupCol.id)
	// 				.slice(0, groupKeys.length + 1);

	// 		sortModel.forEach(function (item) {
	// 			if (grouping && groupColIds.indexOf(item.colId) < 0) {
	// 				// ignore
	// 			} else {
	// 				sortParts.push(item.colId + ' ' + item.sort);
	// 			}
	// 		});
	// 	}

	// 	if (sortParts.length > 0) {
	// 		return ' order by ' + sortParts.join(', ');
	// 	}

	return ""
}

func isDoingGrouping(rowGroupCols []map[string]interface{}, groupKeys []map[string]interface{}) bool {
	// we are not doing grouping if at the lowest level. we are at the lowest level
	// if we are grouping by more columns than we have keys for (that means the user
	// has not expanded a lowest level group, OR we are not grouping at all).
	return len(rowGroupCols) > len(groupKeys)
	// return rowGroupCols.length > groupKeys.length;
}

func createLimitSQL(r RequestAgGrid) (q string) {
	startRow := r.StartRow
	endRow := r.EndRow
	pageSize := endRow - startRow

	return fmt.Sprintf("LIMIT %v OFFSET %v", (pageSize + 1), startRow)
}

// getRowCount(request, results) {
// 	if (results === null || results === undefined || results.length === 0) {
// 		return null;
// 	}
// 	const currentLastRow = request.startRow + results.length;
// 	return currentLastRow <= request.endRow ? currentLastRow : -1;
// }

// cutResultsToPageSize(request, results) {
// 	const pageSize = request.endRow - request.startRow;
// 	if (results && results.length > pageSize) {
// 		return results.splice(0, pageSize);
// 	} else {
// 		return results;
// 	}
// }
