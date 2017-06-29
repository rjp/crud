package crud

import (
	stdsql "database/sql"
	"errors"
	"fmt"
	"github.com/rjp/crud/sql"
)

func Update(exec ExecFn, record interface{}) (stdsql.Result, error) {
	table, err := NewTable(record)
	if err != nil {
		return nil, err
	}

	pk := table.PrimaryKeyField()
	if pk == nil {
		return nil, errors.New(fmt.Sprintf("Table '%s' (%s) doesn't have a primary-key field", table.Name, table.SQLName))
	}

	return exec(sql.UpdateQuery(table.SQLName, pk.SQL.Name, table.SQLUpdateColumnSet()), table.SQLUpdateValueSet()...)
}

func MustUpdate(exec ExecFn, record interface{}) error {
	result, err := Update(exec, record)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("No rows matching")
	}

	return nil
}
