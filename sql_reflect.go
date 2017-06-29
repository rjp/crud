package crud

import (
	"github.com/rjp/crud/meta"
	"github.com/azer/snakecase"
)

func SQLTableNameOf(st interface{}) string {
	return snakecase.SnakeCase(meta.TypeNameOf(st))
}
