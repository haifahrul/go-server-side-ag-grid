package mysql.builder

import {
	"strings"
}

type AgGridQueryBuilder struct{}

// ColumnVO struct
type ColumnVO struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Field       string `json:"field"`
	AggFunc     string `json:"aggFunc"`
}

// RequestAgGrid struct for Ag Grid
type RequestAgGrid struct {
	StartRow     int64                    `json:"startRow"`
	EndRow       int64                    `json:"endRow"`
	RowGroupCols []ColumnVO               `json:"rowGroupCols"`
	ValueCols    []ColumnVO               `json:"valueCols"`
	PivotCols    []ColumnVO               `json:"pivotCols"`
	PivotMode    bool                     `json:"pivotMode"`
	GroupKeys    []string                 `json:"groupKeys"`
	FilterModel  map[string]interface{}   `json:"filterModel"`
	SortModel    []map[string]interface{} `json:"sortModel"`
}

// ResponseAgGrid struct for Ag Grid
type ResponseAgGrid struct {
	LastRow int64   `json:"lastRow"`
	Rows    []Model `json:"rows"`
}

func buildSQL(r RequestAgGrid, table string) string {
	selectSQL := createSelectSQL(r)
	fromSQL := fmt.Sprintf("FROM %s ", table)
	whereSQL := createWhereSQL(r)
	limitSQL := createLimitSQL(r)
	orderBySQL := createOrderBySQL(r)
	groupBySQL := createGroupBySQL(r)

	SQL := fmt.Sprintf("%s %s %s %s %s %s", selectSQL, fromSQL, whereSQL, groupBySQL, orderBySQL, limitSQL)

	log.Println("\n\n------ START QUERY BUILDER -----")
	log.Println(SQL)
	log.Println("======= END QUERY BUILDER ======\n")

	return SQL
}

func createSelectSQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	valueCols := r.ValueCols
	groupKeys := r.GroupKeys

	isDoingGrouping := isDoingGrouping(rowGroupCols, groupKeys)
	if isDoingGrouping {
		groupKeysLength := len(groupKeys)
		rowGroupCol := rowGroupCols[groupKeysLength]
		colsToSelect := make([]interface{}, 0)
		colsToSelect = append(colsToSelect, rowGroupCol.Field)

		for _, v := range valueCols {
			s := fmt.Sprintf(`%s(%s) AS %s`, v.AggFunc, v.Field, v.Field)
			colsToSelect = append(colsToSelect, s)
		}

		strs := make([]string, len(colsToSelect))
		for i, v := range colsToSelect {
			strs[i] = v.(string)
		}
		part := strings.Join(strs, ", ")

		return fmt.Sprintf("SELECT %s", part)
	}

	return "SELECT *"
}

func createFilterSQL(key string, item map[string]interface{}) string {
	switch item["filterType"] {
	case "text":
		return createTextFilterSQL(key, item)
	case "number":
		return createNumberFilterSQL(key, item)
	default:
		log.Println("unkonwn filter type: %s", item["filterType"])
		return ""
	}
}

func createTextFilterSQL(key string, item map[string]interface{}) string {
	switch item["type"] {
	case "equals":
		return fmt.Sprintf(`%s = '%s'`, key, item["filter"])
	case "notEqual":
		return fmt.Sprintf(`%s != '%s'`, key, item["filter"])
	case "contains":
		return fmt.Sprintf(`%s LIKE '%s%s%s'`, key, "%", item["filter"], "%")
	case "notContains":
		return fmt.Sprintf(`%s NOT LIKE '%s%s%s'`, key, "%", item["filter"], "%")
	case "startsWith":
		return fmt.Sprintf(`%s LIKE '%s%s'`, key, item["filter"], "%")
	case "endsWith":
		return fmt.Sprintf(`%s LIKE '%s%s'`, key, "%", item["filter"])
	default:
		log.Println("unknown text filter type: %s", item["type"])
		return "true"
	}
}

func createNumberFilterSQL(key string, item map[string]interface{}) string {
	switch item["type"] {
	case "equals":
		return fmt.Sprintf(`%s = %v`, key, item["filter"])
	case "notEqual":
		return fmt.Sprintf(`%s != %v`, key, item["filter"])
	case "greaterThan":
		return fmt.Sprintf(`%s > %v`, key, item["filter"])
	case "greaterThanOrEqual":
		return fmt.Sprintf(`%s >= %v`, key, item["filter"])
	case "lessThan":
		return fmt.Sprintf(`%s < %v`, key, item["filter"])
	case "lessThanOrEqual":
		return fmt.Sprintf(`%s <= %v`, key, item["filter"])
	case "inRange":
		return fmt.Sprintf(`(%s >= %v AND %s <= %v)`, key, item["filter"], key, item["filterTo"])
	default:
		log.Println("unknown number filter type: %s", item["type"])
		return "true"
	}
}

func createWhereSQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys
	filterModel := r.FilterModel

	whereParts := make([]string, 0)

	if len(groupKeys) > 0 {
		for k, v := range groupKeys {
			colName := rowGroupCols[k].Field
			part := fmt.Sprintf(`%s = "%s"`, colName, v)
			whereParts = append(whereParts, part)
		}
	}

	if filterModel != nil {
		for i, v := range filterModel {
			inRange := v.(map[string]interface{})
			operator := inRange["operator"]
			if operator == "AND" || operator == "OR" {
				partRange := make([]string, 0)
				for i2, v2 := range inRange {
					if i2 == "operator" {
						continue
					}

					createFilterSQL := createFilterSQL(i, v2.(map[string]interface{}))
					partRange = append(partRange, createFilterSQL)
				}

				strs := make([]string, 0)
				for _, v3 := range partRange {
					strs = append(strs, v3)
				}
				part := strings.Join(strs, fmt.Sprintf(" %s ", operator.(string)))

				wherePartRange := fmt.Sprintf(" %s ", part)
				whereParts = append(whereParts, wherePartRange)
			} else {
				createFilterSQL := createFilterSQL(i, v.(map[string]interface{}))
				whereParts = append(whereParts, createFilterSQL)
			}
		}
	}

	if len(whereParts) > 0 {
		strs := make([]string, len(whereParts))
		for i, v := range whereParts {
			strs[i] = v
		}
		part := strings.Join(strs, " AND ")

		return fmt.Sprintf(" WHERE %s ", part)
	}

	return ""
}

func createGroupBySQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys

	isDoingGrouping := isDoingGrouping(rowGroupCols, groupKeys)
	if isDoingGrouping {
		colsToGroupBy := make([]interface{}, 0)
		rowGroupCol := rowGroupCols[len(groupKeys)]
		colsToGroupBy = append(colsToGroupBy, rowGroupCol.Field)

		strs := make([]string, len(colsToGroupBy))
		for i, v := range colsToGroupBy {
			strs[i] = v.(string)
		}
		part := strings.Join(strs, ", ")
		return fmt.Sprintf(` GROUP BY %s`, part)
	}

	// select all columns
	return ""
}

// TODO:
func createOrderBySQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys
	sortModel := r.SortModel
	grouping := isDoingGrouping(rowGroupCols, groupKeys)

	sortParts := make([]string, 0)
	if len(sortModel) != 0 {
		groupColIds := make([]string, 0)
		for _, v := range rowGroupCols {
			id := v.ID
			groupColIds = append(groupColIds, id)
			break
		}

		for _, v := range sortModel {
			var groupColIdsIndexOf int
			for ig, vg := range groupColIds {
				if v["colId"] == vg {
					groupColIdsIndexOf = ig
					break
				} else {
					groupColIdsIndexOf = -1
					break
				}
			}

			if grouping && groupColIdsIndexOf < 0 {
				// ignore
			} else {
				part := fmt.Sprintf("%s %s", v["colId"], v["sort"])
				sortParts = append(sortParts, part)
			}
		}
	}

	if len(sortParts) > 0 {
		strs := make([]string, len(sortParts))
		for i, v := range sortParts {
			strs[i] = v
		}
		part := strings.Join(strs, ", ")
		return fmt.Sprintf(` ORDER BY %s`, part)
	}

	return ""
}

func isDoingGrouping(r []ColumnVO, g []string) bool {
	// we are not doing grouping if at the lowest level. we are at the lowest level
	// if we are grouping by more columns than we have keys for (that means the user
	// has not expanded a lowest level group, OR we are not grouping at all).
	return len(r) > len(g)
}

func createLimitSQL(r RequestAgGrid) string {
	startRow := r.StartRow
	endRow := r.EndRow
	pageSize := endRow - startRow

	return fmt.Sprintf("LIMIT %v OFFSET %v", (pageSize + 1), startRow)
}

func getRowCount(r RequestAgGrid, rows []Model) int64 {
	rowsLength := len(rows)

	log.Println("getRowCount : ", len(rows))

	if len(rows) == 0 {
		return 0
	}

	currentLastRow := r.StartRow + int64(rowsLength)

	if currentLastRow <= r.EndRow {
		return currentLastRow
	}
	return -1
}

func cutResultsToPageSize(r RequestAgGrid, rows []Model) interface{} {
	pageSize := r.EndRow - r.StartRow
	rowsLength := len(rows)

	if rowsLength != 0 && int64(rowsLength) > pageSize {
		// TODO: convert this to go
		// return rows.splice(0, pageSize)
	}
	return rows
}

var mysqlBuilder = MysqlBuilder