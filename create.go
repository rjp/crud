package crud

import (
	"github.com/rjp/crud/sql"
)

func Create(exec ExecFn, record interface{}) error {
	row, err := NewRow(record)
	if err != nil {
		return err
	}

	columns := []string{}
	values := []interface{}{}

	for c, v := range row.SQLValues() {
		columns = append(columns, c)
		values = append(values, v)
	}

	_, err = exec(sql.InsertQuery(record, row.SQLTableName, columns), values...)
	return err
}

