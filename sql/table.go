package sql

import (
	"fmt"
	"reflect"
	"strings"
)

func NewTableQuery(st interface{}, name string, fields []*Options, ifNotExists bool) string {
	ifNotExistsExt := ""
	if ifNotExists {
		ifNotExistsExt = " IF NOT EXISTS"
	}

    return fmt.Sprintf("CREATE TABLE%s `%s` (\n%s%s\n)%s;",
		ifNotExistsExt, name, NewFieldQueries(fields), NewPrimaryKeyQuery(st, fields), NewTableConfigQuery(fields))
}

func NewFieldQueries(fields []*Options) string {
	queries := []string{}

	for _, f := range fields {
		if f.Ignore {
			continue
		}

		queries = append(queries, NewFieldQuery(f))
	}

	return strings.Join(queries, ",\n")
}

func NewFieldQuery(field *Options) string {
	length := ""
	autoIncrement := ""
	required := ""
	defaultValue := ""
	unsigned := ""
	unique := ""

	if field.Length > -1 {
		length = fmt.Sprintf("(%d)", field.Length)
	}

	if field.AutoIncrement > 0 {
		autoIncrement = " AUTO_INCREMENT"
	}

	if field.IsRequired {
		required = " NOT NULL"
	}

	if field.DefaultValue != "" {
		defaultValue = fmt.Sprintf(" DEFAULT %s", field.DefaultValue)
	}

	if field.IsUnsigned {
		unsigned = " UNSIGNED"
	}

	if field.IsUnique {
		unique = " UNIQUE"
	}

	query := fmt.Sprintf("%s%s%s%s%s%s",
		length, required, defaultValue, unsigned, unique, autoIncrement)

	return fmt.Sprintf("  `%s` %s%s", field.Name, field.Type, query)
}

func CallHook(st interface{}, method string, arg string) string {
    v := reflect.ValueOf(st)
    m := v.MethodByName(method)
    if m.IsValid() {
        args := []reflect.Value{reflect.ValueOf(arg)}
        retVal := m.Call(args)
        fish := retVal[0]
        if fish.String() != "" {
            arg = fish.String()
        }
        fmt.Println("CallHook()", arg)
    }
	return arg
}

func NewPrimaryKeyQuery(st interface{}, fields []*Options) string {
	keys := []string{}

	for _, f := range fields {
		if f.IsPrimaryKey {
			keys = append(keys, f.Name)
		}
	}

	if len(keys) == 0 {
		return ""
	}

    sql := fmt.Sprintf(",\n  PRIMARY KEY (`%s`)", strings.Join(keys, "`, `"))
	sql = CallHook(st, "PrimaryKeyHook", sql)
	return sql
}

func NewTableConfigQuery(fields []*Options) string {
	autoIncrement := ""
	for _, f := range fields {
		if f.AutoIncrement > 1 {
			autoIncrement = fmt.Sprintf(" AUTO_INCREMENT=%d", f.AutoIncrement)
		}
	}

	return fmt.Sprintf("%s", autoIncrement)
}

func DropTableQuery(name string, ifExists bool) string {
	ext := ""

	if ifExists {
		ext = " IF EXISTS"
	}

	return fmt.Sprintf("DROP TABLE%s %s", ext, name)
}

func ShowTablesLikeQuery(name string) string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s'", name)
}

func InsertQuery(st interface{}, tableName string, columnNames []string) string {
	var questionMarks string

	if len(columnNames) > 0 {
		questionMarks = strings.Repeat("?,", len(columnNames))
		questionMarks = questionMarks[:len(questionMarks)-1]
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(columnNames, ","), questionMarks)
	sql = CallHook(st, "InsertHook", sql)
	return sql
}

func SelectQuery(tableName string, columnNames []string) string {
	columns := strings.Join(columnNames, ",")
	if columns == "" {
		columns = "*"
	}

	return fmt.Sprintf("SELECT %s FROM %s", columns, tableName)
}

func CompleteSelectQuery(tableName string, columnNames []string, original string) string {
	if strings.HasPrefix(original, "SELECT ") || strings.HasPrefix(original, "select ") {
		return original
	}

	if len(original) > 0 {
		original = " " + original
	}

	return fmt.Sprintf("%s%s", SelectQuery(tableName, columnNames), original)
}

func UpdateQuery(tableName, index string, columnNames []string) string {
	return fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=?", tableName, strings.Join(columnNames, "=?, "), index)
}

func DeleteQuery(tableName, index string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s=?", tableName, index)
}
