package storage

type Scanner interface {
	Scan(args ...any) error
}

type ScannableModel interface {
	Values() []any
}

func Scan(row Scanner, models ...ScannableModel) error {
	values := []any{}
	for _, model := range models {
		values = append(values, model.Values()...)
	}
	return row.Scan(values...)
}
