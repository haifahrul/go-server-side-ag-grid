package go_aggrid

import (
	"fmt"
	"log"
	"strings"
)

// mySQL struct
type mySQL struct{}

// BuildQuery for build query
func (*mySQL) BuildQuery(r RequestAgGrid, table string) string {
	selectSQL := MySQL.createSelectSQL(r)
	fromSQL := fmt.Sprintf("FROM %s ", table)
	whereSQL := MySQL.createWhereSQL(r)
	limitSQL := MySQL.createLimitSQL(r)
	orderBySQL := MySQL.createOrderBySQL(r)
	groupBySQL := MySQL.createGroupBySQL(r)

	SQL := fmt.Sprintf("%s %s %s %s %s %s", selectSQL, fromSQL, whereSQL, groupBySQL, orderBySQL, limitSQL)

	return SQL
}

func (*mySQL) createSelectSQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	valueCols := r.ValueCols
	groupKeys := r.GroupKeys

	isDoingGrouping := MySQL.isDoingGrouping(rowGroupCols, groupKeys)
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

func (*mySQL) createFilterSQL(key string, item map[string]interface{}) string {
	switch item["filterType"] {
	case "text":
		return MySQL.createTextFilterSQL(key, item)
	case "number":
		return MySQL.createNumberFilterSQL(key, item)
	default:
		log.Println("unkonwn filter type: ", item["filterType"])
		return ""
	}
}

func (*mySQL) createTextFilterSQL(key string, item map[string]interface{}) string {
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
		log.Println("unknown text filter type: ", item["type"])
		return "true"
	}
}

func (*mySQL) createNumberFilterSQL(key string, item map[string]interface{}) string {
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
		log.Println("unknown number filter type: ", item["type"])
		return "true"
	}
}

func (*mySQL) createWhereSQL(r RequestAgGrid) string {
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

					createFilterSQL := MySQL.createFilterSQL(i, v2.(map[string]interface{}))
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
				createFilterSQL := MySQL.createFilterSQL(i, v.(map[string]interface{}))
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

func (*mySQL) createGroupBySQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys

	isDoingGrouping := MySQL.isDoingGrouping(rowGroupCols, groupKeys)
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

func (*mySQL) createOrderBySQL(r RequestAgGrid) string {
	rowGroupCols := r.RowGroupCols
	groupKeys := r.GroupKeys
	sortModel := r.SortModel
	grouping := MySQL.isDoingGrouping(rowGroupCols, groupKeys)

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

func (*mySQL) isDoingGrouping(r []ColumnVO, g []string) bool {
	// we are not doing grouping if at the lowest level. we are at the lowest level
	// if we are grouping by more columns than we have keys for (that means the user
	// has not expanded a lowest level group, OR we are not grouping at all).
	return len(r) > len(g)
}

func (*mySQL) createLimitSQL(r RequestAgGrid) string {
	startRow := r.StartRow
	endRow := r.EndRow
	pageSize := endRow - startRow

	return fmt.Sprintf("LIMIT %v OFFSET %v", (pageSize + 1), startRow)
}

// GetRowCount for get row count
func (*mySQL) GetRowCount(r RequestAgGrid, rows int) int64 {
	if rows == 0 {
		return 0
	}

	currentLastRow := r.StartRow + int64(rows)

	if currentLastRow <= r.EndRow {
		return currentLastRow
	}
	return -1
}

// CutResultsToPageSize func
func (*mySQL) CutResultsToPageSize(r RequestAgGrid, rows []interface{}) interface{} {
	pageSize := r.EndRow - r.StartRow
	rowsLength := len(rows)

	if rowsLength != 0 && int64(rowsLength) > pageSize {
		// TODO: convert this to go
		// return rows.splice(0, pageSize)
	}
	return rows
}

// MySQL var
var MySQL = &mySQL{}
