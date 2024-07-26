package models

type Scanner interface {
	Scan(args ...any) error
}

type ValuesFunc func() []any

func Scan(row Scanner, funcs ...ValuesFunc) error {
	values := []any{}
	for _, fn := range funcs {
		values = append(values, fn()...)
	}
	return row.Scan(values...)
}
